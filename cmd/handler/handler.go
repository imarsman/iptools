package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/netip"
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"

	"github.com/alexeyco/simpletable"
	"github.com/imarsman/iptools/cmd/args"

	"github.com/imarsman/iptools/pkg/ipv4subnet"
	"github.com/imarsman/iptools/pkg/ipv4subnet/ipv4util"
	"github.com/imarsman/iptools/pkg/ipv6"
	"github.com/imarsman/iptools/pkg/util"
)

var printer = message.NewPrinter(language.English)

func row(label string, value any) (r []*simpletable.Cell) {
	r = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%v", label)},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%v", value)},
	}
	return
}

// LookupDomain look up IPs for a domain
func LookupDomain(domains []string, mxLookup bool, toJSON, toYAML bool) {
	ipsForDomains := ipv6.NewDomainInfoSet()

	table := simpletable.New()
	for i, domain := range domains {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Type/MX"},
				{Align: simpletable.AlignCenter, Text: "Address/Pref"},
			},
		}
		domainInfo := ipv6.NewDomainInfo()
		domainInfo.Domain = domain

		domainRow := []*simpletable.Cell{
			{},
			{Align: simpletable.AlignLeft, Text: domain},
		}
		table.Body.Cells = append(table.Body.Cells, domainRow)
		addresses, err := util.DomainAddresses(domain)
		if err != nil {
			domainRow = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Span: 2, Text: err.Error()},
			}
			table.Body.Cells = append(table.Body.Cells, domainRow)
		}
		for _, addr := range addresses {
			var addressType string = "ipv4"
			if addr.Is6() {
				addressType = "ipv6"
			}
			addressInfo := ipv6.AddressInfo{}
			addressInfo.Type = addressType
			addressInfo.Address = addr.String()
			domainInfo.AddressInfo = append(domainInfo.AddressInfo, addressInfo)
			table.Body.Cells = append(table.Body.Cells, row(addressType, addr.String()))
		}

		if mxLookup {
			mxRecods, err := util.DomainMXRecods(domain)
			if len(mxRecods) > 0 && err == nil {
				mxRecordRow := []*simpletable.Cell{
					{},
					{Align: simpletable.AlignLeft, Text: "MX Records"},
				}
				table.Body.Cells = append(table.Body.Cells, mxRecordRow)
			}
			for _, mxR := range mxRecods {
				mxRecord := ipv6.MXRecordInfo{}
				mxRecord.Domain = mxR.Host
				mxRecord.Pref = mxR.Pref
				domainInfo.MXRecords = append(domainInfo.MXRecords, mxRecord)

				table.Body.Cells = append(table.Body.Cells, row(mxR.Host, fmt.Sprintf("%d", mxR.Pref)))
			}
		}
		ipsForDomains.DomainInfo = append(ipsForDomains.DomainInfo, domainInfo)
		if i+1 < len(domains) {
			r := []*simpletable.Cell{
				{},
				{},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}
	}
	table.SetStyle(simpletable.StyleCompactLite)

	if toJSON {
		bytes, err := json.MarshalIndent(&ipsForDomains, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else if toYAML {
		bytes, err := yaml.Marshal(&ipsForDomains)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(table.String())
	}
}

// IP4SubnetDescribe describe a subnet
// Needs review and cleanup
// Investigate iptools subnetip4 describe -ip 10.32.0.0 -bits 23 -secondary-bits 24
func IP4SubnetDescribe(ip string, bits int, secondaryBits int) {
	var err error
	var s *ipv4subnet.Subnet

	// Default to 24 bits
	if bits == 0 {
		bits = 24
	}

	var prefixStr string
	prefix, err := netip.ParsePrefix(ip)
	if err == nil {
		bits = prefix.Bits()
		prefixStr = prefix.String()
	} else {
		prefixStr = fmt.Sprintf("%s/%d", ip, bits)
		prefix, err = netip.ParsePrefix(prefixStr)
		if err != nil {

		}
	}

	s, err = ipv4subnet.NewFromPrefix(prefixStr)
	if err != nil {
		os.Exit(1)
	}
	var s2 *ipv4subnet.Subnet
	if secondaryBits != 0 {
		// need to redefine ip since it would have primary bits already.
		ip = prefix.Addr().String()
		prefixStr := fmt.Sprintf("%s/%d", ip, secondaryBits)
		s2, err = ipv4subnet.NewFromPrefix(prefixStr)
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

	table.Body.Cells = append(table.Body.Cells, row("IP Type", util.AddrTypeName(prefix.Addr())))
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

	table.Body.Cells = append(table.Body.Cells, row("in-addr.arpa", ipv4util.Arpa(s.Prefix().Addr())))
	// table.Body.Cells = append(table.Body.Cells, row("in-addr.arpa", ipv4util.InAddrArpa(s.Prefix().Addr())))

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

	prefixStr := fmt.Sprintf("%s/%d", ip, bits)

	prefix, err := netip.ParsePrefix(ip)
	if err == nil {
		bits = prefix.Bits()
		prefixStr = prefix.String()
	} else {
		prefixStr = fmt.Sprintf("%s/%d", ip, bits)
	}

	s, err = ipv4subnet.NewFromPrefix(prefixStr)
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

	var prefixStr string
	prefix, err := netip.ParsePrefix(ip)
	if err == nil {
		bits = prefix.Bits()
		prefixStr = prefix.String()
	} else {
		prefixStr = fmt.Sprintf("%s/%d", ip, bits)
	}

	s, err = ipv4subnet.NewFromPrefix(prefixStr)
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

// IP6SubnetDescribe describe a link-local address
func IP6SubnetDescribe(ip string, bits int, random bool, ip6Type string, json, yaml bool) {
	if ip6Type == "" && ip == "" {
		fmt.Println("If no IP then type must be supplied")
		os.Exit(1)
	} else if ip == "" && !random {
		fmt.Println("No type supplied and -random is false")
		os.Exit(1)
	}

	var addr netip.Addr

	prefix, err := netip.ParsePrefix(ip)
	if err == nil {
		addr = prefix.Addr()
		bits = prefix.Bits()
	} else {
		if bits == 0 {
			bits = 64
		}
		if !random {
			var err error
			addr, err = netip.ParseAddr(ip)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			var err error
			if ip6Type == ipv6.GlobalUnicastName {
				addr, err = ipv6.RandAddrGlobalUnicast()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if ip6Type == ipv6.LinkLocalName {
				addr, err = ipv6.RandAddrLinkLocal()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if ip6Type == ipv6.PrivateName {
				addr, err = ipv6.RandAddrPrivate()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if ip6Type == ipv6.MulticastName {
				addr, err = ipv6.RandAddrMulticast()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if ip6Type == ipv6.LinkLocalMulticastName {
				addr, err = ipv6.RandAddrLinkLocalMulticast()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if ip6Type == ipv6.InterfaceLocalMulticastName {
				addr, err = ipv6.RandAddrInterfaceLocalMulticast()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Println("No valid type specified")
				os.Exit(1)
			}
		}
	}

	if !(addr.IsMulticast() || addr.IsInterfaceLocalMulticast() || addr.IsLinkLocalMulticast()) {
		var prefix netip.Prefix
		if bits != 0 {
			prefix = netip.PrefixFrom(addr, bits)
		}
		ip6SubnetDisplay(addr, prefix, json, yaml)
	} else {
		var prefix netip.Prefix
		if bits != 0 {
			prefix = netip.PrefixFrom(addr, bits)
		}
		ip6SubnetDisplayBasic(addr, prefix, json, yaml)
	}
}

// ip6SubnetDisplay describe a link local IP
func ip6SubnetDisplay(addr netip.Addr, prefix netip.Prefix, toJSON, toYAML bool) {
	var ipSummary = ipv6.IPSummary{}

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

	ipSummary.IPType = util.AddrTypeName(addr)
	fmt.Println("addr", addr.String())
	table.Body.Cells = append(table.Body.Cells, row("IP Type", util.AddrTypeName(addr)))
	ipSummary.TypePrefix = ipv6.AddrTypePrefix(addr).Masked().String()
	table.Body.Cells = append(table.Body.Cells, row("Type Prefix", ipv6.AddrTypePrefix(addr).Masked()))
	ipSummary.IP = addr.String()
	table.Body.Cells = append(table.Body.Cells, row("IP", addr.String()))
	if ipv6.HasType(util.AddrType(addr), ipv6.GlobalUnicast, ipv6.LinkLocalUnicast, ipv6.UniqueLocal, ipv6.Private) {
		solicitedNodeAddr, err := ipv6.AddrSolicitedNodeMulticast(addr)
		if err != nil {
			panic(err)
		}
		ipSummary.SolicitedNodeMulticast = solicitedNodeAddr.String()
		table.Body.Cells = append(table.Body.Cells,
			row(
				"Solicited node multicast", solicitedNodeAddr.String(),
			),
		)
	}

	if (prefix != netip.Prefix{}) {
		ipSummary.Prefix = prefix.Masked().String()
		table.Body.Cells = append(table.Body.Cells, row("Prefix", prefix.Masked()))
	}
	if util.AddrType(addr) == ipv6.GlobalUnicast {
		ipSummary.RoutingPrefix = ipv6.RoutingPrefix(addr)
		table.Body.Cells = append(
			table.Body.Cells, row(
				"Routing Prefix", fmt.Sprintf("%s", fmt.Sprintf("%s", ipv6.RoutingPrefix(addr)))),
		)
	}
	ipSummary.SubnetID = ipv6.AddrSubnet(addr)
	table.Body.Cells = append(table.Body.Cells, row("Subnet ID", fmt.Sprintf("%s", ipv6.AddrSubnet(addr))))
	if ipv6.HasType(util.AddrType(addr), ipv6.GlobalUnicast, ipv6.Private, ipv6.LinkLocalUnicast) {
		number := printer.Sprintf("%.0f", math.Exp2(16))
		ipSummary.Subnets = int64(math.Exp2(16))
		table.Body.Cells = append(table.Body.Cells, row("Subnets", number))
	}
	// Handle global id for appropriate types
	if ipv6.HasType(util.AddrType(addr), ipv6.GlobalUnicast, ipv6.Private) {
		value, err = ipv6.AddrGlobalID(addr)
		if err != nil {
			fmt.Println(err)
		}
		ipSummary.GlobalID = value
		table.Body.Cells = append(table.Body.Cells, row("Global ID", fmt.Sprintf("%s", value)))
	}
	ipSummary.InterfaceID = ipv6.Interface(addr)
	table.Body.Cells = append(table.Body.Cells, row("Interface ID", fmt.Sprintf("%s", ipv6.Interface(addr))))
	if ipv6.HasType(util.AddrType(addr), ipv6.GlobalUnicast, ipv6.Private, ipv6.LinkLocalUnicast) {
		number := printer.Sprintf("%.0f", math.Exp2(64))
		ipSummary.Addresses = int64(math.Exp2(64))
		table.Body.Cells = append(table.Body.Cells, row("Addresses", number))
	}
	if util.AddrType(addr) == ipv6.LinkLocalUnicast {
		ipSummary.DefaultGateway = ipv6.LinkLocalDefaultGateway(addr)
		table.Body.Cells = append(table.Body.Cells, row("Default Gateway", ipv6.LinkLocalDefaultGateway(addr)))
	}
	if ipv6.HasType(util.AddrType(addr), ipv6.GlobalUnicast) {
		ipSummary.Link = ipv6.AddrLink(addr)
		table.Body.Cells = append(table.Body.Cells, row("Link", ipv6.AddrLink(addr)))
	}
	if ipv6.HasType(util.AddrType(addr), ipv6.GlobalUnicast) {
		ipSummary.IPV6Arpa = ipv6.Arpa(addr)
		table.Body.Cells = append(table.Body.Cells, row("ip6.arpa", fmt.Sprintf("%s", ipv6.Arpa(addr))))
	}
	ipSummary.SubnetFirstAddress = ipv6.First(addr).StringExpanded()
	table.Body.Cells = append(table.Body.Cells, row("Subnet first address", ipv6.First(addr).StringExpanded()))
	ipSummary.SubnetLastAddress = ipv6.Last(addr).StringExpanded()
	table.Body.Cells = append(table.Body.Cells, row("Subnet last address", ipv6.Last(addr).StringExpanded()))
	part := strings.Split(ipv6.Addr2BitString(addr), ".")[0]
	part = fmt.Sprintf("%s%s", strings.Repeat("0", 16-len(part)), part)
	ipSummary.FirstAddressFieldBinary = part
	table.Body.Cells = append(table.Body.Cells, row("1st address field binary", part))

	if toJSON {
		bytes, err := json.MarshalIndent(&ipSummary, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else if toYAML {
		bytes, err := yaml.Marshal(&ipSummary)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(table.String())
	}

}

// ip6SubnetDisplay describe a link local IP
func ip6SubnetDisplayBasic(addr netip.Addr, prefix netip.Prefix, toJSON, toYAML bool) {
	var ipSummary = ipv6.IPSummary{}

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

	ipSummary.IPType = util.AddrTypeName(addr)
	table.Body.Cells = append(table.Body.Cells, row("IP Type", util.AddrTypeName(addr)))
	ipSummary.TypePrefix = ipv6.AddrTypePrefix(addr).Masked().String()
	table.Body.Cells = append(table.Body.Cells, row("Type Prefix", ipv6.AddrTypePrefix(addr).Masked()))
	ipSummary.IP = addr.String()
	table.Body.Cells = append(table.Body.Cells, row("IP", addr.String()))
	if (prefix != netip.Prefix{}) {
		table.Body.Cells = append(table.Body.Cells, row("Prefix", prefix.Masked()))
	}
	value, err = ipv6.AddrMulticastNetworkPrefix(addr)
	if err != nil {
		fmt.Println(err)
	}
	ipSummary.NetworkPrefix = fmt.Sprintf("%s", value)
	table.Body.Cells = append(table.Body.Cells, row("Network Prefix", fmt.Sprintf("%s", value)))
	value, err = ipv6.AddrMulticastGroupID(addr)
	if err != nil {
		fmt.Println(err)
	}
	ipSummary.GroupID = fmt.Sprintf("%s", value)
	table.Body.Cells = append(table.Body.Cells, row("Group ID", fmt.Sprintf("%s", value)))
	if ipv6.HasType(util.AddrType(addr), ipv6.Multicast, ipv6.LinkLocalMulticast, ipv6.InterfaceLocalMulticast) {
		number := printer.Sprintf("%.0f", math.Exp2(32))
		ipSummary.Groups = int64(math.Exp2(32))
		table.Body.Cells = append(table.Body.Cells, row("Groups", number))
	}
	part := strings.Split(ipv6.Addr2BitString(addr), ".")[0]
	part = fmt.Sprintf("%s%s", strings.Repeat("0", 16-len(part)), part)
	ipSummary.FirstAddressFieldBinary = part

	table.Body.Cells = append(table.Body.Cells, row("first address field binary", part))

	if toJSON {
		bytes, err := json.MarshalIndent(&ipSummary, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else if toYAML {
		bytes, err := yaml.Marshal(&ipSummary)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(table.String())
	}
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
	if ip6Type == ipv6.GlobalUnicastName {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandAddrGlobalUnicast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == ipv6.LinkLocalName {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandAddrLinkLocal()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == ipv6.PrivateName {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandAddrPrivate()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == ipv6.MulticastName {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandAddrMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == ipv6.InterfaceLocalMulticastName {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandAddrInterfaceLocalMulticast()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(addr.StringExpanded())
		}
	} else if ip6Type == ipv6.LinkLocalMulticastName {
		for i := 0; i < number; i++ {
			addr, err = ipv6.RandAddrLinkLocalMulticast()
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
