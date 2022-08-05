package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
)

// BitStr4 get bit string for IPV4 IP
func BitStr4(ip netip.Addr, separator string) string {
	bytes := ip.As4()
	list := []string{}

	for _, b := range bytes {
		list = append(list, fmt.Sprintf("%08b", b))
	}

	return strings.Join(list, separator)
}

// BinaryIP4StrToBytes convert a binary IP string for an IPV4 IP to a set of 4 bytes
func BinaryIP4StrToBytes(ip string) (list []byte, err error) {
	var parts []string
	if !strings.Contains(ip, ".") {
		if len(ip) != 32 {
			err = fmt.Errorf("invalid binary ip %s", ip)
			return
		}
		var matches bool
		matches, err = regexp.MatchString(`[01]`, ip)
		if err != nil {
			return
		}
		if !matches {
			err = fmt.Errorf("invalid binary ip %s", ip)
			return
		}
		parts = []string{ip[0:8], ip[8:16], ip[16:24], ip[25:32]}
	} else {
		parts = strings.Split(ip, ".")
	}

	for _, part := range parts {
		var i int64
		i, err = strconv.ParseInt(part, 2, 32)
		if err != nil {
			return
		}
		list = append(list, byte(i))
	}

	return
}

// https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func int2ip(ipInt uint32) (netip.Addr, error) {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	addr, ok := netip.ParseAddr(ip.String())

	return addr, ok
}

// WildCardMaskIP4 get mask bits available for addressing
func WildCardMaskIP4(ip netip.Addr) string {
	bytes := ip.As4()
	var list = make([]string, 4, 4)
	for i := range bytes {
		list[i] = fmt.Sprint(^bytes[i])
	}

	return strings.Join(list, `.`)
}

// AddToIPIP4 add count IPs to IP
func AddToIPIP4(startIP netip.Addr, add int32) (addedIP netip.Addr, err error) {
	if !startIP.Next().IsValid() {
		err = fmt.Errorf("ip %v is already max", startIP)
		return
	}
	bytes := startIP.As4()
	slice := bytes[:]
	ipValue := binary.BigEndian.Uint32(slice)
	ipValue += uint32(add)
	addedIP, err = int2ip(ipValue)
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

// InAddrArpaIP4 get the InAddrArpaIP4 version of an IP
func InAddrArpaIP4(ip netip.Addr) string {
	ipStr := ip.String()
	parts := strings.Split(ipStr, `.`)

	reverse(parts)

	return strings.Join(parts, ".")
}

// IPToHexStr convert an IP4 address to a hex string
func IPToHexStr(ip netip.Addr) string {
	bytes := ip.As4()
	slice := bytes[:]
	ipValue := binary.BigEndian.Uint32(slice)

	var ipInt = make(net.IP, 4)
	binary.BigEndian.PutUint32(ipInt, ipValue)

	return fmt.Sprintf("0x%X", ipValue)
}
