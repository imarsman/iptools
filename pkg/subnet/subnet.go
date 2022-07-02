package subnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	Prefix *netaddr.IPPrefix `json:"prefix" yaml:"prefix"`
}

// NewDefaultFromMask parse for 255.255.255.255 starting mask
func NewDefaultFromMask(mask uint8) (subnet *IPV4Subnet, err error) {
	return newSubnet("255.255.255.255", mask, true)
}

// NewFromMask new subnet with prefix 255.255.255.255 from incoming mask
func NewFromMask(mask uint8) (subnet *IPV4Subnet, err error) {
	return newSubnet("255.255.255.255", mask, true)
}

// NewDefaultFromPrefix new default (255.255.255.255) from prefix
func NewDefaultFromPrefix(prefix string) (subnet *IPV4Subnet, err error) {
	parts := strings.Split(prefix, "/")
	if len(parts) != 2 {
		err = errors.New("invalid prefix")
	}
	mask, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	return newSubnet(parts[0], uint8(mask), true)
}

// NewFromPrefix new using incoming prefix
func NewFromPrefix(prefix string) (subnet *IPV4Subnet, err error) {
	parts := strings.Split(prefix, "/")
	if len(parts) != 2 {
		err = errors.New("invalid prefix")
	}
	mask, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	return newSubnet(parts[0], uint8(mask), true)
}

