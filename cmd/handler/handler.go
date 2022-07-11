package handler

import (
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/pkg/subnet"
	"github.com/imarsman/iptools/pkg/util"
	"inet.af/netaddr"
)

// SubnetDescribe describe a subnet
func SubnetDescribe(ip string, mask uint8) {
	var err error
	var s *subnet.IPV4Subnet
	prefix := fmt.Sprintf("%s/%d", ip, mask)
	s, err = subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	r := []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.CIDR())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	if s.Networks() > 0 {
		// get last address for subnet
		last, err := s.Last()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get first address for subnet
		first, err := s.First()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network Address")},
			// {Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", networkAddress.String())},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", first.String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)

		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "IP Address")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", last.String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)

		networkAddress, err := s.NetworkAddress()
		if err != nil {
			return
		}
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Broadcast Address")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", networkAddress.String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network Hosts")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Hosts())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Total Hosts")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.TotalHosts())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	class := string(s.Class())
	if class == `0` {
		class = "Subnet"
	}
	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "IP Class")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", class)},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	ipType := "Public"
	if s.IP().IsPrivate() {
		ipType = "Private"
	}
	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "IP Type")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", ipType)},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Binary Subnet Mask")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.BinarySubnetMask())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Binary ID")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.BinaryID())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	last, err := s.Last()
	if err != nil {
		return
	}

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Hex ID")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", util.IPToHexStr(last))},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "in-addr.arpa")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s.in-addr.arpa", util.InAddrArpa(s.IP()))},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Wildcard Mask")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", util.WildCardMask(s.Prefix().IP()))},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	fmt.Println(table.String())
}

// SubnetRanges divide a subnet into ranges
func SubnetRanges(ip string, bits uint8, secondaryMask uint8) {
	var err error
	var s *subnet.IPV4Subnet
	prefix := fmt.Sprintf("%s/%d", ip, bits)
	s, err = subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ranges := []netaddr.IPRange{}
	var s2 *subnet.IPV4Subnet
	if secondaryMask != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryMask)
		s2, err = subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ranges, err = s.NetworkIPRangesInSubnet(s2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		s2 = s
		ranges, err = s.NetworkIPRanges()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if args.CLIArgs.Subnet.SubnetRanges.Pretty {
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Category"},
				{Align: simpletable.AlignCenter, Text: "Value"},
			},
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.CIDR())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet IP")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.IP().String())},
		}
		if secondaryMask != 0 {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.CIDR())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}
		if bits == secondaryMask {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network Hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Hosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		} else {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			if secondaryMask != 0 {
				r = []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Networks")},
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s2.Networks())},
				}
				table.Body.Cells = append(table.Body.Cells, r)
				r = []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Effective Networks")},
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", len(ranges))},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network Hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Hosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			if secondaryMask != 0 {
				r = []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Sub Network Hosts")},
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s2.Hosts())},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
		}
		fmt.Println()
		table.SetStyle(simpletable.StyleCompactLite)
		fmt.Println(table.String())
	}

	if args.CLIArgs.Subnet.SubnetRanges.Pretty {
		fmt.Println()
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Start"},
				{Align: simpletable.AlignCenter, Text: "End"},
			},
		}
		for _, r := range ranges {
			cell := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", r.From().String())},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", r.To().String())},
			}
			table.Body.Cells = append(table.Body.Cells, cell)
		}
		table.SetStyle(simpletable.StyleCompactLite)
		fmt.Println(table.String())
	} else {
		for _, r := range ranges {
			fmt.Println(r.String())
		}
	}
}

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

	subnets := []*subnet.IPV4Subnet{}
	var s2 *subnet.IPV4Subnet
	if secondaryMask != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryMask)
		s2, err = subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		subnets, err = s.NetworkSubnetsInSubnet(s2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		s2 = s
		subnets, err = s.NetworkSubnetsInSubnet(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if args.CLIArgs.Subnet.SubnetDivide.Pretty {
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Category"},
				{Align: simpletable.AlignCenter, Text: "Value"},
			},
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.Prefix().String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet IP")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.IP().String())},
		}
		if secondaryMask != 0 {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.Prefix().String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}
		if mask == secondaryMask {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network Hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Hosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		} else {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			if secondaryMask != 0 {
				r = []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Networks")},
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s2.Networks())},
				}
				table.Body.Cells = append(table.Body.Cells, r)
				r = []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Effective Networks")},
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", len(subnets))},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network Hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Hosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			if secondaryMask != 0 {
				r = []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Sub Network Hosts")},
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s2.Hosts())},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
		}
		fmt.Println()
		table.SetStyle(simpletable.StyleCompactLite)
		fmt.Println(table.String())
	}

	if args.CLIArgs.Subnet.SubnetDivide.Pretty {
		fmt.Println()
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Subnet"},
			},
		}
		// var subnetsSplit []*subnet.IPV4Subnet
		for _, s := range subnets {
			cell := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.String())},
			}
			table.Body.Cells = append(table.Body.Cells, cell)
		}
		table.SetStyle(simpletable.StyleCompactLite)

		fmt.Println(table.String())
	} else {
		for _, s := range subnets {
			// subnetNew, err := subnet.NewFromIPAndBits(r.From().IPAddr().IP.String(), s2.Prefix.Bits())
			// if err != nil {
			// 	panic(err.Error())
			// }
			fmt.Println(s.String())
		}
	}
}
