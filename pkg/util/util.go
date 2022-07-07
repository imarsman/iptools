package util

import (
	"encoding/binary"
	"fmt"
	"net"
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

// IntToHexStr get hex string from int32
func IntToHexStr(hex uint32) string {
	return fmt.Sprintf("%x", hex)
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

func InAddrArpa(ip netaddr.IP) string {
	ipStr := ip.String()
	parts := strings.Split(ipStr, `.`)

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ".")
}

func IPToHexStr(ip netaddr.IP) string {
	bytes := ip.As4()
	slice := bytes[:]
	ipValue := binary.BigEndian.Uint32(slice)

	var ipInt = make(net.IP, 4)
	binary.BigEndian.PutUint32(ipInt, ipValue)
	return IntToHexStr(ipValue)
}
