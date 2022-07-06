package handler

import (
	"fmt"
	"os"

	"github.com/imarsman/iptools/pkg/subnet"
	"inet.af/netaddr"
)

// SubnetDivide divide a subnet into ranges
func SubnetDivide(ip string, mask uint8, secondaryMask uint8) {
	var err error
	var s *subnet.IPV4Subnet
	prefix := fmt.Sprintf("%s/%d", ip, mask)
	s, err = subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ranges := []netaddr.IPRange{}
	if secondaryMask != 0 {
		var s2 *subnet.IPV4Subnet
		prefix := fmt.Sprintf("%s/%d", ip, secondaryMask)
		s2, err = subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ranges, err = s.NetworkRangesInSubnets(s2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		ranges, err = s.NetworkRanges()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	for _, r := range ranges {
		fmt.Println(r.String())
	}
}
