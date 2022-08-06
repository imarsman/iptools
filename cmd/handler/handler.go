package handler

import (
	"fmt"
	"net/netip"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/pkg/ipv4subnet"
	"github.com/imarsman/iptools/pkg/ipv4subnet/util"
)

// IP4SubnetDescribe describe a subnet
func IP4SubnetDescribe(ip string, mask uint8) {
	var err error
	var s *ipv4subnet.Subnet
	prefix := fmt.Sprintf("%s/%d", ip, mask)
	s, err = ipv4subnet.NewFromPrefix(prefix)
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

	// Get subnet mask
	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet Mask")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.SubnetMask().Addr())},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Wildcard Mask")},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", util.WildCardMask(netip.Addr(s.SubnetMask().Addr())))},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	if s.Networks() > 0 {
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

		networkAddress, err := s.BroadcastAddr()
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
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.BinaryMask())},
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
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s.in-addr.arpa", util.InAddrArpa(s.Prefix().Addr()))},
	}
	table.Body.Cells = append(table.Body.Cells, r)

	fmt.Println(table.String())
}

// IP4SubnetRanges divide a subnet into ranges
func IP4SubnetRanges(ip string, bits uint8, secondaryMask uint8) {
	var err error
	var s *ipv4subnet.Subnet
	prefix := fmt.Sprintf("%s/%d", ip, bits)
	s, err = ipv4subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ranges := []ipv4subnet.Range{}
	var s2 *ipv4subnet.Subnet
	if secondaryMask != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryMask)
		s2, err = ipv4subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ranges, err = s.SecondaryIPRanges(s2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		s2 = s
		ranges, err = s.IPRanges()
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
		table.Body.Cells = append(table.Body.Cells, r)

		networkAddress, err := s.BroadcastAddr()
		if err != nil {
			return
		}
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Broadcast Address")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", networkAddress.String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet Mask")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.SubnetMask().Addr())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		if secondaryMask != 0 {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.CIDR())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet IP")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.IP().String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)

			networkAddress, err := s.BroadcastAddr()
			if err != nil {
				return
			}
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet Broadcast Address")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", networkAddress.String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)

			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet Mask")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.SubnetMask().Addr())},
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

		fmt.Println()
		table = simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Start"},
				{Align: simpletable.AlignCenter, Text: "End"},
			},
		}
		for _, r := range ranges {
			cell := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", r.First().String())},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", r.Last().String())},
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

// IP4SubnetDivide divide a subnet into ranges
func IP4SubnetDivide(ip string, mask uint8, secondaryMask uint8) {
	var err error
	var s *ipv4subnet.Subnet
	prefix := fmt.Sprintf("%s/%d", ip, mask)
	s, err = ipv4subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	subnets := []*ipv4subnet.Subnet{}
	var s2 *ipv4subnet.Subnet
	if secondaryMask != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryMask)
		s2, err = ipv4subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		subnets, err = s.SecondarySubnets(s2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		s2 = s
		subnets, err = s.SecondarySubnets(s)
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
		table.Body.Cells = append(table.Body.Cells, r)

		networkAddress, err := s.BroadcastAddr()
		if err != nil {
			return
		}
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Broadcast Address")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", networkAddress.String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet Mask")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.SubnetMask().Addr())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		if secondaryMask != 0 {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.Prefix().String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)

			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet IP")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.IP().String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)

			networkAddress, err := s.BroadcastAddr()
			if err != nil {
				return
			}
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet Broadcast Address")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", networkAddress.String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)

			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary Subnet Mask")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.SubnetMask().Addr())},
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
			fmt.Println(s.String())
		}
	}
}
