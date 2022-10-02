package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

const (
	// GlobalUnicast IPV6 type
	GlobalUnicast = iota
	// UniqueLocal IPV6 type
	UniqueLocal
	// InterfaceLocalMulticast IPV6 type
	InterfaceLocalMulticast
	// LinkLocalMulticast IPV6 type
	LinkLocalMulticast
	// LinkLocalUnicast IPV6 type
	LinkLocalUnicast
	// Loopback IPV6 type
	Loopback
	// Multicast IPV6 type
	Multicast
	// Private IPV6 type
	Private
	// Unspecified IPV6 type
	Unspecified
	// Unknown IPV6 type
	Unknown
)

// For fun with generics
func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Use crypto/rand to generate a uint64 with value [0,max]
// There will be no error if max is > 0
func randUInt64(max int64) uint64 {
	bigVal, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	inRange := bigVal.Uint64()

	return inRange
}

// AddrDefaultGateway get IP default gateway for IP
func AddrDefaultGateway(addr netip.Addr) []byte {
	bytes := addr.As16()

	return bytes[:6]
}

// AddrSubnetSection get IP section for IP
func AddrSubnetSection(addr netip.Addr) []byte {
	bytes := addr.As16()
	return bytes[6:8]
}

// AddrGeneralPrefixSection get the general prefix section for IP
func AddrGeneralPrefixSection(addr netip.Addr) []byte {
	bytes := addr.As16()
	return bytes[:8]
}

// AddrRoutingPrefixSecion get routing prefix section for IP
func AddrRoutingPrefixSecion(addr netip.Addr) []byte {
	bytes := addr.As16()
	return bytes[:6]
}

// AddrInterfaceSection get interface section for IP
func AddrInterfaceSection(addr netip.Addr) []byte {
	bytes := addr.As16()
	return bytes[8:]
}

// AddrToBitString complete address binary to 16 bit sections
func AddrToBitString(addr netip.Addr) (result string) {
	str := addr.StringExpanded()

	var sb strings.Builder
	parts := strings.Split(str, ":")
	for _, p := range parts {
		value, err := strconv.ParseInt(p, 16, 64)
		if err != nil {
			return ""
		}
		sb.WriteString(fmt.Sprintf("%08b.", value))
	}

	result = sb.String()
	result = result[:len(result)-1]

	return result
}

// IP6Arpa get the IPV6 ARPA address
func IP6Arpa(addr netip.Addr) string {
	addrStr := addr.StringExpanded()
	addrStr = strings.ReplaceAll(addrStr, ":", "")
	addrSlice := strings.Split(addrStr, "")
	reverse(addrSlice)

	addrStr = fmt.Sprintf("%s.ip6.arpa", strings.Join(addrSlice, "."))
	return addrStr
}

// Bytes2Hex get string with two byte sets delimited by colon
func Bytes2Hex(bytes []byte) string {
	var sb strings.Builder
	for i, byte := range bytes {
		part := fmt.Sprintf("%x", byte)
		if len(part) == 1 {
			sb.WriteString("0")
		}
		sb.WriteString(part)
		if (i+1)%2 == 0 && i != 0 && i != (len(bytes)-1) {
			sb.WriteString(":")
		}
	}
	return sb.String()
}

// TypePrefix the prefix for the IP type
func TypePrefix(addr netip.Addr) (prefix netip.Prefix) {
	kind := AddressType(addr)
	var err error
	switch kind {
	// unique local ipv6 address prefix
	case UniqueLocal:
		prefix, err = netip.ParsePrefix("fd00::/8")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case GlobalUnicast:
		prefix, err = netip.ParsePrefix("2000::/3")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case InterfaceLocalMulticast:
		prefix, err = netip.ParsePrefix("FF00::/8")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case LinkLocalMulticast:
		prefix, err = netip.ParsePrefix("ff00::/8")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case LinkLocalUnicast:
		prefix, err = netip.ParsePrefix("fe80::/10")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case Loopback:
		prefix, err = netip.ParsePrefix("::1/128")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case Multicast:
		prefix, err = netip.ParsePrefix("ff00::/8")
		if err != nil {
			prefix = netip.Prefix{}
		}
		// i.e. unique local
	case Private:
		prefix, err = netip.ParsePrefix("fc00::/7")
		if err != nil {
			prefix = netip.Prefix{}
		}
	case Unspecified:
		prefix = netip.Prefix{}
	default:
		prefix = netip.Prefix{}
	}

	return
}

