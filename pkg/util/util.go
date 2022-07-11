package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"inet.af/netaddr"
)

// BitStr4 get bit string for IPV4 IP
func BitStr4(ip netaddr.IP, separator string) string {
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

// Int32ToHexStr get hex string from int32
func Int32ToHexStr(int32In uint32) string {
	return fmt.Sprintf("%x", int32In)
}

// https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func int2ip(ipInt uint32) (netaddr.IP, bool) {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	addr, ok := netaddr.FromStdIP(ip)

	return addr, ok
}

// WildCardMask get mask bits available for addressing
func WildCardMask(ip netaddr.IP) string {
	bytes := ip.As4()
	var list = make([]string, 4, 4)
	for i, b := range bytes {
		if b == 255 {
			list[i] = fmt.Sprint(0)
			continue
		}
		list[i] = fmt.Sprint(250 - b)
	}

	return strings.Join(list, `.`)
}

// AddToIP add count IPs to IP
func AddToIP(startIP netaddr.IP, add int32) (newIP netaddr.IP, err error) {
	if !startIP.Next().IsValid() {
		err = fmt.Errorf("ip %v is already max", startIP)
		return
	}
	bytes := startIP.As4()
	slice := bytes[:]
	ipValue := binary.BigEndian.Uint32(slice)
	ipValue += uint32(add)
	newIP, ok := int2ip(ipValue)
	if !ok {
		err = fmt.Errorf("problem after adding %d to IP %v", add, startIP)
		return
	}
	return newIP, nil
}

// For fun with generics
func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// InAddrArpa get the InAddrArpa version of an IP
func InAddrArpa(ip netaddr.IP) string {
	ipStr := ip.String()
	parts := strings.Split(ipStr, `.`)

	reverse(parts)

	return strings.Join(parts, ".")
}

// IPToHexStr convert an IP4 address to a hex string
func IPToHexStr(ip netaddr.IP) string {
	bytes := ip.As4()
	slice := bytes[:]
	ipValue := binary.BigEndian.Uint32(slice)

	var ipInt = make(net.IP, 4)
	binary.BigEndian.PutUint32(ipInt, ipValue)

	return Int32ToHexStr(ipValue)
}
