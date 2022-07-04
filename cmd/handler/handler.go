package handler

import (
	"fmt"

	"github.com/imarsman/iptools/pkg/subnet"
	"inet.af/netaddr"
)

// SubnetDivide divide a subnet into ranges
func SubnetDivide(ip string, mask uint8, secondaryMask uint8) {
	// prefix := fmt.Sprintf("%s/d", ip, mask)
	s, err := subnet.NewFromIPAndMask(ip, mask)
	if err != nil {

	}
	ranges := []netaddr.IPRange{}
	if secondaryMask != 0 {
		s2, err := subnet.NewFromIPAndMask(ip, secondaryMask)
		ranges, err = s.NetworkRangesInSubnets(s2)
		if err != nil {

		}
		// fmt.Println("secondary mask", secondaryMask, s2, ranges)
	} else {
		ranges, err = s.NetworkRanges()
		if err != nil {

		}
	}
	for _, r := range ranges {
		fmt.Println(r.String())
	}
}