// AddressType get address type as int
func AddressType(addr netip.Addr) int {
	switch {
	case strings.HasPrefix(addr.StringExpanded(), "fd00"):
		return UniqueLocal
	case addr.IsGlobalUnicast(): // 2001
		return GlobalUnicast
	case addr.IsInterfaceLocalMulticast(): // fe80::/10
		return InterfaceLocalMulticast
	case addr.IsLinkLocalMulticast(): // ff00::/8 ff02
		return LinkLocalMulticast
	case addr.IsLinkLocalUnicast(): // fe80::/10
		return LinkLocalUnicast
	case addr.IsLoopback(): // ::1/128
		return Loopback
	case addr.IsMulticast(): // ff00::/8
		return Multicast
	case addr.IsPrivate(): // fc00::/7
		return Private
	case addr.IsUnspecified():
		return Unspecified
	default:
		return Unknown
	}

}

// AddressTypeName the type of address for the subnet
// https://www.networkacademy.io/ccna/ipv6/ipv6-address-types
func AddressTypeName(addr netip.Addr) string {
	switch AddressType(addr) {
	case UniqueLocal:
		return "Unique local"
	case GlobalUnicast:
		return "Global unicast"
	case InterfaceLocalMulticast:
		return "Interface local multicast"
	case LinkLocalMulticast:
		return "Link local muticast"
	case LinkLocalUnicast:
		return "Link local unicast"
	case Loopback:
		return "Loopback"
	case Multicast:
		return "Multicast"
	case Private:
		return "Private"
	case Unspecified:
		return "Unspecified"
	default:
		return "Unknown"
	}
}

// bytesToMacAddr transform a 6 byte array to a mac address
func bytesToMacAddr(bytes [6]byte) string {
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])

	return macAddress
}

func randomMacBytes2MacAddrBytes(mac [6]byte) ([6]byte, error) {
	addr := net.HardwareAddr(mac[:])
	var bytes [6]byte
	copy(bytes[:], addr[:6])

	return bytes, nil
}

