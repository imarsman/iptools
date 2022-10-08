package ip6util

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

	return
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

// ByteSlice2Hex get string with two byte sets delimited by colon
func ByteSlice2Hex(bytes []byte) string {
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

// bytes2MacAddr transform a 6 byte array to a mac address
func bytes2MacAddr(bytes [6]byte) string {
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])

	return macAddress
}

func randomMacBytesForInterface() (bytes [6]byte, err error) {
	var mac [6]byte
	_, err = rand.Read(mac[:])
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Set second last bit to 1 to indicate globally unique
	// See also https://www.rfc-editor.org/rfc/rfc4291.html#section-2.5.1

	// fmt.Printf("%08b\n", mac[0])
	mac[0] |= (1 << (2 - 1))
	addr := net.HardwareAddr(mac[:])
	// fmt.Printf("%08b\n", addr[0])

	copy(bytes[:], addr[:6])

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

	return bitRangeHex(addr, start, end)
}

// MulticastNetworkPrefix get prefix specific to multicast (at end of IP before Group ID)
func MulticastNetworkPrefix(addr netip.Addr) (hex string) {
	start := 32
	end := 32 + 64

	return bitRangeHex(addr, start, end)
}

// MulticastGroupID id from range for multicast addresses
func MulticastGroupID(addr netip.Addr) (hex string) {
	start := 128 - 32
	end := 128

	return bitRangeHex(addr, start, end)
}

// hexStringToDelimited make a string of hex digits into an ipv6 colon delimited string
func hexStringToDelimited(input string) (hex string) {
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

func hexStringToBytes(hex string) (rangeBytes [8]byte, err error) {
	hex = strings.ReplaceAll(hex, ":", "")
	value, err := strconv.ParseUint(hex, 64, 16)
	if err != nil {
		return
	}
	binary.LittleEndian.PutUint64(rangeBytes[:], value)

	return
}

func bitRangeHex(addr netip.Addr, start, end int) (hex string) {
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
		return hexStringToDelimited(zeroes)
		// don't add zeroes if the section is not at the beginning of the IP
	} else if len(dataStr) < expectedLen && end == 48 {
		// need at least 40 bytes/ 10 hex chars for global id
		prefix := strings.Repeat("0", expectedLen-len(dataStr))
		dataStr = fmt.Sprintf("%s%s", prefix, dataStr)
	}

	hex = hexStringToDelimited(dataStr)

	return
}

// RandomSubnet get a random subnet for IPV6
func RandomSubnet() uint16 {
	rand := randUInt64(65_536)

	return uint16(rand)
}

// RandomAddrGlobalUnicast get a global unicast random IPV6 address
func RandomAddrGlobalUnicast() (addr netip.Addr, err error) {
	addr, err = randomGlobalUnicast()
	if err != nil {
		return
	}

	return
}

// RandomAddrLinkLocal get a link-local random IPV6 address
func RandomAddrLinkLocal() (addr netip.Addr, err error) {
	addr, err = randomLinkLocal()
	if err != nil {
		return
	}

	return
}

// RandomAddrPrivate get a unique local random IPV6 address
func RandomAddrPrivate() (addr netip.Addr, err error) {
	addr, err = randomPrivate()
	if err != nil {
		return
	}

	return
}

// RandomAddrMulticast get a random multicast address
func RandomAddrMulticast() (addr netip.Addr, err error) {
	addr, err = randomMulticast()
	if err != nil {
		return
	}

	return
}

// RandomAddrLinkLocalMulticast get a random link local multicast address
func RandomAddrLinkLocalMulticast() (addr netip.Addr, err error) {
	addr, err = randomLinkLocalMulticast()
	if err != nil {
		return
	}

	return
}

// RandomAddrInterfaceLocalMulticast get a random interface local multicast address
func RandomAddrInterfaceLocalMulticast() (addr netip.Addr, err error) {
	addr, err = randomInterfaceLocalMulticast()
	if err != nil {
		return
	}

	return
}

// AddrSolicitedNodeMulticast get solicited node multicast address for incoming unicast address
func AddrSolicitedNodeMulticast(addr netip.Addr) (newAddr netip.Addr, err error) {
	if !(HasType(AddressType(addr), GlobalUnicast, LinkLocalUnicast, UniqueLocal, Private)) {
		err = errors.New("not a unicast address")
		return
	}
	// we need the last six characters from the address or 24 bits or 3 bytes
	start := 104
	end := 128

	// get range hex
	unique := bitRangeHex(addr, start, end)
	// strip colons
	unique = strings.ReplaceAll(unique, ":", "")

	// hex.DecodeString requires an even number of characters to make bytes
	// and we definitely need 6 characters/3 bytes
	if len(unique) < 6 {
		unique = fmt.Sprintf("%s%s", strings.Repeat("0", 6-len(unique)), unique)
	}

	// get bytes for decoded string - it will be up to 3 bytes
	rangeBytes, err := hex.DecodeString(unique)
	if err != nil {
		panic(err)
	}
	// copy bytes in range to an array of length 3
	var source [3]byte
	copy(source[:], rangeBytes)

	// Create the IP based on the rule for solicited node multicast
	// ff02::1:ffca:2fdf
	ipBytes := []byte{
		0xff, 0x2,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x1,
		0xff, source[0],
		source[1], source[2],
	}

	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	newAddr = netip.AddrFrom16(addrBytes)

	return
}

