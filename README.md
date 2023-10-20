
# iprange

This package helps you to work with IP ranges,Both IPv4 and IPv6 are supported.

IP range doesn't contain network and broadcast IP addresses if the format is `IPv4 CIDR`, `IPv4 subnet mask`
or `IPv6 CIDR`. 


## Installation

Install iprange with go mod

```bash
  go get github.com/NSObjects/iprange  
```
    
## Supported formats

- `IPv4 address` (192.0.2.1)
- `IPv4 range` (192.0.2.0-192.0.2.10)
- `IPv4 CIDR` (192.0.2.0/24)
- `IPv4 subnet mask` (192.0.2.0/255.255.255.0)
- `IPv6 address` (2001:db8::1)
- `IPv6 range` (2001:db8::-2001:db8::10)
- `IPv6 CIDR` (2001:db8::/64)

## Usage/Examples

### Range 

```go
parseRange, err := iprange.ParseRange("192.0.2.1-192.0.2.10")
if err != nil {
	panic(err)
}
fmt.Println(parseRange.Ips())
}
```

#### CIDR

```go 
parseRange, err := iprange.ParseRange("192.0.2.1/24")
if err != nil {
	panic(err)
}
fmt.Println(parseRange.Ips())
}
```

### subnet mask

```go
parseRange, err := iprange.ParseRange("192.0.2.1/255.255.255.255")
if err != nil {
	panic(err)
}
fmt.Println(parseRange.Ips())
}
```
