package handler

import (
	"fmt"
	"net/netip"
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/alexeyco/simpletable"
	"github.com/imarsman/iptools/cmd/args"

	"github.com/imarsman/iptools/pkg/ipv4subnet"
	"github.com/imarsman/iptools/pkg/ipv4subnet/ipv4util"
	"github.com/imarsman/iptools/pkg/ipv6"
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
func IP4SubnetDescribe(ip string, bits int, secondaryBits int) {
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
	table.Body.Cells = append(table.Body.Cells, row("Broadcast Address Hex ID", ipv4util.IPToHexStr(s.Last())))
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

	table.Body.Cells = append(table.Body.Cells, row("in-addr.arpa", ipv4util.InAddrArpa(s.Prefix().Addr())))

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
func IP4SubnetRanges(ip string, bits int, secondaryBits int) {
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
	if args.CLIArgs.IP4Subnet.SubnetRanges.Pretty {
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
func IP4SubnetDivide(ip string, bits int, secondaryBits int) {
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
	if args.CLIArgs.IP4Subnet.SubnetDivide.Pretty {
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

	if args.CLIArgs.IP4Subnet.SubnetDivide.Pretty {
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

const typeGlobalUnicast = "global-unicast"
const typeLinkLocal = "link-local"
const typeUniqueLocal = "unique-local"
const typePrivate = "private"
const typeMulticast = "multicast"
const typeInterfaceLocalMulticast = "interface-local-multicast"
const typeLinkLocalMulticast = "link-local-multicast"

// IP6SubnetDescribe describe a link-local address
func IP6SubnetDescribe(ip string, bits int, random bool, ip6Type string) {
	if bits == 0 {
		bits = 64
	}
	if ip6Type == "" && ip == "" {
		fmt.Println("If no IP then type must be supplied")
		os.Exit(1)
	} else if ip == "" && !random {
		fmt.Println("No type supplied and -random is false")
		os.Exit(1)
	}
	var addr netip.Addr
	if !random {
		var err error
		addr, err = netip.ParseAddr(ip)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		var err error
		if ip6Type == typeGlobalUnicast {
			addr, err = ipv6.RandomAddrGlobalUnicast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeLinkLocal {
			addr, err = ipv6.RandomAddrLinkLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typePrivate {
			addr, err = ipv6.RandomAddrPrivate()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeMulticast {
			addr, err = ipv6.RandomAddrMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeLinkLocalMulticast {
			addr, err = ipv6.RandomAddrLinkLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeInterfaceLocalMulticast {
			addr, err = ipv6.RandomAddrInterfaceLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("No valid type specified")
			os.Exit(1)
		}
	}
	// s, err := ipv6subnet.NewFromAddrAndBits(addr, bits)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	if !(addr.IsMulticast() || addr.IsInterfaceLocalMulticast() || addr.IsLinkLocalMulticast()) {
		prefix := netip.PrefixFrom(addr, bits)
		ip6SubnetDisplay(addr, prefix)
	} else {
		prefix := netip.PrefixFrom(addr, bits)
		ip6SubnetDisplayBasic(addr, prefix)
	}
}

// ip6SubnetDisplay describe a link local IP
func ip6SubnetDisplay(addr netip.Addr, prefix netip.Prefix) {
	var value string
	var err error
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, row("IP Type", ipv6.AddrTypeName(addr)))
	table.Body.Cells = append(table.Body.Cells, row("Type Prefix", ipv6.AddrTypePrefix(addr).Masked()))
	table.Body.Cells = append(table.Body.Cells, row("IP", addr.String()))
	if ipv6.HasType(ipv6.AddrType(addr), ipv6.GlobalUnicast, ipv6.LinkLocalUnicast, ipv6.UniqueLocal, ipv6.Private) {
		solicitedNodeAddr, err := ipv6.AddrSolicitedNodeMulticast(addr)
		if err != nil {
			panic(err)
		}
		table.Body.Cells = append(table.Body.Cells,
			row(
				"Solicited node multicast", solicitedNodeAddr.String(),
			),
		)
	}

	table.Body.Cells = append(table.Body.Cells, row("Prefix", prefix.Masked()))
	if ipv6.AddrType(addr) == ipv6.GlobalUnicast {
		table.Body.Cells = append(
			table.Body.Cells, row(
				"Routing Prefix", fmt.Sprintf("%s", fmt.Sprintf("%s", ipv6.RoutingPrefixString(addr)))),
		)
	}
	// Handle global id for appropriate types
	if ipv6.HasType(ipv6.AddrType(addr), ipv6.GlobalUnicast, ipv6.Private) {
		value, err = ipv6.GlobalID(addr)
		if err != nil {
			fmt.Println(err)
		}
		table.Body.Cells = append(table.Body.Cells, row("Global ID", fmt.Sprintf("%s", value)))
	}
	table.Body.Cells = append(table.Body.Cells, row("Interface ID", fmt.Sprintf("%s", ipv6.InterfaceString(addr))))
	table.Body.Cells = append(table.Body.Cells, row("Subnet ID", fmt.Sprintf("%s", ipv6.SubnetString(addr))))
	if ipv6.AddrType(addr) == ipv6.LinkLocalUnicast {
		table.Body.Cells = append(table.Body.Cells, row("Default Gateway", ipv6.LinkLocalDefaultGateway(addr)))
	}
	if ipv6.HasType(ipv6.AddrType(addr), ipv6.GlobalUnicast) {
		table.Body.Cells = append(table.Body.Cells, row("Link", ipv6.AddrLink(addr)))
	}
	if ipv6.IsARPA(addr) {
		table.Body.Cells = append(table.Body.Cells, row("ip6.arpa", fmt.Sprintf("%s", ipv6.Arpa(addr))))
	}
	table.Body.Cells = append(table.Body.Cells, row("Subnet first address", ipv6.First(addr).StringExpanded()))
	table.Body.Cells = append(table.Body.Cells, row("Subnet last address", ipv6.Last(addr).StringExpanded()))
	part := strings.Split(ipv6.AddrToBitString(addr), ".")[0]
	part = fmt.Sprintf("%s%s", strings.Repeat("0", 16-len(part)), part)
	table.Body.Cells = append(table.Body.Cells, row("first address field binary", part))

	fmt.Println(table.String())
}

// ip6SubnetDisplay describe a link local IP
func ip6SubnetDisplayBasic(addr netip.Addr, prefix netip.Prefix) {
	var value string
	var err error
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, row("IP Type", ipv6.AddrTypeName(addr)))
	table.Body.Cells = append(table.Body.Cells, row("Type Prefix", ipv6.AddrTypePrefix(addr).Masked()))
	table.Body.Cells = append(table.Body.Cells, row("IP", addr.String()))
	table.Body.Cells = append(table.Body.Cells, row("Prefix", prefix.Masked()))
	value, err = ipv6.MulticastNetworkPrefix(addr)
	if err != nil {
		fmt.Println(err)
	}
	table.Body.Cells = append(table.Body.Cells, row("Network Prefix", fmt.Sprintf("%s", value)))
	value, err = ipv6.MulticastGroupID(addr)
	if err != nil {
		fmt.Println(err)
	}
	table.Body.Cells = append(table.Body.Cells, row("Group ID", fmt.Sprintf("%s", value)))
	part := strings.Split(ipv6.AddrToBitString(addr), ".")[0]
	part = fmt.Sprintf("%s%s", strings.Repeat("0", 16-len(part)), part)
	table.Body.Cells = append(table.Body.Cells, row("first address field binary", part))

	fmt.Println(table.String())
}

// IP6RandomIPs produce list of random IPs
func IP6RandomIPs(ip6Type string, number int) {
	if number == 0 {
		number = 10
	}
	if number > 100 {
		number = 100
	}
	var addr netip.Addr
	var err error
	if ip6Type == typeGlobalUnicast {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandomAddrGlobalUnicast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeLinkLocal {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandomAddrLinkLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typePrivate {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandomAddrPrivate()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeMulticast {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandomAddrMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeInterfaceLocalMulticast {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandomAddrInterfaceLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeLinkLocalMulticast {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandomAddrLinkLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else {
		fmt.Println("No valid type specified")
		os.Exit(1)
	}
}
