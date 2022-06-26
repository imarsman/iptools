package subnet

import (
	"encoding/json"
	"errors"
	"math"

	"gopkg.in/yaml.v2"
	"inet.af/netaddr"
)

const (
	octets   = 4
	octetMax = 255
)

// Subnet an IP subnet
// https://www.calculator.net/ip-Subnet-calculator.html?cclass=any&csubnet=16&cip=10.0.0.0&ctype=ipv4&printit=0&x=43&y=21
// https://www.calculator.net/ip-Subnet-calculator.html
type Subnet struct {
	Prefix *netaddr.IPPrefix `json:"prefix" yaml:"prefix"`
}

// NewSubnet create new subnet from prefix. Use canonical form if usemasked is true
func NewSubnet(prefix string, usemasked bool) (subnet *Subnet, err error) {
	return newSubnet(prefix, usemasked)
}

func newSubnet(prefix string, usemaked bool) (subnet *Subnet, err error) {
	subnet = new(Subnet)

	pfx, err := netaddr.ParseIPPrefix(prefix)
	if err != nil {
		return
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

// UsableRange omit first and last IPs
func UsableRange() (usableRange netaddr.IPRange) {

	return
}

// JSON get JSON for subnet
func (s *Subnet) JSON() (bytes []byte, err error) {
	bytes, err = json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	return bytes, nil
}

// YAML get YAML for subnet
func (s *Subnet) YAML() (bytes []byte, err error) {
	bytes, err = yaml.Marshal(s)
	if err != nil {
		return
	}

	return bytes, nil
}

// ClassBits byte used for subnet
func (s *Subnet) ClassBits() uint8 {
	return s.Prefix.Bits()
}

// ClassPartialBits bits used in mask block
func (s *Subnet) ClassPartialBits() uint8 {
	return s.Prefix.Bits() % 8
}

// UsableHosts number of usable hosts
func (s *Subnet) UsableHosts() float64 {
	return math.Exp2(float64(s.ClassPartialBits())) - 2
}

// Hosts bits remaining in mask block
func (s *Subnet) Hosts() int64 {
	return int64(math.Exp2(float64(s.ClassPartialBits())))
}

// ClassUsableHosts how many usable hosts in subnet?
func (s *Subnet) ClassUsableHosts() int64 {
	return s.Hosts() - 2
}

// MaxBitsForClass maximum bits for subnet range for the class
func (s *Subnet) MaxBitsForClass() uint8 {
	if s.Prefix.Bits() <= 8 {
		return 8
	} else if s.Prefix.Bits() <= 16 {
		return 16
	} else if s.Prefix.Bits() <= 24 {
		return 24
	}
	return 32
}

func (s *Subnet) TotalHosts() int64 {
	return int64(s.BlockSize()) * s.Hosts()
}

// BlockSize get size of blocks
func (s *Subnet) BlockSize() uint8 {
	return uint8(math.Exp2(float64(s.MaxBitsForClass() - s.Prefix.Bits())))
}

// // SubnetSize number of hosts per subnet
// func (s *Subnet) SubnetSize() uint8 {
// 	return 1
// }

// Subnets number of subnets
func (s *Subnet) Subnets() uint8 {
	return s.BlockSize()
}

// UsableIPs get usable ips for subnet
func (s *Subnet) UsableIPs() (ips []netaddr.IP, err error) {
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
func (s *Subnet) IPs() (ips []netaddr.IP, err error) {
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
func (s *Subnet) Range() (r netaddr.IPRange, err error) {
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

// Blocks the set of equally sized subnet blocks for subnet
func (s *Subnet) Blocks() (r []netaddr.IPRange, err error) {
	r = []netaddr.IPRange{}
	ip := s.Prefix.IP()
	ipStart := ip

	for j := 0; j < int(s.BlockSize()); j++ {
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
