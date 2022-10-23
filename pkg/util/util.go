package util

import (
	"fmt"
	"net"
	"net/netip"
)

// GetDomainAddresses get addresses for a domain
func GetDomainAddresses(domain string) (addresses []netip.Addr, err error) {
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

// GetDomainMXRecods get MX records for a domain
func GetDomainMXRecods(domain string) (mxRecods []*net.MX, err error) {
	mxRecods, err = net.LookupMX(domain)

	return
}