// randomMacAddress make a random mac address of a 6 byte array
func randomMacAddress() (buf [6]byte, err error) {
	buf = [6]byte{}
	_, err = rand.Read(buf[:])
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

// GlobalID get subsection of bits in network part of IP
func GlobalID(addr netip.Addr) (hex string) {
	start := TypePrefix(addr).Bits() + 1
	end := 48

	// for unique local account for L bit
	if AddressType(addr) == UniqueLocal {
		start = start + 1
	}
	// if AddressType(addr) == Multicast {
	// 	start = 128 - 32
	// 	end = 128
	// }

	bytes := addr.As16()
	var arr [8]byte
	copy(arr[:], bytes[0:7])
	data := binary.BigEndian.Uint64(arr[:])

	var dataStr string
	if AddressType(addr) == Multicast || AddressType(addr) == InterfaceLocalMulticast {
		var arr [8]byte
		copy(arr[:], bytes[8:])
		data := binary.BigEndian.Uint64(arr[:])
		data = data << 32
		data = data >> 32
		dataStr = strconv.FormatUint(data, 16)
	} else {
		data = data << start
		data = data >> uint64(64+start-end)
		dataStr = strconv.FormatUint(data, 16)
	}
	if data == 0 {
		return "0000:0000"
	}

	parts := strings.Split(dataStr, "")
	reverse(parts)

	var sb strings.Builder

	for i, letter := range parts {
		sb.WriteString(letter)
		// add colon every 4th letter unless at very end
		if (i+1)%4 == 0 && i != len(parts)-1 {
			sb.WriteString(":")
		}
	}

	parts = strings.Split(sb.String(), "")
	reverse(parts)

	hex = strings.Join(parts, "")

	return
}

// RandomSubnet get a random subnet for IPV6
func RandomSubnet() uint16 {
	rand := randUInt64(65_536)

	return uint16(rand)
}

// mac2GlobalUnicast transform a mac address to a globaal unicast address
func mac2GlobalUnicast(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	inRange := randUInt64(63-32) + 32

	// db8:cafe
	ip := []byte{
		byte(inRange), 0x01, 0xd, 0xb8, 0xca, 0xfe, byte(randUInt64(256)), byte(randUInt64(256)),
		mac[0], mac[1], mac[2], 0xff, 0xfe, mac[3], mac[4], mac[5], // insert ff:fe in the middle
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

// mac2GlobalUnicast transform a mac address to a globaal unicast address
func mac2UniqueLocal(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	// fc00::/8 is currently not defined
	ip := []byte{
		0xfd, 0x0, byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), // prepend with fd00::
		mac[0], mac[1], mac[2], 0xff, 0xfe, mac[3], mac[4], mac[5], // insert ff:fe in the middle
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

// mac2LinkLocal transform a mac address to a link local address
func mac2LinkLocal(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	// link local can be fc00::/7 to fdff::/7
	ip := []byte{
		0xfe, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // prepend with fe80::
		mac[0], mac[1], mac[2], 0xff, 0xfe, mac[3], mac[4], mac[5], // insert ff:fe in the middle
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

// mac2LinkLocal transform a mac address to a link local address
// multicast is tricky - this is not properly implemented in terms of network prefix
// and group id
func mac2Multicast(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	flags := []string{"0", "1", "3", "7"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	scopes := []string{"1", "2", "3", "4", "5", "8", "e", "f"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	flagAndScope, err := strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)

	ip := []byte{
		0xff, byte(flagAndScope),
		0x0, 0x0, byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)),
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

func mac2LinkLocalMulticast(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	flags := []string{"0", "1", "3", "7"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	scopes := []string{"2"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	flagAndScope, err := strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)

	ip := []byte{
		0xff, byte(flagAndScope),
		0x0, 0x0, byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)),
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

func mac2InterfaceLocalMulticast(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	flags := []string{"0", "1", "3", "7"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	scopes := []string{"1"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	flagAndScope, err := strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)

	ip := []byte{
		0xff, byte(flagAndScope),
		0x0, 0x0, byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)),
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

// RandomAddrGlobalUnicast get a global unicast random IPV6 address
func RandomAddrGlobalUnicast() (addr netip.Addr, err error) {
	bytes, err := randomMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := randomMacBytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2GlobalUnicast(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}

// RandomAddrLinkLocal get a link-local random IPV6 address
func RandomAddrLinkLocal() (addr netip.Addr, err error) {
	bytes, err := randomMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := randomMacBytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2LinkLocal(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}

// RandomAddrUniqueLocal get a unique local random IPV6 address
func RandomAddrUniqueLocal() (addr netip.Addr, err error) {
	bytes, err := randomMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := randomMacBytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2UniqueLocal(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}

// RandomAddrMulticast get a random multicast address
func RandomAddrMulticast() (addr netip.Addr, err error) {
	bytes, err := randomMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := randomMacBytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2Multicast(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}

// RandomAddrLinkLocalMulticast get a random link local multicast address
func RandomAddrLinkLocalMulticast() (addr netip.Addr, err error) {
	bytes, err := randomMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := randomMacBytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2LinkLocalMulticast(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}

// RandomAddrInterfaceLocalMulticast get a random interface local multicast address
func RandomAddrInterfaceLocalMulticast() (addr netip.Addr, err error) {
	bytes, err := randomMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := randomMacBytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2InterfaceLocalMulticast(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}
