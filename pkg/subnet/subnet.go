package subnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/imarsman/iptools/pkg/util"
	"gopkg.in/yaml.v2"
	"inet.af/netaddr"
)

const (
	octets   = 4
	octetMax = 255
)

// IPV4Subnet an IP subnet
type IPV4Subnet struct {
	name   string
	prefix netaddr.IPPrefix
	ip     netaddr.IP
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
func (s *IPV4Subnet) Prefix() netaddr.IPPrefix {
	return s.prefix
}

// IP get IP for subnet
func (s *IPV4Subnet) IP() netaddr.IP {
	return s.ip
}

// String get string representing subnet (cidr notation)
func (s *IPV4Subnet) String() string {
	return s.prefix.String()
}

// JSON get JSON for subnet
func (s *IPV4Subnet) JSON() (bytes []byte, err error) {
	bytes, err = json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	return bytes, nil
}

// YAML get YAML for subnet
func (s *IPV4Subnet) YAML() (bytes []byte, err error) {
	bytes, err = yaml.Marshal(s)
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
	p, err := netaddr.ParseIPPrefix(prefix)
	if err != nil {
		return
	}

	subnet, err = newSubnet(p.IP().String(), p.Bits())
	if err != nil {
		return
	}
	subnet.name = name

	return
}

// NewFromPrefix new using incoming prefix
func NewFromPrefix(prefix string) (subnet *IPV4Subnet, err error) {
	p, err := netaddr.ParseIPPrefix(prefix)
	if err != nil {
		return
	}

	subnet, err = newSubnet(p.IP().String(), p.Bits())
	if err != nil {
		return
	}

	return
}

// newSubnet new subnet with prefix ip and network bits
func newSubnet(ip string, bits uint8) (subnet *IPV4Subnet, err error) {
	subnet = new(IPV4Subnet)

	errMsg := "invalid prefix"

	addressIP, err := netaddr.ParseIP(ip)
	if err != nil {
		return
	}

	var pfx netaddr.IPPrefix
	prefixStr := fmt.Sprintf("%s/%d", ip, bits)
	var pfxPre netaddr.IPPrefix
	pfxPre, err = netaddr.ParseIPPrefix(prefixStr)
	if err != nil {
		return
	}
	prefixStr = fmt.Sprintf("%s/%d", pfxPre.Masked().IP().String(), bits)
	pfx = pfxPre.Masked()

	if !pfx.IsValid() {
		return nil, errors.New(errMsg)
	}

	if pfx.IP().Is6() {
		return nil, errors.New("subnet too large for current implementation")
	}
	subnet.ip = addressIP
	subnet.prefix = pfx

	return
}

// BroadcastAddress get broadcast address for subnet, i.e. the max IP
func (s *IPV4Subnet) BroadcastAddress() (ip netaddr.IP, err error) {
	return s.Last()
}

// CIDR get CIDR notation for subnet
func (s *IPV4Subnet) CIDR() (cidr string) {
	cidr = s.Prefix().Masked().String()

	return
}

// BinarySubnetMask get dot delimited subnet mask in binary
func (s *IPV4Subnet) BinarySubnetMask() (mask string) {
	mask = util.BitStr4(s.Prefix().Masked().IP(), `.`)

	return
}

// BinaryID get the starting IP for subnet as binary
func (s *IPV4Subnet) BinaryID() (mask string) {
	mask = util.BitStr4(s.IP(), ``)

	return
}

// classOcted get octet for prefix IP
func (s *IPV4Subnet) classOctet() int {
	bits := s.Prefix().Masked().IP().As4()

	if s.maxBitsForClass() == 8 {
		return int(bits[0])
	} else if s.maxBitsForClass() == 16 {
		return int(bits[1])
	} else if s.maxBitsForClass() == 24 {
		return int(bits[2])
	}
	return int(bits[3])
}

// PrefixBits byte used for subnet
func (s *IPV4Subnet) PrefixBits() uint8 {
	return s.Prefix().Bits()
}

// ClassNetworkBits bits not used for hosts in class block
func (s *IPV4Subnet) ClassNetworkBits() uint8 {
	return s.Prefix().Bits() - s.startBitsForClass()
}

// ClassHostBits bits used for network in class block
func (s *IPV4Subnet) ClassHostBits() uint8 {
	return s.maxBitsForClass() - s.Prefix().Bits()
}

// TotalHosts total hosts in subnet
func (s *IPV4Subnet) TotalHosts() int64 {
	return s.Hosts() * s.Networks()
}

// Hosts bits remaining in mask block
func (s *IPV4Subnet) Hosts() int64 {
	if s.Prefix().Bits()%8 == 0 {
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

// startBitsForClass starting bits for subnet range for the class
func (s *IPV4Subnet) startBitsForClass() uint8 {
	if s.Prefix().Bits() < 8 {
		return 0
	} else if s.Prefix().Bits() < 16 {
		return 8
	} else if s.Prefix().Bits() < 24 {
		return 16
	}
	return 24
}

// maxBitsForClass maximum bits for subnet range for the class
func (s *IPV4Subnet) maxBitsForClass() uint8 {
	if s.Prefix().Bits() <= 8 {
		return 8
	} else if s.Prefix().Bits() <= 16 {
		return 16
	} else if s.Prefix().Bits() <= 24 {
		return 24
	}
	return 32
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
func (s *IPV4Subnet) First() (ip netaddr.IP, err error) {
	ip = s.IP()
	return
}

// Last get last IP for subnet
func (s *IPV4Subnet) Last() (ip netaddr.IP, err error) {
	ip, err = util.AddToIP(s.IP(), int32(s.TotalHosts()-1))
	if err != nil {
		return
	}

	return
}

// NetworkAddress get last IP for subnet
func (s *IPV4Subnet) NetworkAddress() (ip netaddr.IP, err error) {
	return s.Last()
}

// Networks number of subnets
func (s *IPV4Subnet) Networks() int64 {
	bits := s.Prefix().Bits() - s.startBitsForClass()

	return int64(math.Exp2(float64(bits)))
}

// UsableIPs get usable ips for subnet
func (s *IPV4Subnet) UsableIPs() (ips []netaddr.IP, err error) {
	errMsg := "empty ip list for subnet"
	ips, err = s.IPs()
	if err != nil {
		return
	}
	if len(ips) == 0 {
		err = errors.New(errMsg)
		ips = []netaddr.IP{}
		return
	}
	ips = ips[1 : len(ips)-1]

	return
}

// IPs get ips for subnet
func (s *IPV4Subnet) IPs() (ips []netaddr.IP, err error) {
	errMsg := "empty ip list for subnet"
	ip := s.Prefix().IP()
	ips = append(ips, ip)

	for j := 0; j < int(s.TotalHosts()); j++ {
		ip = ip.Next()
		if (ip == netaddr.IP{}) {
			err = errors.New(errMsg)
			ips = []netaddr.IP{}
			return
		}
		ips = append(ips, ip)
	}

	return
}

// UsableIPRange get range of IPs usable for hosts
func (s *IPV4Subnet) UsableIPRange() (r netaddr.IPRange, err error) {
	ip := s.Prefix().IP()
	startIP := ip
	ip, err = util.AddToIP(ip, int32(s.TotalHosts()))
	if err != nil {
		return
	}
	r = netaddr.IPRangeFrom(startIP.Next(), ip.Prior())

	return
}

// IPRange get subnet range
func (s *IPV4Subnet) IPRange() (r netaddr.IPRange, err error) {
	ip := s.Prefix().IP()
	startIP := ip
	ip, err = util.AddToIP(ip, int32(s.TotalHosts()))
	if err != nil {
		return
	}
	r = netaddr.IPRangeFrom(startIP, ip)

	return
}

// NetworkSubnetsInSubnet set of subnets in the context of parent subnet
func (s *IPV4Subnet) NetworkSubnetsInSubnet(childSubnet *IPV4Subnet) (subnets []*IPV4Subnet, err error) {
	return s.networkSubnets(childSubnet)
}

// NetworkSubnets the set of equally sized subnets for subnet
func (s *IPV4Subnet) NetworkSubnets() (subnets []*IPV4Subnet, err error) {
	return s.networkSubnets(s)
}

// networkSubnets split a subnet into smaller child subnets
func (s *IPV4Subnet) networkSubnets(childSubnet *IPV4Subnet) (subnets []*IPV4Subnet, err error) {
	ranges, err := s.networkIPRanges(childSubnet)
	if err != nil {
		return
	}

	for _, r := range ranges {
		prefix := fmt.Sprintf("%s/%d", r.From(), childSubnet.Prefix().Bits())
		s, err = NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		subnets = append(subnets, s)
	}

	return
}

// NetworkIPRangesInSubnet set of ranges in the context of parent subnet
func (s *IPV4Subnet) NetworkIPRangesInSubnet(childSubnet *IPV4Subnet) (ranges []netaddr.IPRange, err error) {
	return s.networkIPRanges(childSubnet)
}

// NetworkIPRanges the set of equally sized ranges for subnet
func (s *IPV4Subnet) NetworkIPRanges() (ranges []netaddr.IPRange, err error) {
	return s.networkIPRanges(s)
}

// networkIPRanges get the ranges for a subnet splitting by child subnet (can be self)
func (s *IPV4Subnet) networkIPRanges(childSubnet *IPV4Subnet) (ranges []netaddr.IPRange, err error) {
	// Can't subdivide to smaller prefixed subnet
	if childSubnet.Prefix().Bits() < s.Prefix().Bits() {
		err = fmt.Errorf("Subnet to split to has more bits %d than parent %d", s.Prefix().Bits(), childSubnet.Prefix().Bits())
		return
	}
	ranges = []netaddr.IPRange{}
	ip := s.IP()
	ipStart := ip

	ratio := int(math.Exp2(float64(childSubnet.Prefix().Bits() - s.Prefix().Bits())))
	for j := 0; j < int(s.Networks()); j++ {
		for r := 0; r < ratio; r++ {
			ip, err := util.AddToIP(ipStart, int32(childSubnet.Hosts()-1))
			if err != nil {
				//return
			}
			ranges = append(ranges, netaddr.IPRangeFrom(ipStart, ip))
			ipStart = ip.Next()
		}
	}

	return
}
