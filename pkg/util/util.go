package util

import (
	"fmt"
	"net"
	"net/netip"
)

const (
	// GlobalUnicast IPV6 type
	GlobalUnicast = iota
	// UniqueLocal IPV6 type
	UniqueLocal
	// LinkLocalUnicast IPV6 type
	LinkLocalUnicast
	// Loopback IPV6 type
	Loopback
	// Multicast IPV6 type
	Multicast
	// InterfaceLocalMulticast IPV6 type
	InterfaceLocalMulticast
	// LinkLocalMulticast IPV6 type
	LinkLocalMulticast
	// Private IPV6 type
	Private
	// Unspecified IPV6 type
	Unspecified
	// Unknown IPV6 type
	Unknown
)

// DomainAddresses get addresses for a domain
func DomainAddresses(domain string) (addresses []netip.Addr, err error) {
	ipList, err := net.LookupIP(domain)

	if err == nil {
		for _, ip := range ipList {
			// ip will be either 4 or 16 bytes. net.IP is an alias for a byte slice
			addr, ok := netip.AddrFromSlice(ip)
			if !ok {
				err = fmt.Errorf("could not obtain address from %s", ip.String())
				return
			}
			addresses = append(addresses, addr)
		}
	}

	return
}

// AddrType get address type as int
func AddrType(addr netip.Addr) int {
	switch {
	// case strings.HasPrefix(addr.StringExpanded(), "fd00"):
	// 	return UniqueLocal
	case addr.IsInterfaceLocalMulticast(): // fe80::/10
		return InterfaceLocalMulticast
	case addr.IsLinkLocalMulticast(): // ff00::/8 ff02
		return LinkLocalMulticast
	case addr.IsLinkLocalUnicast(): // fe80::/10
		return LinkLocalUnicast
	case addr.IsLoopback(): // ::1/128
		return Loopback
	case addr.IsPrivate(): // fc00::/7
		return Private
	case addr.IsGlobalUnicast(): // 2001
		return GlobalUnicast
	case addr.IsMulticast(): // ff00::/8
		return Multicast
	case addr.IsUnspecified():
		return Unspecified
	default:
		return Unknown
	}

}

// AddrTypeName the type of address for the subnet
// https://www.networkacademy.io/ccna/ipv6/ipv6-address-types
func AddrTypeName(addr netip.Addr) string {
	switch AddrType(addr) {
	case UniqueLocal:
		return "Unique local"
	case GlobalUnicast:
		return "Global unicast"
	case InterfaceLocalMulticast:
		return "Interface local multicast"
	case LinkLocalMulticast:
		return "Link local muticast"
	case LinkLocalUnicast:
		return "Link local unicast"
	case Loopback:
		return "Loopback"
	case Multicast:
		return "Multicast"
	case Private:
		return "Private"
	case Unspecified:
		return "Unspecified"
	default:
		return "Unknown"
	}
}

// DomainMXRecods get MX records for a domain
func DomainMXRecods(domain string) (mxRecods []*net.MX, err error) {
	mxRecods, err = net.LookupMX(domain)

	return
}
