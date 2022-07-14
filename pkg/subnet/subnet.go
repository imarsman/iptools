package subnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"net/netip"

	"github.com/imarsman/iptools/pkg/util"
	"gopkg.in/yaml.v2"
)

const (
	octetBits = 8
)

// IPV4Range added because ranges were removed when netaddr joined core Go
type IPV4Range struct {
	first netip.Addr
	last  netip.Addr
}

// NewIPV4Range make a new range
func NewIPV4Range(first, last netip.Addr) IPV4Range {
	var r = IPV4Range{}
	r.first = first
	r.last = last

	return r
}

// First get first address in range
func (r *IPV4Range) First() netip.Addr {
	return r.first
}

// Last get last address in range
func (r *IPV4Range) Last() netip.Addr {
	return r.last
}

// String return string version of range
func (r *IPV4Range) String() string {
	return fmt.Sprintf("%s-%s", r.first.String(), r.last.String())
}

// IPV4Subnet an IP subnet
type IPV4Subnet struct {
	name   string
	prefix netip.Prefix
}

// Name subnet name
func (s *IPV4Subnet) Name() string {
	return s.name
}

// SetName set subnet name
func (s *IPV4Subnet) SetName(name string) {
	s.name = name
}

// Prefix get prefix for subnet
func (s *IPV4Subnet) Prefix() netip.Prefix {
	return s.prefix
}

// IP get IP for subnet
func (s *IPV4Subnet) IP() netip.Addr {
	return s.prefix.Addr()
}

// String get string representing subnet (cidr notation)
func (s *IPV4Subnet) String() string {
	return s.prefix.String()
}

// JSON get JSON for subnet
func (s *IPV4Subnet) JSON() (bytes []byte, err error) {
	var prefix = s.prefix.String()
	bytes, err = json.MarshalIndent(&prefix, "", "  ")
	if err != nil {
		return
	}

	return bytes, nil
}

// YAML get YAML for subnet
func (s *IPV4Subnet) YAML() (bytes []byte, err error) {
	var prefix = s.prefix.String()
	bytes, err = yaml.Marshal(&prefix)
	if err != nil {
		return
	}

	return bytes, nil
}

// NewNamedFromIPAndBits new with name using incoming prefix ip and network bits
func NewNamedFromIPAndBits(ip string, bits uint8, name string) (subnet *IPV4Subnet, err error) {
	subnet, err = newSubnet(ip, uint8(bits))
	if err != nil {
		return
	}
	subnet.name = name

	return
}

// NewFromIPAndBits new using incoming prefix ip and network bits
func NewFromIPAndBits(ip string, bits uint8) (subnet *IPV4Subnet, err error) {
	return newSubnet(ip, uint8(bits))
}

// NewNamedFromPrefix new with name using incoming prefix
func NewNamedFromPrefix(prefix string, name string) (subnet *IPV4Subnet, err error) {
	p, err := netip.ParseAddr(prefix)
	if err != nil {
		return
	}

	subnet, err = newSubnet(p.String(), uint8(p.BitLen()))
	if err != nil {
		return
	}
	subnet.name = name

	return
}

// NewFromPrefix new using incoming prefix
func NewFromPrefix(prefix string) (subnet *IPV4Subnet, err error) {
	p, err := netip.ParsePrefix(prefix)
	if err != nil {
		return
	}

	subnet, err = newSubnet(p.Addr().String(), uint8(p.Bits()))
	if err != nil {
		return
	}

	return
}

// newSubnet new subnet with prefix ip and network bits
func newSubnet(ip string, bits uint8) (subnet *IPV4Subnet, err error) {
	errMsg := "invalid prefix"

	subnet = new(IPV4Subnet)

	var pfx netip.Prefix
	pfx, err = netip.ParsePrefix(fmt.Sprintf("%s/%d", ip, bits))
	if err != nil {
		return
	}
	pfx = pfx.Masked()

	if !pfx.IsValid() {
		return nil, errors.New(errMsg)
	}

	if pfx.Addr().Is6() {
		return nil, errors.New("subnet too large for current implementation")
	}
	subnet.prefix = pfx

	return
}

// BroadcastAddr get broadcast address for subnet, i.e. the max IP
func (s *IPV4Subnet) BroadcastAddr() (ip netip.Addr, err error) {
	return s.Last()
}