// randomGlobalUnicast transform a mac address to a globaal unicast address
// https://support.lenovo.com/ca/en/solutions/ht509925-how-to-convert-a-mac-address-into-an-ipv6-link-local-address-eui-64
func randomGlobalUnicast() (addr netip.Addr, err error) {
	macAddrBytes, err := randomMacBytesForInterface()
	if err != nil {
		return
	}
	s := bytes2MacAddr(macAddrBytes)
	mac, err := net.ParseMAC(s)
	if err != nil {
		return
	}

	inRange := randUInt64(63-32) + 32

	// db8:cafe
	ipBytes := []byte{
		byte(inRange), 0x01,
		0xd, 0xb8,
		0xca, 0xfe,
		byte(randUInt64(256)), byte(randUInt64(256)),
		mac[0], mac[1],
		mac[2], 0xff,
		0xfe, mac[3],
		mac[4], mac[5],
	}

	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	addr = netip.AddrFrom16(addrBytes)

	return
}

// mac2GlobalUnicast transform a mac address to a globaal unicast address
func randomPrivate() (addr netip.Addr, err error) {
	macAddrBytes, err := randomMacBytesForInterface()
	if err != nil {
		return
	}
	s := bytes2MacAddr(macAddrBytes)
	mac, err := net.ParseMAC(s)
	if err != nil {
		return
	}

	// fc00::/7 is currently not defined
	ipBytes := []byte{
		0xfd, byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)), // prepend with fd00::
		mac[0], mac[1],
		mac[2], 0xff,
		0xfe, mac[3],
		mac[4], mac[5],
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	addr = netip.AddrFrom16(addrBytes)

	return
}

// randomLinkLocal transform a mac address to a link local address
func randomLinkLocal() (addr netip.Addr, err error) {
	macAddrBytes, err := randomMacBytesForInterface()
	if err != nil {
		return
	}
	s := bytes2MacAddr(macAddrBytes)
	mac, err := net.ParseMAC(s)
	if err != nil {
		return
	}

	// link local has prefix FE80::/10
	ipBytes := []byte{
		0xfe, 0x80,
		0x0, 0x0,
		0x0, 0x0,
		0x0, 0x0,
		mac[0], mac[1],
		mac[2], 0xff,
		0xfe, mac[3],
		mac[4], mac[5],
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	addr = netip.AddrFrom16(addrBytes)

	return
}

// mac2LinkLocal transform a mac address to a link local address
// multicast is tricky - this is not properly implemented in terms of network prefix
// and group id
func randomMulticast() (addr netip.Addr, err error) {
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
	ipBytes := []byte{
		0xff, byte(flagAndScope),
		0x0, 0x0,
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	addr = netip.AddrFrom16(addrBytes)

	return
}

func randomLinkLocalMulticast() (addr netip.Addr, err error) {
	// flag for 0 is reserved currently
	flags := []string{"1", "2", "3"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	// a single scope applies to link local
	scopes := []string{"2"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	flagAndScope, err := strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)

	ipBytes := []byte{
		0xff, byte(flagAndScope),
		0x0, 0x0,
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	addr = netip.AddrFrom16(addrBytes)

	return
}

func randomInterfaceLocalMulticast() (addr netip.Addr, err error) {
	// flag for 0 is reserved currently
	flags := []string{"1", "2", "3"}
	element := randUInt64(int64(len(flags))) + 1
	flagStr := flags[element-1]

	// a single scope applies to interface local
	scopes := []string{"1"}
	element = randUInt64(int64(len(scopes))) + 1
	scopeStr := scopes[element-1]

	flagAndScope, err := strconv.ParseInt(fmt.Sprintf("%s%s", flagStr, scopeStr), 16, 64)

	ipBytes := []byte{
		0xff, byte(flagAndScope),
		0x0, 0x0,
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
		byte(randUInt64(256)), byte(randUInt64(256)),
	}
	var addrBytes [16]byte
	copy(addrBytes[:], ipBytes)
	addr = netip.AddrFrom16(addrBytes)

	return
}
