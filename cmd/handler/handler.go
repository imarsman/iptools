package handler

import (
	"fmt"
	"net/netip"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/alexeyco/simpletable"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/pkg/ipv4subnet"
	"github.com/imarsman/iptools/pkg/ipv4subnet/util"
)

var printer = message.NewPrinter(language.English)

func row(label string, value any) (r []*simpletable.Cell) {
	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%v", label)},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%v", value)},
	}
	return
}

// IP4SubnetDescribe describe a subnet
func IP4SubnetDescribe(ip string, bits uint8, secondaryBits uint8) {
	var err error
	var s *ipv4subnet.Subnet
	prefix := fmt.Sprintf("%s/%d", ip, bits)
	s, err = ipv4subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var s2 *ipv4subnet.Subnet
	if secondaryBits != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryBits)
		s2, err = ipv4subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		s2 = s
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, row("Subnet", s.CIDR()))
	table.Body.Cells = append(table.Body.Cells, row("Subnet IP", s.IP().String()))
	table.Body.Cells = append(table.Body.Cells, row("Broadcast Address", s.BroadcastAddr().String()))
	table.Body.Cells = append(table.Body.Cells, row("Subnet Mask", s.SubnetMask()))
	table.Body.Cells = append(table.Body.Cells, row("Wildcard Mask", s.WildcardMask()))

	class := string(s.Class())
	if class == `0` {
		class = "Subnet"
	}
	table.Body.Cells = append(table.Body.Cells, row("IP Class", class))

	ipType := "Public"
	if s.IP().IsPrivate() {
		ipType = "Private"
	}
	table.Body.Cells = append(table.Body.Cells, row("IP Type", ipType))
	table.Body.Cells = append(table.Body.Cells, row("Binary Subnet Mask", s.BinaryMask()))
	table.Body.Cells = append(table.Body.Cells, row("Binary ID", s.BinaryID()))

	last := s.Last()
	if !last.IsValid() {
		fmt.Println(fmt.Errorf("invalid address %s", netip.Addr{}))
	}

	table.Body.Cells = append(table.Body.Cells, row("Hex ID", util.IPToHexStr(last)))
	table.Body.Cells = append(table.Body.Cells, row("in-addr.arpa", util.InAddrArpa(s.Prefix().Addr())))

	if secondaryBits != 0 {
		table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet", s2.CIDR()))
		table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet IP", s2.IP().String()))

		if !s2.BroadcastAddr().IsValid() {
			fmt.Println(fmt.Errorf("invalid address %s", netip.Addr{}))
			return
		}
		table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Broadcast Address", s2.BroadcastAddr().String()))
		table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Mask", s2.SubnetMask()))
		table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Wildcard Mask", s2.WildcardMask()))
	}

	if secondaryBits == 0 {
		table.Body.Cells = append(table.Body.Cells, row("Networks", s.Networks()))

		table.Body.Cells = append(table.Body.Cells, row("Network Hosts", printer.Sprintf("%d", s.Hosts())))
	} else {
		table.Body.Cells = append(table.Body.Cells, row("Networks", s.Networks()))
		table.Body.Cells = append(table.Body.Cells, row("Secondary Networks", s2.Networks()))
		table.Body.Cells = append(table.Body.Cells, row("Effective Networks", s.EffectiveNetworks(s2)))

		table.Body.Cells = append(table.Body.Cells, row("Network Hosts", printer.Sprintf("%d", s.Hosts())))
		table.Body.Cells = append(table.Body.Cells, row("Secondary Network Hosts", printer.Sprintf("%d", s2.Hosts())))
	}

	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}

