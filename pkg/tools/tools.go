package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"inet.af/netaddr"
)

// BytesForMask get bytes from hex string for CIDR base address
func BytesForMask(encodedMask string) (bytes []byte, err error) {
	// get base64 bytes by converting from encoded string
	bytes, err = base64.StdEncoding.DecodeString(encodedMask)
	if err != nil {
		return
	}
	if len(bytes) != net.IPv4len {
		err = fmt.Errorf("did not decode %s to %d bytes", encodedMask, net.IPv4len)
		return
	}
	return
}

// CIDR get CIDR notation string for a range
func CIDR(maskLen int) (cidrString string, err error) {
	cidrString = fmt.Sprintf("250.250.250.250/%d", maskLen)
	p, err := netaddr.ParseIPPrefix(cidrString)
	cidrString = fmt.Sprintf("%s/%d", p.Masked(), p.Bits())

	// ipMask := net.CIDRMask(maskLen, 32)
	// bytes, err := hex.DecodeString(ipMask.String())
	// if err != nil {
	// 	return
	// }

	// cidrString = fmt.Sprintf("%d.%d.%d.%d/%d", bytes[0], bytes[1], bytes[2], bytes[3], maskLen)

	return
}

// DecodeMaskBase64 get CIDR value from encoded hex string
func DecodeMaskBase64(encodedMask string, verbose bool) (cidr int, err error) {
	bytes, err := BytesForMask(encodedMask)
	if err != nil {
		return
	}
	// using net package
	ip := make(net.IP, net.IPv4len)
	ipMask := net.IPMask(ip)
	copy(ip[:], bytes)

	// uses netaddr
	// bytes := [4]byte{}
	// copy(bytes[:], b)
	// netAddrIP := netaddr.IPFrom4(bytes)
	// ipMask := net.IPMask(netAddrIP.IPAddr().IP.To4())
	cidr, _ = ipMask.Size()
	if verbose {
		cidrString := fmt.Sprintf("%d.%d.%d.%d/%d", bytes[0], bytes[1], bytes[2], bytes[3], cidr)
		fmt.Println(cidrString)
		fmt.Println("decoded bytes", bytes, "cidr mask", cidr)
	}

	return
}

// DecodeCIDRIP get bitlength for a subnet mask from JSON encoded int
// This uses the json library
func DecodeCIDRIP(encodedMask string, verbose bool) (cidr int, err error) {
	if strings.HasPrefix(encodedMask, `"`) {
		encodedMask = strings.ReplaceAll(encodedMask, `"`, "")
		// encodedMask = fmt.Sprintf(`"%s"`, encodedMask)
	}
	var bytes = []byte{}
	err = json.Unmarshal([]byte(encodedMask), &bytes)
	if err != nil {
		return
	}
	if len(bytes) != net.IPv4len {
		err = fmt.Errorf("did not decode %s to %d bytes", encodedMask, net.IPv4len)
		return
	}

	netAddrIP := make(net.IP, net.IPv4len)
	copy(netAddrIP[:], bytes)

	ipMask := net.IPMask(netAddrIP)
	cidr, _ = ipMask.Size()

	return
}

// IPsForCIDR get IP range from CIDR string
func IPsForCIDR(cidr string) (rng netaddr.IPRange, err error) {
	ip, _, err := net.ParseCIDR(cidr)
	end, err := netaddr.ParseIP("255.255.255.255")
	if err != nil {
		return
	}
	ipNetaddr, _ := netaddr.FromStdIP(ip)
	rng = netaddr.IPRangeFrom(ipNetaddr, end)

	return
}
