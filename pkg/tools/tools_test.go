package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/matryer/is"
	"inet.af/netaddr"
)

func expandInterfaceToMatch(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = expandInterfaceToMatch(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = expandInterfaceToMatch(v)
		}
	}
	return i
}

func TestEncodedMask(t *testing.T) {
	is := is.New(t)

	// "fffff000"
	mask := net.CIDRMask(20, 32)
	t.Log("mask", mask)
	// Marshal takes a byte slice
	// The fact we can send a net.IPMask, which is a type aliases to []byte
	// indicates that we are marshalling a byte slice, not the string output
	// of the CIDR mask (a hexidecimal string)
	new, err := json.Marshal(mask)
	is.NoErr(err)
	t.Log(string(new))
	// Another demo of marshalling bytes
	bytes := []byte(mask)
	new, err = json.Marshal([]byte(bytes))
	is.NoErr(err)
	t.Log(string(new))

	// stringValue := strings.ReplaceAll(string(new), `"`, "")
	stringValue := string(new)
	var obj interface{}
	t.Log([]byte(stringValue))
	err = json.Unmarshal([]byte(stringValue), &obj)
	is.NoErr(err)
	obj = expandInterfaceToMatch(obj)
	t.Log("unmarshalled bytes", obj)

	hexBytes, err := json.Marshal([]byte("fffff000"))
	is.NoErr(err)
	t.Log("hex string", string(hexBytes))

}

func TestEncodeIP(t *testing.T) {
	is := is.New(t)
	ip := `{ "ip": "255.255.240.000" }`

	t.Log("unmarshalling", ip)
	var obj = make(map[string]string)
	err := json.Unmarshal([]byte(ip), &obj)
	is.NoErr(err)
	t.Log("ip key from resulting map", obj[`ip`])
}

// TestDecodeBytes test base64 decoding
func TestDecodeBytes(t *testing.T) {
	is := is.New(t)
	encoded := `///wAA==`

	cidr, err := DecodeMaskBase64(encoded, true)
	is.NoErr(err)

	t.Logf("cidr %d\n", cidr)
}

func TestFromHex(t *testing.T) {
	is := is.New(t)
	b, err := base64.StdEncoding.DecodeString(`///wAA==`)
	is.NoErr(err)
	t.Log("hex to bytes", b)
	bytes := [4]byte{}
	copy(bytes[:], b)
	netAddrIP := netaddr.IPFrom4(bytes)
	ipMask := net.IPMask(netAddrIP.IPAddr().IP.To4())
	cidr, _ := ipMask.Size()
	t.Log("cidr", cidr)
}

func TestPrefix(t *testing.T) {
	is := is.New(t)
	p, err := netaddr.ParseIPPrefix("250.250.250.250/26")
	is.NoErr(err)
	t.Logf("Mask IP %s", p.Masked())
	t.Logf("bits %d", p.Bits())
}

func TestIPPrefix(t *testing.T) {
	is := is.New(t)
	ip, err := netaddr.ParseIP("255.255.255.255")
	is.NoErr(err)
	ipPrefix := netaddr.IPPrefixFrom(ip, 23)
	t.Log("range", ipPrefix.Range())
	// ip2, err := netaddr.ParseIP("255.255.255.255")
	// t.Log("start", ipPrefix.Masked())

	// https://www.calculator.net/ip-subnet-calculator.html?cclass=any&csubnet=23&cip=255.255.254.0&ctype=ipv4&printit=0&x=79&y=18
	// 2^9 = 512
	// http://www.sput.nl/internet/netmask-table.html
	// https://www.ittsystems.com/introduction-to-subnetting/#wbounce-modal
	// 2^required subnets size - reference block size
	// so /23 is 2^24-23=1
	for i := 0; i < 8; i += 2 {
		ip2, err := netaddr.ParseIP("255.255.255.255")
		is.NoErr(err)
		testIP, err := ip2.Prefix(23 + uint8(i))
		fmt.Println("next", testIP.Masked())
	}

	// // ipPrefix2 := netaddr.IPPrefixFrom(ip2, 24)
	// testIP, err := ipPrefix2.IP().Prefix(26)
	// is.NoErr(err)
	// t.Log("added", testIP.Masked())
	// t.Log("range", testIP.Range())
	// t.Log(ipPrefix2.Masked())

	t.Log(ipPrefix.Masked())
	t.Log(ipPrefix.IsValid())

	var b netaddr.IPSetBuilder
	b.AddPrefix(ipPrefix)

	s, _ := b.IPSet()

	fmt.Println("Ranges:")
	for _, r := range s.Ranges() {
		fmt.Printf("  %s - %s\n", r.From(), r.To())
	}

	fmt.Println("Prefixes:")
	for _, p := range s.Prefixes() {
		fmt.Printf("  %s\n", p)
	}
}
