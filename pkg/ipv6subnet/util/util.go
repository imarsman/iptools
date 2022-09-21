package util

import (
	"fmt"
	"math/rand"
	"net"
	"net/netip"
	"strings"
	"time"
)

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

// mac2LinkLocal transform a mac address to a link local address
func mac2LinkLocal(s string) (netip.Addr, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return netip.Addr{}, err
	}

	// Invert the bit at the index 6 (counting from 0)
	mac[0] ^= (1 << (2 - 1))

	ip := []byte{
		0xfe, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // prepend with fe80::
		mac[0], mac[1], mac[2], 0xff, 0xfe, mac[3], mac[4], mac[5], // insert ff:fe in the middle
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
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

	// db8:cafe
	ip := []byte{
		// 0x20, 0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // prepend with fe80::
		// 0x20, 0x01, 0x0, 0x0, 0x0, 0x0, byte(rand.Intn(256)), byte(rand.Intn(256)), // prepend with fe80::
		0x20, 0x01, 0xd, 0xb8, 0xca, 0xfe, byte(rand.Intn(256)), byte(rand.Intn(256)), // prepend with 2001::
		mac[0], mac[1], mac[2], 0xff, 0xfe, mac[3], mac[4], mac[5], // insert ff:fe in the middle
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ip)

	return netip.AddrFrom16(addrBytes), nil
}

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

// // BitStr4 get bit string for IPV4 IP
// func BitStr4(ip netip.Addr, separator string) string {
// 	bytes := ip.As4()
// 	list := []string{}

// 	for _, b := range bytes {
// 		list = append(list, fmt.Sprintf("%08b", b))
// 	}

// 	return strings.Join(list, separator)
// }

// // BinaryIP4StrToBytes convert a binary IP string for an IPV4 IP to a set of 4 bytes
// func BinaryIP4StrToBytes(ip string) (list []byte, err error) {
// 	var parts []string
// 	if !strings.Contains(ip, ".") {
// 		if len(ip) != 32 {
// 			err = fmt.Errorf("invalid binary ip %s", ip)
// 			return
// 		}
// 		var matches bool
// 		matches, err = regexp.MatchString(`[01]`, ip)
// 		if err != nil {
// 			return
// 		}
// 		if !matches {
// 			err = fmt.Errorf("invalid binary ip %s", ip)
// 			return
// 		}
// 		parts = []string{ip[0:8], ip[8:16], ip[16:24], ip[25:32]}
// 	} else {
// 		parts = strings.Split(ip, ".")
// 	}

// 	for _, part := range parts {
// 		var i int64
// 		i, err = strconv.ParseInt(part, 2, 32)
// 		if err != nil {
// 			return
// 		}
// 		list = append(list, byte(i))
// 	}

// 	return
// }

// https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
// func int2ip(ipInt uint32) (netip.Addr, error) {
// 	ip := make(net.IP, 4)
// 	binary.BigEndian.PutUint32(ip, ipInt)
// 	addr, ok := netip.ParseAddr(ip.String())

// 	return addr, ok
// }

// WildCardMask get mask bits available for addressing
// func WildCardMask(ip netip.Addr) string {
// 	bytes := ip.As4()
// 	var list = make([]string, 4, 4)
// 	for i, b := range bytes {
// 		if b == 255 {
// 			list[i] = fmt.Sprint(0)
// 			continue
// 		}
// 		list[i] = fmt.Sprint(250 - b)
// 	}

// 	return strings.Join(list, `.`)
// }

// AddToIPIP4 add count IPs to IP
// func AddToIPIP4(startIP netip.Addr, add int32) (addedIP netip.Addr, err error) {
// 	if !startIP.Next().IsValid() {
// 		err = fmt.Errorf("ip %v is already max", startIP)
// 		return
// 	}
// 	bytes := startIP.As4()
// 	slice := bytes[:]
// 	ipValue := binary.BigEndian.Uint32(slice)
// 	ipValue += uint32(add)
// 	addedIP, err = int2ip(ipValue)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// For fun with generics
func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// InAddrArpa get the InAddrArpa version of an IP
// This has not been tested
// result will be [dot separated].ip6.arpa
// e.g. 1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa
// func InAddrArpa(ip netip.Addr) string {
// 	ipStr := ip.StringExpanded()

// 	var chars = make([]string, 0, 32)
// 	for _, r := range strings.Join(strings.Split(ipStr, `:`), "") {
// 		chars = append(chars, string(r))
// 	}

// 	reverse(chars)

// 	return fmt.Sprintf("%s.%s", strings.Join(chars, "."), ".ip6.arpa")
// }
