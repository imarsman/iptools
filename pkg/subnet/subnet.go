package subnet

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"inet.af/netaddr"
)

const (
	octets   = 4
	octetMax = 255
)

// IPV4Subnet an IP subnet
// https://www.calculator.net/ip-IPV4Subnet-calculator.html?cclass=any&csubnet=16&cip=10.0.0.0&ctype=ipv4&printit=0&x=43&y=21
// https://www.calculator.net/ip-IPV4Subnet-calculator.html
type IPV4Subnet struct {
	Name   string
	Prefix netaddr.IPPrefix `json:"prefix" yaml:"prefix"`
	IP     netaddr.IP
}

// NewDefaultFromMask parse for 255.255.255.255 starting mask
// func NewDefaultFromMask(mask uint8) (subnet *IPV4Subnet, err error) {
// 	return newSubnet("255.255.255.255", mask, true)
// }

// NewFromMask new subnet with prefix 255.255.255.255 from incoming mask
func NewFromMask(mask uint8) (subnet *IPV4Subnet, err error) {
	return newSubnet("255.255.255.255", mask, true)
}

// NewDefaultFromPrefix new default (255.255.255.255) from prefix
// func NewDefaultFromPrefix(prefix string) (subnet *IPV4Subnet, err error) {
// 	errMsg := "invalid prefix"

// 	parts := strings.Split(prefix, "/")
// 	if len(parts) != 2 {
// 		err = errors.New(errMsg)
// 	}
// 	mask, err := strconv.Atoi(parts[1])
// 	if err != nil {
// 		return
// 	}
// 	return newSubnet(parts[0], uint8(mask), true)
// }

// NewFromIPAndMask new using incoming prefix ip and mask
func NewFromIPAndMask(ip string, mask uint8) (subnet *IPV4Subnet, err error) {
	return newSubnet(ip, uint8(mask), true)
}

// NewFromPrefix new using incoming prefix
func NewFromPrefix(prefix string) (subnet *IPV4Subnet, err error) {
	errMsg := "invalid prefix"

	parts := strings.Split(prefix, "/")
	if len(parts) != 2 {
		err = errors.New(errMsg)
	}
	mask, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	return newSubnet(parts[0], uint8(mask), true)
}

// newSubnet new subnet with prefix string, masked and boolean flag
func newSubnet(ip string, mask uint8, usemask bool) (subnet *IPV4Subnet, err error) {
	errMsg := "invalid prefix"

	subnet = new(IPV4Subnet)
	addressIP, err := netaddr.ParseIP(ip)
	if err != nil {
		return
	}
	subnet.IP = addressIP

	subnetAddress := "255.255.255.0"
	var pfx netaddr.IPPrefix
	// if usemask {
	prefixStr := fmt.Sprintf("%s/%d", subnetAddress, mask)
	var pfxPre netaddr.IPPrefix
	pfxPre, err = netaddr.ParseIPPrefix(prefixStr)
	if err != nil {
		return
	}
	prefixStr = fmt.Sprintf("%s/%d", pfxPre.Masked().IP().String(), mask)
	pfx = pfxPre.Masked()

	if !pfx.IsValid() {
		return nil, errors.New(errMsg)
	}

	if pfx.IP().Is6() {
		return nil, errors.New("subnet too large for current implementation")
	}

	subnet.Prefix = pfx

	return subnet, nil
}

func (s *IPV4Subnet) classMask() int {
	bits := s.Prefix.Masked().IP().As4()

	if s.maxBitsForClass() == 8 {
		return int(bits[0])
	} else if s.maxBitsForClass() == 16 {
		return int(bits[1])
	} else if s.maxBitsForClass() == 24 {
		return int(bits[2])
	}
	return int(bits[3])
}

