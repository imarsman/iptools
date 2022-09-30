package ipv6subnet

import (
	// "crypto/rand"

	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"strings"

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

// Last get last IP for subnet
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

// GlobalIDString get the string global ID a hex string
func (s *Subnet) GlobalIDString() string {
	return util.Bytes2Hex([]byte(util.GlobalID(s.Addr())))
}

// SubnetString get the string subnet section as a hex string
func (s *Subnet) SubnetString() string {
	return util.Bytes2Hex(util.AddrSubnetSection(s.Addr()))
}

// DefaultGatewayString get the default gateway as a hex string
func (s *Subnet) DefaultGatewayString() string {
	return fmt.Sprintf("%s::%d", util.Bytes2Hex(util.AddrDefaultGateway(s.Addr())), 1)
}

// Link link version of address
func (s *Subnet) Link() (url string) {
	return fmt.Sprintf("http://[%s]/", s.Addr().String())
}

// TypePrefix prefix make a prefix for a type
func (s *Subnet) TypePrefix() (prefix netip.Prefix) {

	return util.TypePrefix(s.Addr())
}

// RoutingPrefixString get the routing prefix as a hex string
func (s *Subnet) RoutingPrefixString() string {
	if strings.HasPrefix(s.Addr().StringExpanded(), "fd00") {
		return fmt.Sprintf("%s::/%d", "fd00", 48)
	}
	return fmt.Sprintf("%s::/%d", util.Bytes2Hex(util.AddrRoutingPrefixSecion(s.Addr())), 48)
}

// InterfaceString get the string representation in hex of the interface bits
func (s *Subnet) InterfaceString() string {
	return util.Bytes2Hex(util.AddrInterfaceSection(s.Addr()))
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
