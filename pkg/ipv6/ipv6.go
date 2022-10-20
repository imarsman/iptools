package ipv6

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

const (
	// GlobalUnicastName name for global unicast type
	GlobalUnicastName = "global-unicast"
	// LinkLocalName name for link local type
	LinkLocalName = "link-local"
	// UniqueLocalName name for unique local type
	UniqueLocalName = "unique-local"
	// PrivateName name for private type
	PrivateName = "private"
	// MulticastName name for multicast type
	MulticastName = "multicast"
	// InterfaceLocalMulticastName name for interface local multicast type
	InterfaceLocalMulticastName = "interface-local-multicast"
	// LinkLocalMulticastName name for link local multicast type
	LinkLocalMulticastName = "link-local-multicast"
)

const (
	// GlobalUnicast IPV6 type
	GlobalUnicast = iota
	// UniqueLocal IPV6 type
	UniqueLocal
	// LinkLocalUnicast IPV6 type
	LinkLocalUnicast
	// Loopback IPV6 type
	Loopback
	// Multicast IPV6 type
	Multicast
	// InterfaceLocalMulticast IPV6 type
	InterfaceLocalMulticast
	// LinkLocalMulticast IPV6 type
	LinkLocalMulticast
	// Private IPV6 type
	Private
	// Unspecified IPV6 type
	Unspecified
	// Unknown IPV6 type
	Unknown
)

var typePrefixes = make(map[int]string)

func init() {
	typePrefixes[GlobalUnicast] = "2000::/3"
	typePrefixes[UniqueLocal] = "fd00::/8"
	typePrefixes[LinkLocalUnicast] = "fe80::/10"
	typePrefixes[Loopback] = "::1/128"
	typePrefixes[Multicast] = "ff00::/8"
	typePrefixes[InterfaceLocalMulticast] = "FF00::/8"
	typePrefixes[LinkLocalMulticast] = "FF00::/8"
	typePrefixes[Private] = "fc00::/7"
}

// AddrTypePrefix the prefix for the IP type
func AddrTypePrefix(addr netip.Addr) (prefix netip.Prefix) {
	addType := AddrType(addr)

	p, ok := typePrefixes[addType]

	switch ok {
	case true:
		var err error
		prefix, err = netip.ParsePrefix(p)
		if err != nil {
			prefix = netip.Prefix{}
		}
	default:
		prefix = netip.Prefix{}
	}

	return
}

// HasType is address const one of a list of candidates
func HasType(t int, candidates ...int) (hasType bool) {
	for _, candidate := range candidates {
		if t == candidate {
			hasType = true
			break
		}
	}

	return
}

// For fun with generics
func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Use crypto/rand to generate a uint64 with value [0,max]
// There will be no error if max is > 0
func randUInt64(max int64) uint64 {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	inRange := bigInt.Uint64()

	return inRange
}

// RandAddrGlobalUnicast get a global unicast random IPV6 address
func RandAddrGlobalUnicast() (addr netip.Addr, err error) {
	macAddrBytes, err := randomMacBytesForInterface(true, true)
	if err != nil {
		return
	}
	macAddr := bytes2MacAddr(macAddrBytes)
	var mac net.HardwareAddr
	mac, err = net.ParseMAC(macAddr)
	if err != nil {
		return
	}

	inRange := randUInt64(63-32) + 32

	addrBytes := [16]byte{
		byte(inRange), 0x01,
		0xd, 0xb8,
		0xca, 0xfe,
		byte(randUInt64(256)), byte(randUInt64(256)),
		mac[0], mac[1],
		mac[2], 0xff,
		0xfe, mac[3],
		mac[4], mac[5],
	}

	addr = netip.AddrFrom16(addrBytes)

	return
}

// RandAddrLinkLocal get a link-local random IPV6 address
func RandAddrLinkLocal() (addr netip.Addr, err error) {
	macAddrBytes, err := randomMacBytesForInterface(true, true)
	if err != nil {
		return
	}
	macAddr := bytes2MacAddr(macAddrBytes)
	mac, err := net.ParseMAC(macAddr)
	if err != nil {
		return
	}

	// link local has prefix FE80::/10
	addrBytes := [16]byte{
		0xfe, 0x80,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x0,
		mac[0], mac[1],
		mac[2], 0xff,
		0xfe, mac[3],
		mac[4], mac[5],
	}
	addr = netip.AddrFrom16(addrBytes)

	return
}

