// SPDX-License-Identifier: GPL-3.0-or-later

package iprange

import (
	"fmt"
	"math/big"
	"net/netip"
)

// Family represents IP Range address-family.
type Family uint8

const (
	// V4Family is IPv4 address-family.
	V4Family Family = iota
	// V6Family is IPv6 address-family.
	V6Family
)

// Range represents an IP range.
type Range interface {
	Family() Family
	Contains(ip netip.Addr) bool
	Size() *big.Int
	fmt.Stringer
	Ips() []string
}

// New returns new IP Range.
// If it is not a valid range (start and end IPs have different address-families, or start > end),
// New returns nil.
func New(start, end netip.Addr) Range {
	return ipRange{start: start, end: end}
}

type ipRange struct {
	start netip.Addr
	end   netip.Addr
}

func (i ipRange) Family() Family {
	if i.start.Is6() {
		return V6Family
	}

	return V4Family
}

func (i ipRange) Contains(ip netip.Addr) bool {
	if i.start.Compare(ip) <= 0 && i.end.Compare(ip) >= 0 {
		return true
	}
	return false
}

func (i ipRange) Size() *big.Int {
	if i.start.Is4() && i.end.Is4() {
		return big.NewInt(v4ToInt(i.end.AsSlice()) - v4ToInt(i.start.AsSlice()) + 1)
	} else if i.start.Is6() && i.end.Is6() {
		size := big.NewInt(0)
		size.Add(size, big.NewInt(0).SetBytes(i.end.AsSlice()))
		size.Sub(size, big.NewInt(0).SetBytes(i.start.AsSlice()))
		size.Add(size, big.NewInt(1))
		return size
	}

	return big.NewInt(0)
}

func v4ToInt(ip []byte) int64 {
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

func (i ipRange) String() string {
	return fmt.Sprintf("%s-%s", i.start, i.end)
}

func (i ipRange) Ips() []string {
	return rangeIps(i.start, i.end)
}

func rangeIps(start, end netip.Addr) []string {

	var ips []string
	for start.Compare(end) <= 0 {
		ips = append(ips, start.String())
		start = start.Next()
	}

	return ips
}