// CIDR get CIDR notation for subnet
func (s *IPV4Subnet) CIDR() (cidr string) {
	cidr = s.Prefix().Masked().String()

	return
}

// BinaryMask get dot delimited subnet mask in binary
func (s *IPV4Subnet) BinaryMask() (mask string) {
	mask = util.BitStr4(s.Prefix().Masked().Addr(), `.`)

	return
}

// BinaryID get the starting IP for subnet as binary
func (s *IPV4Subnet) BinaryID() (mask string) {
	mask = util.BitStr4(s.IP(), ``)

	return
}

// classOcted get octet for prefix IP
func (s *IPV4Subnet) classOctet() int {
	bits := s.Prefix().Masked().Addr().As4()

	if s.maxClassBits() == octetBits {
		return int(bits[0])
	} else if s.maxClassBits() == 2*octetBits {
		return int(bits[1])
	} else if s.maxClassBits() == 3*octetBits {
		return int(bits[2])
	}
	return int(bits[3])
}

// ClassNetworkBits bits not used for hosts in class block
func (s *IPV4Subnet) ClassNetworkBits() uint8 {
	return uint8(s.Prefix().Bits()) - s.startClassBits()
}

// ClassHostBits bits used for network in class block
func (s *IPV4Subnet) ClassHostBits() uint8 {
	return s.maxClassBits() - uint8(s.Prefix().Bits())
}

// TotalHosts total hosts in subnet
func (s *IPV4Subnet) TotalHosts() int64 {
	return s.Hosts() * s.Networks()
}

// Hosts bits remaining in mask block
func (s *IPV4Subnet) Hosts() int64 {
	if s.Prefix().Bits()%8 == 0 {
		// This may or may not be the proper solution
		if s.Prefix().Bits() == 32 {
			return 1
		}
		return int64(
			(math.Exp2(float64(32) - float64(s.Prefix().Bits()))) / float64(s.Networks()),
		)
	}
	return int64((math.Exp2(float64(32 - s.Prefix().Bits()))))
}

// UsableHosts number of usable hosts
func (s *IPV4Subnet) UsableHosts() int64 {
	if s.Hosts() < 2 {
		return 0
	}
	return s.Hosts() - 2
}

// startClassBits starting bits for subnet range for the class
func (s *IPV4Subnet) startClassBits() uint8 {
	if s.Prefix().Bits() < octetBits {
		return 0
	} else if s.Prefix().Bits() < 2*octetBits {
		return octetBits
	} else if s.Prefix().Bits() < 3*octetBits {
		return 2 * octetBits
	}
	return 3 * octetBits
}

// maxClassBits maximum bits for subnet range for the class
func (s *IPV4Subnet) maxClassBits() uint8 {
	if s.Prefix().Bits() <= octetBits {
		return octetBits
	} else if s.Prefix().Bits() <= 2*octetBits {
		return 2 * octetBits
	} else if s.Prefix().Bits() <= 3*octetBits {
		return 3 * octetBits
	}
	return 4 * octetBits
}

// Class get network class, a, b, or c
func (s *IPV4Subnet) Class() (class rune) {
	parts := s.IP().As4()
	bitStr := fmt.Sprintf("%08b", parts[0])

	// https://stackoverflow.com/a/34257287/2694971
	if strings.HasPrefix(bitStr, `0`) {
		return 'A'
	} else if strings.HasPrefix(bitStr, `10`) {
		return 'B'
	} else if strings.HasPrefix(bitStr, `110`) {
		return 'C'
	} else if strings.HasPrefix(bitStr, `1110`) {
		return 'D'
	} else if strings.HasPrefix(bitStr, `1111`) {
		return 'E'
	}

	return '0'
}

// First get first IP for subnet
func (s *IPV4Subnet) First() (ip netip.Addr, err error) {
	ip = s.prefix.Addr()
	return
}

// Last get last IP for subnet
func (s *IPV4Subnet) Last() (ip netip.Addr, err error) {
	ip, err = util.AddToIP(s.prefix.Addr(), int32(s.TotalHosts()-1))
	if err != nil {
		return
	}

	return
}

// NetworkAddr get last IP for subnet
func (s *IPV4Subnet) NetworkAddr() (ip netip.Addr, err error) {
	return s.Last()
}

