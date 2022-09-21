package ipv6subnet

import (
	// "crypto/rand"

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

// First get firt IP from subnet
func (s *Subnet) First() netip.Addr {
	addr := s.prefix.Addr()
	bytes := addr.As16()
	bytes[8] = 0x0
	bytes[9] = 0x0
	bytes[10] = 0x0
	bytes[11] = 0x0
	bytes[12] = 0x0
	bytes[13] = 0x0
	bytes[14] = 0x0
	bytes[15] = 0x0

	addr = netip.AddrFrom16(bytes)

	return addr
}

// First get firt IP from subnet
func (s *Subnet) Last() netip.Addr {
	addr := s.prefix.Addr()
	bytes := addr.As16()
	bytes[8] = 0xff
	bytes[9] = 0xff
	bytes[10] = 0xff
	bytes[11] = 0xff
	bytes[12] = 0xff
	bytes[13] = 0xff
	bytes[14] = 0xff
	bytes[15] = 0xff

	addr = netip.AddrFrom16(bytes)

	return addr
}

// // AddressType the type of address for the subnet
// // https://www.networkacademy.io/ccna/ipv6/ipv6-address-types
// func (s *Subnet) AddressType() string {
// 	switch {
// 	case s.Addr().IsGlobalUnicast(): // 2001
// 		return "Global unicast"
// 	case s.Addr().IsInterfaceLocalMulticast(): // fe80::/10
// 		return "Interface local multicast"
// 	case s.Addr().IsLinkLocalMulticast(): // ff00::/8 ff02
// 		return "Link local muticast"
// 	case s.Addr().IsLinkLocalUnicast(): // fe80::/10
// 		return "Link local unicast"
// 	case s.Addr().IsLoopback(): // ::1/128
// 		return "Loopback"
// 	case s.Addr().IsMulticast(): // ff00::/8
// 		return "Multicast"
// 	case s.Addr().IsPrivate(): // fc00::/7
// 		return "Private"
// 	case s.Addr().IsUnspecified():
// 		return "Unspecified"
// 	default:
// 		return "Unknown"
// 	}
// }

// SubnetString get the string representation in hex of the subnet bits
func (s *Subnet) SubnetString() string {
	bytes := s.Addr().Next().AsSlice()
	return util.Bytes2Hex(bytes[6:8])
}

func (s *Subnet) RoutingPrefixString() string {
	bytes := s.Addr().Next().AsSlice()
	return util.Bytes2Hex(bytes[:6])
}

// InterfaceString get the string representation in hex of the interface bits
func (s *Subnet) InterfaceString() string {
	start := s.prefix.Bits() / 8
	bytes := s.Addr().AsSlice()
	return util.Bytes2Hex(bytes[start:])
}

// PrefixString get the string representation in hex of the address prefix
func (s *Subnet) PrefixString() string {
	end := s.prefix.Bits() / 8
	bytes := s.Addr().AsSlice()
	return util.Bytes2Hex(bytes[:end])
}

// String get string representing subnet of the subnet prefix
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