// RandAddrPrivate get a unique local random IPV6 address
func RandAddrPrivate() (addr netip.Addr, err error) {
	macAddrBytes, err := randomMacBytesForInterface(true, true)
	if err != nil {
		return
	}
	macAddr := bytes2MacAddr(macAddrBytes)
	var mac net.HardwareAddr
	mac, err = net.ParseMAC(macAddr)
	if err != nil {
		return
	}

	// Setting last bit (called the L bit) to 1 ensures 0xfd, which is supported
	// The L bit needs to be 1
	first := byte(0xfc)
	first |= 0x1

	// fc00::/7 is currently not defined
	addrBytes := [16]byte{
		first, byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), // prepend with fd00::
		mac[0], mac[1],
		mac[2], 0xff,
		0xfe, mac[3],
		mac[4], mac[5],
	}
	addr = netip.AddrFrom16(addrBytes)

	return
}

// RandAddrMulticast get a random multicast address
func RandAddrMulticast() (addr netip.Addr, err error) {
	// flag for 0 is reserved currently
	flags := []string{"1", "2", "3"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	// scope 1 is interface-local and defined in interfaceLocalMulticast
	// scope 2 is link-local multicast defined in randomLinkLocalMulticast
	scopes := []string{"3", "4", "5", "8", "e", "f"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	flagAndScope, err := strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)

	// multicast has prefix ff00::/8
	addrBytes := [16]byte{
		0xff, byte(flagAndScope),
		0x0, 0x0,
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
	}
	addr = netip.AddrFrom16(addrBytes)

	return
}

// RandAddrLinkLocalMulticast get a random link local multicast address
func RandAddrLinkLocalMulticast() (addr netip.Addr, err error) {
	// flag for 0 is reserved currently
	flags := []string{"1", "2", "3"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	// a single scope applies to link local
	scopes := []string{"2"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	var flagAndScope int64
	flagAndScope, err = strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)
	if err != nil {
		return
	}

	addrBytes := [16]byte{
		0xff, byte(flagAndScope),
		0x0, 0x0,
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
	}
	addr = netip.AddrFrom16(addrBytes)

	return
}

// RandAddrInterfaceLocalMulticast get a random interface local multicast address
func RandAddrInterfaceLocalMulticast() (addr netip.Addr, err error) {
	// flag for 0 is reserved currently
	flags := []string{"1", "2", "3"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	// a single scope applies to interface local
	scopes := []string{"1"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	// get hex value for flag plus scope
	var flagAndScope int64
	flagAndScope, err = strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)
	if err != nil {
		return
	}

	addrBytes := [16]byte{
		0xff, byte(flagAndScope),
		0x0, 0x0,
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
	}
	addr = netip.AddrFrom16(addrBytes)

	return
}

// AddrSolicitedNodeMulticast get solicited node multicast address for incoming unicast address
// EUI-64 compliance
func AddrSolicitedNodeMulticast(addr netip.Addr) (newAddr netip.Addr, err error) {
	if !(HasType(AddrType(addr), GlobalUnicast, LinkLocalUnicast, UniqueLocal, Private)) {
		err = errors.New("not a unicast address")
		return
	}
	// we need the last six characters from the address or 24 bits or 3 bytes
	start := 104
	end := 128

	// get range hex
	var lowOrder24Bits string
	lowOrder24Bits, err = bitRangeHex(addr, start, end)
	if err != nil {
		return
	}
	// strip colons
	lowOrder24Bits = strings.ReplaceAll(lowOrder24Bits, ":", "")

	// hex.DecodeString requires an even number of characters to make bytes
	// and we definitely need 6 characters/3 bytes
	if len(lowOrder24Bits) < 6 {
		lowOrder24Bits = fmt.Sprintf("%s%s", strings.Repeat("0", 6-len(lowOrder24Bits)), lowOrder24Bits)
	}

	// get bytes for decoded string - it will be up to 3 bytes
	rangeBytes, err := hex.DecodeString(lowOrder24Bits)
	if err != nil {
		return
	}

	// copy bytes in range to an array of length 3
	var source [3]byte
	copy(source[:], rangeBytes)

	// Create the IP based on the rule for solicited node multicast
	// ff02::1:ffca:2fdf
	addrBytes := [16]byte{
		0xff, 0x2,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x1,
		0xff, source[0],
		source[1], source[2],
	}

	newAddr = netip.AddrFrom16(addrBytes)

	return
}

// AddrGlobalID get subsection of bits in network part of IP
func AddrGlobalID(addr netip.Addr) (hex string, err error) {
	start := AddrTypePrefix(addr).Bits() + 1
	end := 48
	hex, err = bitRangeHex(addr, start, end)
	// error would be from range > 64 and should not happen
	if err != nil {
		return
	}

	return
}

// AddrMulticastNetworkPrefix get prefix specific to multicast (at end of IP before Group ID)
func AddrMulticastNetworkPrefix(addr netip.Addr) (hex string, err error) {
	start := 32
	end := 32 + 64

	hex, err = bitRangeHex(addr, start, end)
	// error would be from range > 64 and should not happen
	if err != nil {
		return
	}

	return
}

// AddrMulticastGroupID id from range for multicast addresses
func AddrMulticastGroupID(addr netip.Addr) (hex string, err error) {
	start := 128 - 32
	end := 128

	hex, err = bitRangeHex(addr, start, end)
	// error would be from range > 64 and should not happen
	if err != nil {
		return
	}

	return
}

// hex2Delimited make a string of hex digits into an ipv6 colon delimited string
func hex2Delimited(input string) (hex string) {
	input = strings.ReplaceAll(input, ":", "")
	parts := strings.Split(input, "")
	reverse(parts)

	var sb strings.Builder

	// break output into groups of 4 separated by colon
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

// hex2Bytes convert hex string to an int then to an array of bytes
func hex2Bytes(hex string) (rangeBytes [8]byte, err error) {
	hex = strings.ReplaceAll(hex, ":", "")
	value, err := strconv.ParseUint(hex, 64, 16)
	if err != nil {
		return
	}
	binary.LittleEndian.PutUint64(rangeBytes[:], value)

	return
}

// bitRangeHex get hex value for a range of an IP's bits
func bitRangeHex(addr netip.Addr, start, end int) (hex string, err error) {
	if (end - start) > 64 {
		err = errors.New("end-start > 64")
		return
	}
	expectedLen := 10
	startByte := start / 8
	endByte := (end / 8) + 1
	if endByte == 17 {
		endByte = 16
	}

	bytes := addr.As16()
	var arr [8]byte

	if endByte == 16 {
		copy(arr[:], bytes[startByte:])
	} else {
		copy(arr[:], bytes[startByte:endByte])
	}

	var dataStr string
	copy(arr[:], bytes[startByte:])
	data := binary.BigEndian.Uint64(arr[:])

	// remainder is the part of a byte that does not start at the boundary
	// e.g. 3
	remainder := start % 8
	data = data << remainder
	// dataStr = strconv.FormatUint(data, 16)
	data = data >> (64 - (end - start))

	dataStr = strconv.FormatUint(data, 16)
	if data == 0 {
		zeroes := strings.Repeat("0", (end-start)/8)
		return hex2Delimited(zeroes), nil
		// don't add zeroes if the section is not at the beginning of the IP
	} else if len(dataStr) < expectedLen && end == 48 {
		// need at least 40 bytes/ 10 hex chars for global id
		prefix := strings.Repeat("0", expectedLen-len(dataStr))
		dataStr = fmt.Sprintf("%s%s", prefix, dataStr)
	}

	hex = hex2Delimited(dataStr)

	return
}

// addrDefaultGateway get IP default gateway for IP
func addrDefaultGateway(addr netip.Addr) []byte {
	bytes := addr.As16()

	return bytes[:6]
}

// addrSubnetSection get IP section for IP
func addrSubnetSection(addr netip.Addr) []byte {
	bytes := addr.As16()

	return bytes[6:8]
}

// addrGeneralPrefixSection get the general prefix section for IP
func addrGeneralPrefixSection(addr netip.Addr) []byte {
	bytes := addr.As16()

	return bytes[:8]
}

// addrRoutingPrefixSecion get routing prefix section for IP
func addrRoutingPrefixSecion(addr netip.Addr) []byte {
	bytes := addr.As16()

	return bytes[:6]
}

// RoutingPrefix get the routing prefix as a hex string
func RoutingPrefix(addr netip.Addr) string {
	return fmt.Sprintf("%s::/%d", byteSlice2Hex(addrRoutingPrefixSecion(addr)), 48)
}

// Interface get the string representation in hex of the interface bits
func Interface(addr netip.Addr) string {
	return byteSlice2Hex(addrInterfaceSection(addr))
}

// addrInterfaceSection get interface section for IP
func addrInterfaceSection(addr netip.Addr) []byte {
	bytes := addr.As16()
	return bytes[8:]
}

// Addr2BitString complete address binary to 16 bit sections
func Addr2BitString(addr netip.Addr) (result string) {
	str := addr.StringExpanded()

	var sb strings.Builder
	parts := strings.Split(str, ":")
	for _, p := range parts {
		var value int64
		value, err := strconv.ParseInt(p, 16, 64)
		if err != nil {
			return
		}
		sb.WriteString(fmt.Sprintf("%08b.", value))
	}

	result = sb.String()
	result = result[:len(result)-1]

	return
}

// Arpa get the IPV6 ARPA address
func Arpa(addr netip.Addr) (addrStr string) {
	if !HasType(AddrType(addr), GlobalUnicast) {
		return
	}

	addrStr = addr.StringExpanded()
	addrStr = strings.ReplaceAll(addrStr, ":", "")
	addrSlice := strings.Split(addrStr, "")
	reverse(addrSlice)

	addrStr = fmt.Sprintf("%s.ip6.arpa", strings.Join(addrSlice, "."))

	return
}

// byteSlice2Hex get string with two byte sets delimited by colon
func byteSlice2Hex(bytes []byte) string {
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

// AddrType get address type as int
func AddrType(addr netip.Addr) int {
	switch {
	// case strings.HasPrefix(addr.StringExpanded(), "fd00"):
	// 	return UniqueLocal
	case addr.IsInterfaceLocalMulticast(): // fe80::/10
		return InterfaceLocalMulticast
	case addr.IsLinkLocalMulticast(): // ff00::/8 ff02
		return LinkLocalMulticast
	case addr.IsLinkLocalUnicast(): // fe80::/10
		return LinkLocalUnicast
	case addr.IsLoopback(): // ::1/128
		return Loopback
	case addr.IsPrivate(): // fc00::/7
		return Private
	case addr.IsGlobalUnicast(): // 2001
		return GlobalUnicast
	case addr.IsMulticast(): // ff00::/8
		return Multicast
	case addr.IsUnspecified():
		return Unspecified
	default:
		return Unknown
	}

}

// AddrTypeName the type of address for the subnet
// https://www.networkacademy.io/ccna/ipv6/ipv6-address-types
func AddrTypeName(addr netip.Addr) string {
	switch AddrType(addr) {
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

// IsARPA is the address relevant to ARPA addressing?
func IsARPA(addr netip.Addr) (is bool) {
	if HasType(AddrType(addr), GlobalUnicast) {
		is = true
	}

	return
}

// AddrLink get link for address
func AddrLink(addr netip.Addr) (url string) {
	if !HasType(AddrType(addr), GlobalUnicast) {
		return
	}
	return fmt.Sprintf("http://[%s]/", addr.String())
}

// AddrSubnet get the string subnet section as a hex string
func AddrSubnet(addr netip.Addr) string {
	return byteSlice2Hex(addrSubnetSection(addr))
}

// LinkLocalDefaultGateway get default gateway for link local
func LinkLocalDefaultGateway(addr netip.Addr) string {
	gateway := fmt.Sprintf("%s::%d", byteSlice2Hex(addrDefaultGateway(addr)), 1)
	gateway = strings.ReplaceAll(gateway, "0000:", "")

	return gateway
}

// First get firt IP from subnet
func First(addr netip.Addr) netip.Addr {
	bytes := addr.As16()

	for i := 8; i <= 15; i++ {
		bytes[i] = 0x0
	}

	addr = netip.AddrFrom16(bytes)

	return addr
}

// Last get last IP for subnet
func Last(addr netip.Addr) netip.Addr {
	bytes := addr.As16()

	for i := 8; i <= 15; i++ {
		bytes[i] = 0xff
	}

	addr = netip.AddrFrom16(bytes)

	return addr
}

// bytes2MacAddr transform a 6 byte array to a mac address
func bytes2MacAddr(bytes [6]byte) string {
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])

	return macAddress
}

// Get random mac address with global bool flag
// The goal here is to implement EUI-64
// https://community.cisco.com/t5/networking-knowledge-base/understanding-ipv6-eui-64-bit-address/ta-p/3116953
func randomMacBytesForInterface(local, unicast bool) (bytes [6]byte, err error) {
	var mac [6]byte
	_, err = rand.Read(mac[:])
	if err != nil {
		return
	}
	// Best example using random first byte
	// http://cisco.num.edu.mn/CCNA_R&S1/course/module8/8.2.4.5/8.2.4.5.html
	// https://www.practical-go-lessons.com/chap-27-enum-iota-and-bitmask
	// https://developer.epages.com/blog/tech-stories/how-to-generate-mac-addresses/
	// https://community.cisco.com/t5/networking-knowledge-base/understanding-ipv6-eui-64-bit-address/ta-p/3116953
	// https://packetlife.net/blog/2008/aug/4/eui-64-ipv6/
	// https://www.geeksforgeeks.org/ipv6-eui-64-extended-unique-identifier/

	// The first three bytes are vendor-specific and consistent across vendor
	// For example, Cisco's OUI is 20:37:06
	// The 7th of the first byte is 1 for local and 0 for global
	// The 8th of the first byte is 0 for unicast and 1 for multicast

	// The EUI-64 standard specifies that when converting a mac address to an IPV6 interface ID the 7th bit of the first
	// byte of the mac address must be flipped. A 0 becomes a 1 and vice versa. This is not used currently but the
	// flipping of the bits is part of the instructions for converting a mac address to an IPV6 interface ID.

	// Things I don't quite understand
	// - the local/global bit - why flip it instead of just making it one or the other?
	//     - this may be because vendor-produced mac addresses always have this bit set to 0, so "flipping" it
	// 		 would always set to value to 1
	// 	       - this algorithm does not use the first three groupings based on vendor, it uses random values

	// Rules for creating an interface ID from a mac ID are
	// - flip the 7th bit in the first byte (or set to 1 if the first byte is random)
	// - set the 8th bit of the first byte to 1 to indicate unicast
	// - between the third and fourth bytes add 0xFF and OXFe
	//     - this is apparently a disallowed combination in a mac address and allows the fact that the interface ID
	//       started with a mac address to be known
	// - the result will be an 8 byte array, or 64 bits, corresponding to the interface ID length for IPV6

	// fmt.Printf("%08b\n", mac[0])

	// https://en.wikipedia.org/wiki/MAC_address#Ranges_of_group_and_locally_administered_addresses
	// with local == true and unicast == true
	// 01111010 becomes
	// 01111010

	switch local {
	case true:
		// Set 7th bit to 1
		mac[0] |= 0x2
	default:
		// Set 7th bit to 0
		mac[0] &^= 0x2
	}

	switch unicast {
	case true:
		// Set 8th bit to 0 (clear the bit)
		mac[0] &^= 0x1
	default:
		// Set 8th bit to 1
		mac[0] |= 0x1
	}

	addr := net.HardwareAddr(mac[:])
	// fmt.Printf("%08b\n", addr[0])

	copy(bytes[:], addr[:6])

	return
}

// addrRandSubnetID get a random subnet for IPV6
func addrRandSubnetID() uint16 {
	rand := randUInt64(65_536)

	return uint16(rand)
}
