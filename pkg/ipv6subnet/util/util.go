package util

import (
	"fmt"
	"math/rand"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

// AddrGlobalID get IP global ID
func AddrGlobalID(addr netip.Addr) []byte {
	bytes := addr.As16()
	return bytes[2:6]
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

// AddressType the type of address for the subnet
// https://www.networkacademy.io/ccna/ipv6/ipv6-address-types
func AddressType(addr netip.Addr) string {
	switch {
	case strings.HasPrefix(addr.StringExpanded(), "fd00"):
		return "Unique local"
	case addr.IsGlobalUnicast(): // 2001
		return "Global unicast"
	case addr.IsInterfaceLocalMulticast(): // fe80::/10
		return "Interface local multicast"
	case addr.IsLinkLocalMulticast(): // ff00::/8 ff02
		return "Link local muticast"
	case addr.IsLinkLocalUnicast(): // fe80::/10
		return "Link local unicast"
	case addr.IsLoopback(): // ::1/128
		return "Loopback"
	case addr.IsMulticast(): // ff00::/8
		return "Multicast"
	case addr.IsPrivate(): // fc00::/7
		return "Private"
	case addr.IsUnspecified():
		return "Unspecified"
	default:
		return "Unknown"
	}
}

// bytesToMacAddr transform a 6 byte array to a mac address
func bytesToMacAddr(bytes []byte) string {
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])

	return macAddress
}

func bytes2MacAddrBytes(mac [6]byte) ([]byte, error) {
	addr := net.HardwareAddr(mac[:])

	return addr, nil
}

// makeMacAddress make a random mac address of a 6 byte array
func makeMacAddress() (buf [6]byte, err error) {
	rand.Seed(time.Now().Unix())
	buf = [6]byte{}
	_, err = rand.Read(buf[:])
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

// RandomSubnet get a random subnet for IPV6
func RandomSubnet() uint16 {
	rand.Seed(time.Now().Unix())

	max := 65_536

	rand := rand.Intn(max)

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

	rand.Seed(time.Now().Unix())

	// global unicast must start with 2000 to 3fff
	// The first 8 bytes define the range
	// min: 32 = hex 20
	// max: 63 = hex 3f
	inRange := rand.Intn(63-32) + 32

	// db8:cafe
	ip := []byte{
		// 0x20, 0x01, 0xd, 0xb8, 0xca, 0xfe, byte(rand.Intn(256)), byte(rand.Intn(256)), // prepend with 2001::
		byte(inRange), 0x01, 0xd, 0xb8, 0xca, 0xfe, byte(rand.Intn(256)), byte(rand.Intn(256)),
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

	rand.Seed(time.Now().Unix())

	ip := []byte{
		0xfd, 0x0, byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), // prepend with fd00::
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

// https://en.wikipedia.org/wiki/IPv6_address

// RandomAddrGlobalUnicast get a global unicast random IPV6 address
func RandomAddrGlobalUnicast() (addr netip.Addr, err error) {
	bytes, err := makeMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := bytes2MacAddrBytes(bytes)
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
	bytes, err := makeMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := bytes2MacAddrBytes(bytes)
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
	bytes, err := makeMacAddress()
	if err != nil {
		return
	}
	macAddrBytes, err := bytes2MacAddrBytes(bytes)
	if err != nil {
		return
	}
	addr, err = mac2UniqueLocal(bytesToMacAddr(macAddrBytes))
	if err != nil {
		return
	}

	return
}

// For fun with generics
func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
