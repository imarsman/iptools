package util

import (
	"fmt"

	"inet.af/netaddr"
)

// BitStr4 get bit string for IPV4 IP
func BitStr4(ip netaddr.IP) string {
	bytes := ip.As4()
	output := ""
	str := fmt.Sprintf("%08b", bytes[0])
	output = output + str
	str = fmt.Sprintf("%08b", bytes[1])
	output = output + "." + str
	str = fmt.Sprintf("%08b", bytes[2])
	output = output + "." + str
	str = fmt.Sprintf("%08b", bytes[3])
	output = output + "." + str

	return output
}