// Networks number of subnets
func (s *IPV4Subnet) Networks() int64 {
	bits := uint8(s.Prefix().Bits()) - s.startClassBits()

	return int64(math.Exp2(float64(bits)))
}

// UsableIPs get usable ips for subnet
func (s *IPV4Subnet) UsableIPs() (ips []netip.Addr, err error) {
	errMsg := "empty ip list for subnet"
	ips, err = s.IPs()
	if err != nil {
		return
	}
	if len(ips) == 0 {
		err = errors.New(errMsg)
		ips = []netip.Addr{}
		return
	}
	ips = ips[1 : len(ips)-1]

	return
}

// IPs get ips for subnet
func (s *IPV4Subnet) IPs() (ips []netip.Addr, err error) {
	errMsg := "empty ip list for subnet"
	ip := s.Prefix().Addr()
	ips = append(ips, ip)

	for j := 0; j < int(s.TotalHosts()); j++ {
		ip = ip.Next()
		if (ip == netip.Addr{}) {
			err = errors.New(errMsg)
			ips = []netip.Addr{}
			return
		}
		ips = append(ips, ip)
	}

	return
}

// UsableIPRange get range of IPs usable for hosts
func (s *IPV4Subnet) UsableIPRange() (r IPV4Range, err error) {
	ip := s.Prefix().Addr()
	startIP := ip
	ip, err = util.AddToIP(ip, int32(s.TotalHosts()))
	if err != nil {
		return
	}
	r = NewIPV4Range(startIP.Next(), ip.Prev())

	return
}

// IPRange get subnet range
func (s *IPV4Subnet) IPRange() (r IPV4Range, err error) {
	ip := s.Prefix().Addr()
	startIP := ip
	ip, err = util.AddToIP(ip, int32(s.TotalHosts()))
	if err != nil {
		return
	}
	r = NewIPV4Range(startIP, ip)

	return
}

// SecondarySubnets set of subnets in the context of parent subnet
func (s *IPV4Subnet) SecondarySubnets(secondarySubnet *IPV4Subnet) (subnets []*IPV4Subnet, err error) {
	return s.subnets(secondarySubnet)
}

// Subnets the set of equally sized subnets for subnet
func (s *IPV4Subnet) Subnets() (subnets []*IPV4Subnet, err error) {
	return s.subnets(s)
}

// subnets split a subnet into smaller secondary subnets
func (s *IPV4Subnet) subnets(secondarySubnet *IPV4Subnet) (subnets []*IPV4Subnet, err error) {
	ranges, err := s.ipRanges(secondarySubnet)
	if err != nil {
		return
	}

	for _, r := range ranges {
		prefix := fmt.Sprintf("%s/%d", r.First(), secondarySubnet.Prefix().Bits())
		s, err = NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		subnets = append(subnets, s)
	}

	return
}

// SecondaryIPRanges set of ranges in the context of parent subnet
func (s *IPV4Subnet) SecondaryIPRanges(secondarySubnet *IPV4Subnet) (ranges []IPV4Range, err error) {
	return s.ipRanges(secondarySubnet)
}

// IPRanges the set of equally sized ranges for subnet
func (s *IPV4Subnet) IPRanges() (ranges []IPV4Range, err error) {
	return s.ipRanges(s)
}

// ipRanges get the ranges for a subnet splitting by secondary subnet (can be self)
func (s *IPV4Subnet) ipRanges(secondarySubnet *IPV4Subnet) (ranges []IPV4Range, err error) {
	// Can't subdivide to smaller prefixed subnet
	if secondarySubnet.Prefix().Bits() < s.Prefix().Bits() {
		err = fmt.Errorf("Subnet to split to has more bits %d than parent %d", s.Prefix().Bits(), secondarySubnet.Prefix().Bits())
		return
	}
	ranges = []IPV4Range{}
	ip := s.IP()
	ipStart := ip

	ratio := int(math.Exp2(float64(secondarySubnet.Prefix().Bits() - s.Prefix().Bits())))
	for j := 0; j < int(s.Networks()); j++ {
		for r := 0; r < ratio; r++ {
			ip, err := util.AddToIP(ipStart, int32(secondarySubnet.Hosts()-1))
			if err != nil {
				//return
			}
			ranges = append(ranges, NewIPV4Range(ipStart, ip))
			ipStart = ip.Next()
		}
	}

	return
}
