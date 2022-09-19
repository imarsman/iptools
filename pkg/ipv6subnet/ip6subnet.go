package ipv6subnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"

	"github.com/imarsman/iptools/pkg/ipv6subnet/util"
	"gopkg.in/yaml.v2"
)

// Subnet an IP subnet
type Subnet struct {
	name   string
	prefix netip.Prefix
}

// Name subnet name
func (s *Subnet) Name() string {
	return s.name
}

// SetName set subnet name
func (s *Subnet) SetName(name string) {
	s.name = name
}

// Prefix get prefix for subnet
func (s *Subnet) Prefix() netip.Prefix {
	return s.prefix
}

// Addr get Addr for subnet
func (s *Subnet) Addr() netip.Addr {
	return s.prefix.Addr()
}

func (s *Subnet) SubnetString() string {
	bytes := s.Addr().Next().AsSlice()
	return util.Bytes2Hex(bytes[6:8])
}

func (s *Subnet) InterfaceString() string {
	start := s.prefix.Bits() / 8
	bytes := s.Addr().AsSlice()
	return util.Bytes2Hex(bytes[start:])
}

func (s *Subnet) PrefixString() string {
	end := s.prefix.Bits() / 8
	bytes := s.Addr().AsSlice()
	return util.Bytes2Hex(bytes[:end])
}

// String get string representing subnet (cidr notation)
func (s *Subnet) String() string {
	return s.prefix.String()
}

// JSON get JSON for subnet
func (s *Subnet) JSON() (bytes []byte, err error) {
	var prefix = s.prefix.String()
	bytes, err = json.MarshalIndent(&prefix, "", "  ")
	if err != nil {
		return
	}

	return bytes, nil
}

// YAML get YAML for subnet
func (s *Subnet) YAML() (bytes []byte, err error) {
	var prefix = s.prefix.String()
	bytes, err = yaml.Marshal(&prefix)
	if err != nil {
		return
	}

	return bytes, nil
}

// NewNamedFromIPAndBits new with name using incoming prefix ip and network bits
func NewNamedFromIPAndBits(addr string, bits int, name string) (subnet *Subnet, err error) {
	subnet, err = newSubnet(addr, bits)
	if err != nil {
		return
	}
	subnet.name = name

	return
}

// NewFromIPAndBits new using incoming prefix ip and network bits
func NewFromIPAndBits(addr string, bits int) (subnet *Subnet, err error) {
	return newSubnet(addr, bits)
}

// NewNamedFromPrefix new with name using incoming prefix
func NewNamedFromPrefix(prefix string, name string) (subnet *Subnet, err error) {
	p, err := netip.ParseAddr(prefix)
	if err != nil {
		return
	}

	subnet, err = newSubnet(p.String(), p.BitLen())
	if err != nil {
		return
	}
	subnet.name = name

	return
}

// NewFromPrefix new using incoming prefix
func NewFromPrefix(prefix string) (subnet *Subnet, err error) {
	p, err := netip.ParsePrefix(prefix)
	if err != nil {
		return
	}

	subnet, err = newSubnet(p.Addr().String(), p.Bits())
	if err != nil {
		return
	}

	return
}

// newSubnet new subnet with prefix ip and network bits
func newSubnet(address string, bits int) (subnet *Subnet, err error) {
	subnet = new(Subnet)
	addr, err := netip.ParseAddr(address)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if addr.Is4() {
		return nil, errors.New("subnet too large for current implementation")
	}

	subnet.prefix = netip.PrefixFrom(addr, bits)
	if err != nil {
		return
	}

	return
}