// newSubnet new subnet with prefix string, masked and boolean flag
func newSubnet(address string, mask uint8, usemasked bool) (subnet *IPV4Subnet, err error) {
	subnet = new(IPV4Subnet)

	var pfx netaddr.IPPrefix
	if usemasked {
		prefixStr := fmt.Sprintf("%s/%d", address, mask)
		var pfxPre netaddr.IPPrefix
		pfxPre, err = netaddr.ParseIPPrefix(prefixStr)
		if err != nil {
			return
		}
		prefixStr = fmt.Sprintf("%s/%d", pfxPre.Masked().IP().String(), mask)
		pfx, err = netaddr.ParseIPPrefix(pfxPre.Masked().String())
		if err != nil {
			return
		}
	} else {
		prefixStr := fmt.Sprintf("%s/%d", address, mask)
		pfx, err = netaddr.ParseIPPrefix(prefixStr)
	}

	if !pfx.IsValid() {
		return nil, errors.New("invalid prefix")
	}

	if pfx.IP().Is6() {
		return nil, errors.New("subnet too large for current implementation")
	}

	subnet.Prefix = &pfx

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

// Hosts bits remaining in mask block
func (s *IPV4Subnet) Hosts() int64 {
	if s.Prefix.Bits()%8 == 0 {
		return int64((math.Exp2(float64(32) - float64(s.Prefix.Bits()))) / float64(s.NetworkCount()))
	}
	return int64((math.Exp2(float64(32 - s.Prefix.Bits()))))
}

// UsableHosts number of usable hosts
func (s *IPV4Subnet) UsableHosts() int64 {
	if s.Hosts() < 2 {
		return 0
	}
	return s.Hosts() - 2
}

// TotalUsableHosts total number of usable hosts
func (s *IPV4Subnet) TotalUsableHosts() int64 {
	return s.Hosts() - 2
}

// // Hosts bits remaining in mask block
// func (s *Subnet) subBlockBits() int64 {
// 	return int64(s.Prefix.Bits()) - int64(s.ClassPartialBits())
// }

// ClassUsableHosts how many usable hosts in subnet?
func (s *IPV4Subnet) ClassUsableHosts() int64 {
	return s.Hosts() - 2
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

// BlockSize get size of blocks
// func (s *Subnet) BlockSize() uint8 {
// 	// number of subnets is number of bits past block border
// 	return uint8(math.Exp2(float64(s.maxBitsForClass() - s.Prefix.Bits())))
// }

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

// NetworkCount number of subnets
func (s *IPV4Subnet) NetworkCount() int64 {
	bits := s.Prefix.Bits() - s.startBitsForClass()
	// bits := s.Prefix.Bits() - 8

	return int64(math.Exp2(float64(bits)))
}

// UsableIPs get usable ips for subnet
func (s *IPV4Subnet) UsableIPs() (ips []netaddr.IP, err error) {
	ips, err = s.IPs()
	if err != nil {
		return
	}
	if len(ips) == 0 {
		err = errors.New("empty ip list for subnet")
		ips = []netaddr.IP{}
		return
	}
	ips = ips[1 : len(ips)-1]

	return
}

// IPs get ips for subnet
func (s *IPV4Subnet) IPs() (ips []netaddr.IP, err error) {
	ip := s.Prefix.IP()
	ips = append(ips, ip)

	for j := 0; j < int(s.Hosts()); j++ {
		ip = ip.Next()
		if (ip == netaddr.IP{}) {
			err = errors.New("empty ip list for subnet")
			ips = []netaddr.IP{}
			return
		}
		ips = append(ips, ip)
	}

	return
}

// Range get subnet range
func (s *IPV4Subnet) Range() (r netaddr.IPRange, err error) {
	ip := s.Prefix.IP()
	startIP := ip
	for j := 0; j < int(s.Hosts()); j++ {
		ip = ip.Next()
		if (ip == netaddr.IP{}) {
			err = errors.New("empty ip list in subnet range")
			r = netaddr.IPRange{}
			return
		}
	}
	r = netaddr.IPRangeFrom(startIP, ip)

	return
}

// NetworkRanges the set of equally sized subnet blocks for subnet
func (s *IPV4Subnet) NetworkRanges() (r []netaddr.IPRange, err error) {
	r = []netaddr.IPRange{}
	ip := s.Prefix.IP()
	ipStart := ip

	for j := 0; j < int(s.NetworkCount()); j++ {
		for j := 0; j < int(s.Hosts()); j++ {
			ip = ip.Next()
			if (ip == netaddr.IP{}) {
				err = errors.New("empty ip list in subnet range")
				r = []netaddr.IPRange{}
				return
			}
		}
		r = append(r, netaddr.IPRangeFrom(ipStart, ip))
		ipStart = ip.Next()
	}

	return
}

func (s *IPV4Subnet) String() string {
	return s.Prefix.String()
}

// Subdivide subnet into child network sized networks
func (s *IPV4Subnet) networks(childSubnet *IPV4Subnet) (subnets []*IPV4Subnet, err error) {
	fmt.Println(s.Prefix.Bits(), childSubnet.Prefix.Bits())
	if s.Prefix.Bits() > childSubnet.Prefix.Bits() {
		err = fmt.Errorf("Subnet to split to has more bits %d than parent %d", s.Prefix.Bits(), childSubnet.Prefix.Bits())
		return
	}

	if s.Prefix.Bits() == 1 {
		err = fmt.Errorf("Can't subdivide")
		return
	}
	ip := s.Prefix.IP()
	ipStart := ip

	// ratio := int(math.Exp2(float64(s.NetworkCount()) - float64(childSubnet.NetworkCount())))
	ratio := int(math.Exp2(float64(childSubnet.Prefix.Bits() - s.Prefix.Bits())))
	for j := 0; j < int(s.NetworkCount()); j++ {
		for r := 0; r < ratio; r++ {
			for j := 0; j < int(childSubnet.Hosts()); j++ {
				ip = ip.Next()
				if (ip == netaddr.IP{}) {
					err = errors.New("empty ip list in subnet range")
					subnets = []*IPV4Subnet{}
					return
				}
			}
			subnet, err := newSubnet(ipStart.String(), childSubnet.Prefix.Bits(), false)
			if err != nil {
				return []*IPV4Subnet{}, err
			}
			ipStart = ip
			subnets = append(subnets, subnet)
		}
	}

	return
}

// NetworksInSubnets get networks split into child subnet sized networks
func (s *IPV4Subnet) NetworksInSubnets(childSubnet *IPV4Subnet) (subnets []*IPV4Subnet, err error) {
	return s.networks(childSubnet)
}

// Networks get all networks for subnet in subnet sized networks
func (s *IPV4Subnet) Networks() (subnets []*IPV4Subnet, err error) {
	return s.networks(s)
}
