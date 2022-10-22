package util

import (
	"net"
	"net/netip"
)

// GetAddresses get addresses for a domain
func GetAddresses(domain string) (addresses []netip.Addr, err error) {
	ipList, err := net.LookupIP(domain)

	if err == nil {
		for _, ip := range ipList {
			var addr netip.Addr
			addr, err = netip.ParseAddr(ip.String())
			if err != nil {
				return
			}
			addresses = append(addresses, addr)

			addr, err = netip.ParseAddr(ip.String())
			if err != nil {
				return
			}
		}

	}

	return
}