// UsableRange omit first and last IPs
func UsableRange() (usableRange netaddr.IPRange) {

	return
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

// PrefixBits byte used for subnet
func (s *IPV4Subnet) PrefixBits() uint8 {
	return s.Prefix.Bits()
}

// ClassHostItentifierBits bits not used in mask block
func (s *IPV4Subnet) ClassHostItentifierBits() uint8 {
	return s.Prefix.Bits() - s.startBitsForClass()
}

// ClassNetworkPrefixBits bits used in mask block
func (s *IPV4Subnet) ClassNetworkPrefixBits() uint8 {
	return s.maxBitsForClass() - s.Prefix.Bits()
}

// TotalHosts total hosts in subnet
func (s *IPV4Subnet) TotalHosts() int64 {
	return s.NetworkHosts() * s.Networks()
}

// NetworkHosts bits remaining in mask block
func (s *IPV4Subnet) NetworkHosts() int64 {
	if s.Prefix.Bits()%8 == 0 {
		return int64((math.Exp2(float64(32) - float64(s.Prefix.Bits()))) / float64(s.Networks()))
	}
	return int64((math.Exp2(float64(32 - s.Prefix.Bits()))))
}

// UsableNetworkHosts number of usable hosts
func (s *IPV4Subnet) UsableNetworkHosts() int64 {
	if s.NetworkHosts() < 2 {
		return 0
	}
	return s.NetworkHosts() - 2
}

// maxBitsForClass maximum bits for subnet range for the class
func (s *IPV4Subnet) startBitsForClass() uint8 {
	if s.Prefix.Bits() < 8 {
		return 0
	} else if s.Prefix.Bits() < 16 {
		return 8
	} else if s.Prefix.Bits() < 24 {
		return 16
	}
	return 24
}

// maxBitsForClass maximum bits for subnet range for the class
func (s *IPV4Subnet) maxBitsForClass() uint8 {
	if s.Prefix.Bits() <= 8 {
		return 8
	} else if s.Prefix.Bits() <= 16 {
		return 16
	} else if s.Prefix.Bits() <= 24 {
		return 24
	}
	return 32
}

// Class get network class, a, b, or c
func (s *IPV4Subnet) Class() (class rune) {
	switch s.maxBitsForClass() {
	case 8:
		return 'A'
	case 16:
		return 'B'
	case 24:
		return 'C'
	default:
		return '0'
	}
}

// NetworkAddress get last IP for subnet
func (s *IPV4Subnet) NetworkAddress() (ip netaddr.IP, err error) {
	return s.Last()
}

// First get first IP for subnet
func (s *IPV4Subnet) First() (ip netaddr.IP, err error) {
	ip = s.IP
	return
}

// Last get last IP for subnet
func (s *IPV4Subnet) Last() (ip netaddr.IP, err error) {
	ip, err = addToIP(s.IP, int32(s.TotalHosts()-1))
	if err != nil {
		return
	}

	return
}

// func (s *IPV4Subnet) IPAddress() (ip netaddr.IP, err error) {
// 	ip, err = addToIP(s.IP, int32(s.TotalHosts()-1))
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// Networks number of subnets
func (s *IPV4Subnet) Networks() int64 {
	bits := s.Prefix.Bits() - s.startBitsForClass()

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
	ip := s.Prefix.IP()
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

// https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func int2ip(ipInt uint32) (netaddr.IP, bool) {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	addr, ok := netaddr.FromStdIP(ip)

	return addr, ok
}

func addToIP(startIP netaddr.IP, add int32) (newIP netaddr.IP, err error) {
	if !startIP.Next().IsValid() {
		err = fmt.Errorf("ip %v is already max", startIP)
		return
	}
	bytes := startIP.As4()
	slice := bytes[:]
	ipValue := binary.BigEndian.Uint32(slice)
	ipValue += uint32(add)
	newIP, ok := int2ip(ipValue)
	if !ok {
		err = fmt.Errorf("problem after adding %d to IP %v", add, startIP)
		return
	}
	return newIP, nil
}

// UsableIPRange get range of IPs usable for hosts
func (s *IPV4Subnet) UsableIPRange() (r netaddr.IPRange, err error) {
	ip := s.Prefix.IP()
	startIP := ip
	ip, err = addToIP(ip, int32(s.TotalHosts()))
	if err != nil {
		return
	}
	r = netaddr.IPRangeFrom(startIP.Next(), ip.Prior())

	return
}

// IPRange get subnet range
func (s *IPV4Subnet) IPRange() (r netaddr.IPRange, err error) {
	ip := s.Prefix.IP()
	startIP := ip
	ip, err = addToIP(ip, int32(s.TotalHosts()))
	if err != nil {
		return
	}
	r = netaddr.IPRangeFrom(startIP, ip)

	return
}

func (s *IPV4Subnet) networkRanges(childSubnet *IPV4Subnet) (ranges []netaddr.IPRange, err error) {
	// Can't subdivide to smaller prefixed subnet
	if childSubnet.Prefix.Bits() < s.Prefix.Bits() {
		err = fmt.Errorf("Subnet to split to has more bits %d than parent %d", s.Prefix.Bits(), childSubnet.Prefix.Bits())
		return
	}
	ranges = []netaddr.IPRange{}
	// ip := s.Prefix.IP()
	ip := s.IP
	ipStart := ip

	ratio := int(math.Exp2(float64(childSubnet.Prefix.Bits() - s.Prefix.Bits())))
	for j := 0; j < int(s.Networks()); j++ {
		for r := 0; r < ratio; r++ {
			ip, err = addToIP(ipStart, int32(childSubnet.NetworkHosts()-1))
			if err != nil {
				return
			}
			ranges = append(ranges, netaddr.IPRangeFrom(ipStart, ip))
			ipStart = ip.Next()
		}
	}

	return
}

// NetworkRangesInSubnets set of ranges in the context of subnets of a specified size
func (s *IPV4Subnet) NetworkRangesInSubnets(childSubnet *IPV4Subnet) (ranges []netaddr.IPRange, err error) {
	return s.networkRanges(childSubnet)
}

// NetworkRanges the set of equally sized subnet blocks for subnet
func (s *IPV4Subnet) NetworkRanges() (ranges []netaddr.IPRange, err error) {
	return s.networkRanges(s)
}

// String get string representing subnet (cidr notation)
func (s *IPV4Subnet) String() string {
	return s.Prefix.String()
}