// IP4SubnetRanges divide a subnet into ranges
func IP4SubnetRanges(ip string, bits uint8, secondaryBits uint8) {
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
	if secondaryBits != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryBits)
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
		table.Body.Cells = append(table.Body.Cells, row("Subnet", s.CIDR()))
		table.Body.Cells = append(table.Body.Cells, row("Subnet IP", s.IP().String()))

		if !s.BroadcastAddr().IsValid() {
			fmt.Println(fmt.Errorf("invalid address %s", netip.Addr{}))
			return
		}

		table.Body.Cells = append(table.Body.Cells, row("Broadcast Address", s.BroadcastAddr().String()))
		table.Body.Cells = append(table.Body.Cells, row("Subnet Mask", s.SubnetMask()))
		table.Body.Cells = append(table.Body.Cells, row("Wildcard Mask", s.WildcardMask()))

		if secondaryBits != 0 {
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet", s2.CIDR()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet IP", s2.IP().String()))

			if !s2.BroadcastAddr().IsValid() {
				fmt.Println(fmt.Errorf("invalid address %s", netip.Addr{}))
				return
			}

			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Broadcast Address", s2.BroadcastAddr().String()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Mask", s2.SubnetMask()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Wildcard Mask", s2.WildcardMask()))
		}
		if secondaryBits == 0 {
			table.Body.Cells = append(table.Body.Cells, row("Networks", s.Networks()))
			table.Body.Cells = append(table.Body.Cells, row("Network Hosts", printer.Sprintf("%d", printer.Sprintf("%d", s.Hosts()))))
		} else {
			table.Body.Cells = append(table.Body.Cells, row("Networks", s.Networks()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Networks", s2.Networks()))
			table.Body.Cells = append(table.Body.Cells, row("Effective Networks", len(ranges)))
			table.Body.Cells = append(table.Body.Cells, row("Network Hosts", printer.Sprintf("%d", s.Hosts())))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Network Hosts", printer.Sprintf("%d", s2.Hosts())))
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
			table.Body.Cells = append(table.Body.Cells, row(r.First().String(), r.Last().String()))
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
func IP4SubnetDivide(ip string, bits uint8, secondaryBits uint8) {
	var err error
	var s *ipv4subnet.Subnet
	prefix := fmt.Sprintf("%s/%d", ip, bits)
	s, err = ipv4subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	subnets := []*ipv4subnet.Subnet{}
	var s2 = s
	if secondaryBits != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryBits)
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
		subnets, err = s.Subnets()
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
		table.Body.Cells = append(table.Body.Cells, row("Subnet", s.Prefix().String()))
		table.Body.Cells = append(table.Body.Cells, row("Subnet IP", s.IP().String()))

		if !s.BroadcastAddr().IsValid() {
			fmt.Println(fmt.Errorf("invalid address %s", netip.Addr{}))
			return
		}
		table.Body.Cells = append(table.Body.Cells, row("Broadcast Address", s.BroadcastAddr().String()))
		table.Body.Cells = append(table.Body.Cells, row("Subnet Mask", s.SubnetMask()))
		table.Body.Cells = append(table.Body.Cells, row("Wildcard Mask", s.WildcardMask()))

		if secondaryBits != 0 {
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet", s2.Prefix().String()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet IP", s2.IP().String()))

			if !s2.BroadcastAddr().IsValid() {
				fmt.Println(fmt.Errorf("invalid address %s", netip.Addr{}))
				return
			}
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Broadcast Address", s2.BroadcastAddr().String()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Mask", s2.SubnetMask()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Subnet Wildcard Mask", s2.WildcardMask()))
		}

		if secondaryBits == 0 {
			table.Body.Cells = append(table.Body.Cells, row("Networks", s.Networks()))
			table.Body.Cells = append(table.Body.Cells, row("Network Hosts", s.Hosts()))
		} else {
			table.Body.Cells = append(table.Body.Cells, row("Netorks", s.Networks()))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Networks", s2.Networks()))
			table.Body.Cells = append(table.Body.Cells, row("Effective Networks", s.EffectiveNetworks(s2)))
			table.Body.Cells = append(table.Body.Cells, row("Network Hosts", printer.Sprintf("%d", s.Hosts())))
			table.Body.Cells = append(table.Body.Cells, row("Secondary Network Hosts", printer.Sprintf("%d", s2.Hosts())))
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
				{Align: simpletable.AlignCenter, Text: "Subnets"},
			},
		}
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
