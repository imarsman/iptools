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
	ip4util "github.com/imarsman/iptools/pkg/ipv4subnet/ip4util"

	"github.com/imarsman/iptools/pkg/ipv6subnet"
	ip6util "github.com/imarsman/iptools/pkg/ipv6subnet/ip6util"
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
	table.Body.Cells = append(table.Body.Cells, row("Broadcast Address Hex ID", ip4util.IPToHexStr(s.Last())))
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

	table.Body.Cells = append(table.Body.Cells, row("in-addr.arpa", ip4util.InAddrArpa(s.Prefix().Addr())))

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
			addr, err = ip6util.RandomAddrGlobalUnicast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeLinkLocal {
			addr, err = ip6util.RandomAddrLinkLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeUniqueLocal {
			addr, err = ip6util.RandomAddrUniqueLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeMulticast {
			addr, err = ip6util.RandomAddrMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeLinkLocalMulticast {
			addr, err = ip6util.RandomAddrLinkLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if ip6Type == typeInterfaceLocalMulticast {
			addr, err = ip6util.RandomAddrInterfaceLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("No valid type specified")
			os.Exit(1)
		}
	}
	s, err := ipv6subnet.NewFromIPAndBits(addr.StringExpanded(), bits)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !(s.Addr().IsMulticast() || s.Addr().IsInterfaceLocalMulticast() || s.Addr().IsLinkLocalMulticast()) {
		ip6SubnetDisplay(s)
	} else {
		ip6SubnetDisplayBasic(s)
	}
}

// ip6SubnetDisplay describe a link local IP
func ip6SubnetDisplay(s *ipv6subnet.Subnet) {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, row("IP Type", ip6util.AddressTypeName(s.Addr())))
	table.Body.Cells = append(table.Body.Cells, row("Type Prefix", s.TypePrefix().Masked()))
	table.Body.Cells = append(table.Body.Cells, row("IP", s.Addr().String()))
	if ip6util.HasType(ip6util.AddressType(s.Addr()), ip6util.GlobalUnicast, ip6util.LinkLocalUnicast, ip6util.UniqueLocal) {
		solicitedNodeAddr, err := ip6util.SolicitedNodeMulticast(s.Addr())
		if err != nil {
			panic(err)
		}
		table.Body.Cells = append(table.Body.Cells,
			row(
				"Solicited node multicast", solicitedNodeAddr.String(),
			),
		)
	}

	table.Body.Cells = append(table.Body.Cells, row("Prefix", s.Prefix().Masked()))
	if ip6util.AddressType(s.Addr()) == ip6util.GlobalUnicast {
		table.Body.Cells = append(table.Body.Cells, row("Routing Prefix", fmt.Sprintf("%s", s.RoutingPrefix())))
	}
	if ip6util.HasType(ip6util.AddressType(s.Addr()), ip6util.GlobalUnicast, ip6util.UniqueLocal) {
		table.Body.Cells = append(table.Body.Cells, row("Global ID", fmt.Sprintf("%s", ip6util.GlobalID(s.Addr()))))
	}
	table.Body.Cells = append(table.Body.Cells, row("Interface ID", fmt.Sprintf("%s", s.InterfaceString())))
	table.Body.Cells = append(table.Body.Cells, row("Subnet ID", fmt.Sprintf("%s", s.SubnetString())))
	if ip6util.AddressType(s.Addr()) == ip6util.LinkLocalUnicast {
		table.Body.Cells = append(table.Body.Cells, row("Default Gateway", s.LinkLocalDefaultGateway()))
	}
	table.Body.Cells = append(table.Body.Cells, row("Link", s.Link()))
	if s.IsARPA() {
		table.Body.Cells = append(table.Body.Cells, row("ip6.arpa", fmt.Sprintf("%s", ip6util.IP6Arpa(s.Addr()))))
	}
	table.Body.Cells = append(table.Body.Cells, row("Subnet first address", s.First().StringExpanded()))
	table.Body.Cells = append(table.Body.Cells, row("Subnet last address", s.Last().StringExpanded()))
	part := strings.Split(ip6util.AddrToBitString(s.Addr()), ".")[0]
	part = fmt.Sprintf("%s%s", strings.Repeat("0", 16-len(part)), part)
	table.Body.Cells = append(table.Body.Cells, row("first address field binary", part))

	fmt.Println(table.String())
}

// ip6SubnetDisplay describe a link local IP
func ip6SubnetDisplayBasic(s *ipv6subnet.Subnet) {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Category"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, row("IP Type", ip6util.AddressTypeName(s.Addr())))
	table.Body.Cells = append(table.Body.Cells, row("Type Prefix", s.TypePrefix().Masked()))
	table.Body.Cells = append(table.Body.Cells, row("IP", s.Addr().String()))
	table.Body.Cells = append(table.Body.Cells, row("Prefix", s.Prefix().Masked()))
	table.Body.Cells = append(table.Body.Cells, row("Network Prefix", fmt.Sprintf("%s", ip6util.MulticastNetworkPrefix(s.Addr()))))
	table.Body.Cells = append(table.Body.Cells, row("Group ID", fmt.Sprintf("%s", ip6util.MulticastGroupID(s.Addr()))))
	part := strings.Split(ip6util.AddrToBitString(s.Addr()), ".")[0]
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
			addr, err = ip6util.RandomAddrGlobalUnicast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeLinkLocal {
		for i := 0; i < number; i++ {
			addr, err = ip6util.RandomAddrLinkLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeUniqueLocal {
		for i := 0; i < number; i++ {
			addr, err = ip6util.RandomAddrUniqueLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeMulticast {
		for i := 0; i < number; i++ {
			addr, err = ip6util.RandomAddrMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeInterfaceLocalMulticast {
		for i := 0; i < number; i++ {
			addr, err = ip6util.RandomAddrInterfaceLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == typeLinkLocalMulticast {
		for i := 0; i < number; i++ {
			addr, err = ip6util.RandomAddrLinkLocalMulticast()
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
