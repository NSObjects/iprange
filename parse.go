// SPDX-License-Identifier: GPL-3.0-or-later

package iprange

import (
	"fmt"
	"go4.org/netipx"
	"net"
	"net/netip"
	"regexp"
	"strings"
)

// ParseRanges parses s as a space separated list of IP Ranges, returning the result and an error if any.
// IP Range can be in IPv4 address ("192.0.2.1"), IPv4 range ("192.0.2.0-192.0.2.10")
// IPv4 CIDR ("192.0.2.0/24"), IPv4 subnet mask ("192.0.2.0/255.255.255.0"),
// IPv6 address ("2001:db8::1"), IPv6 range ("2001:db8::-2001:db8::10"),
// or IPv6 CIDR ("2001:db8::/64") form.
// IPv4 CIDR, IPv4 subnet mask and IPv6 CIDR ranges don't include network and broadcast addresses.
func ParseRanges(s string) ([]Range, error) {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil, nil
	}

	var ranges []Range
	for _, v := range parts {
		r, err := ParseRange(v)
		if err != nil {
			return nil, err
		}

		if r != nil {
			ranges = append(ranges, r)
		}
	}
	return ranges, nil
}

var (
	reRange      = regexp.MustCompile("^[0-9a-f.:-]+$")           // addr | addr-addr
	reCIDR       = regexp.MustCompile("^[0-9a-f.:]+/[0-9]{1,3}$") // addr/prefix_length
	reSubnetMask = regexp.MustCompile("^[0-9.]+/[0-9.]{7,}$")     // v4_addr/mask
)

// ParseRange parses s as an IP Range, returning the result and an error if any.
// The string s can be in IPv4 address ("192.0.2.1"), IPv4 range ("192.0.2.0-192.0.2.10")
// IPv4 CIDR ("192.0.2.0/24"), IPv4 subnet mask ("192.0.2.0/255.255.255.0"),
// IPv6 address ("2001:db8::1"), IPv6 range ("2001:db8::-2001:db8::10"),
// or IPv6 CIDR ("2001:db8::/64") form.
// IPv4 CIDR, IPv4 subnet mask and IPv6 CIDR ranges don't include network and broadcast addresses.
func ParseRange(s string) (Range, error) {
	s = strings.ToLower(s)
	if s == "" {
		return nil, nil
	}

	var r Range
	switch {
	case reRange.MatchString(s):
		r = parseRange(s)
	case reCIDR.MatchString(s):
		r = parseCIDR(s)
	case reSubnetMask.MatchString(s):
		r = parseSubnetMask(s)
	}

	if r == nil {
		return nil, fmt.Errorf("ip range (%s) invalid syntax", s)
	}
	return r, nil
}

func parseRange(s string) Range {
	ipr, err := netipx.ParseIPRange(s)
	if err != nil {
		addr, err := netip.ParseAddr(s)
		if err != nil {
			return nil
		}
		if addr.IsValid() {
			return New(addr, addr)
		}
		return nil
	}

	return New(ipr.From(), ipr.To())
}

func parseCIDR(s string) Range {
	addr, err := netip.ParsePrefix(s)
	if err != nil {
		return nil
	}
	r := netipx.RangeOfPrefix(addr)
	if r.From().Next().Compare(r.To().Prev()) < 0 {
		return parseRange(fmt.Sprintf("%s-%s", r.From().Next().String(), r.To().Prev().String()))
	}

	return parseRange(fmt.Sprintf("%s-%s", r.From().String(), r.To().String()))
}

func parseSubnetMask(s string) Range {
	idx := strings.LastIndexByte(s, '/')
	if idx == -1 {
		return nil
	}

	address, mask := s[:idx], s[idx+1:]

	stringMask := net.IPMask(net.ParseIP(mask).To4())
	length, bits := stringMask.Size()

	if length == 0 && bits == 0 {
		return nil
	}

	return parseCIDR(fmt.Sprintf("%s/%d", address, length))
}
