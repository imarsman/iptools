package ipv6subnet

import (
	"net/netip"
	"testing"

	"github.com/imarsman/iptools/pkg/ipv6subnet/util"
	"github.com/matryer/is"
)

// Sample IP addresses
// 2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b

func TestNewSubnet(t *testing.T) {
	is := is.New(t)

	s, err := NewFromIPAndBits("2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b", 64)
	is.NoErr(err)
	t.Log(s.prefix.Masked())
	t.Log(s.prefix)
	t.Log(s.prefix.Addr().StringExpanded())

	t.Log("subnet", s.SubnetString())
	t.Log("interface", s.InterfaceString())
	t.Log("prefix", s.PrefixString())
	t.Log("is global unicast", s.Addr().IsGlobalUnicast())
}

func TestNetworks(t *testing.T) {
	// is := is.New(t)
	// p, err := netip.ParsePrefix("192.24.12.0/18")
	// is.NoErr(err)
	// t.Log(p.Addr())
	// t.Log(p.Masked())
	// s, err := NewFromPrefix(p.Masked().String())
	// networks, err := s.IPRanges()
	// is.NoErr(err)
	// t.Log(networks)
}

func TestDifferingSubnets(t *testing.T) {
	// is := is.New(t)
	// s, err := NewFromPrefix("10.0.0.0/23")
	// is.NoErr(err)
	// t.Log("starting from", s)
	// t.Log("hosts per network", s.Hosts())
	// t.Log("network count", s.Networks())
	// networks, err := s.IPRanges()
	// is.NoErr(err)
	// t.Log(networks)
	// s, err = NewFromPrefix("10.0.0.0/24")
	// t.Log("starting from", s)
	// t.Log("hosts per network", s.Hosts())
	// t.Log("network count", s.Networks())
	// networks, err = s.IPRanges()
	// for _, r := range networks {
	// 	t.Log(r.String())
	// }

	// // t.Log(networks)
	// is.NoErr(err)
}

type parentChild struct {
	subnetIP string
	parent   int
	child    int
}

func TestSecondarySubnets(t *testing.T) {
	// is := is.New(t)

	// parentChildSet := []parentChild{}
	// subnetIP := "10.0.0.0"
	// parentChildSet = append(parentChildSet, parentChild{subnetIP: subnetIP, parent: 23, child: 23})
	// parentChildSet = append(parentChildSet, parentChild{subnetIP: subnetIP, parent: 24, child: 24})
	// parentChildSet = append(parentChildSet, parentChild{subnetIP: subnetIP, parent: 23, child: 24})

	// for _, item := range parentChildSet {
	// 	prefix := fmt.Sprintf("%s/%d", item.subnetIP, item.parent)
	// 	subnet, err := NewFromPrefix(prefix)

	// 	prefix = fmt.Sprintf("%s/%d", item.subnetIP, item.child)
	// 	secondarySubnet, err := NewFromPrefix(prefix)
	// 	is.NoErr(err)
	// 	t.Log("starting from", subnet)
	// 	t.Log("child subnet", secondarySubnet)
	// 	t.Log("hosts per network", subnet.Hosts())
	// 	t.Log("child subnet hosts per network", secondarySubnet.Hosts())
	// 	t.Log("network count", subnet.Networks())
	// 	start := time.Now()
	// 	networks, err := subnet.SecondaryIPRanges(secondarySubnet)
	// 	is.NoErr(err)
	// 	t.Log("run took", time.Since(start))
	// 	t.Log("total networks", len(networks))
	// 	// t.Log("networks", networks)

	// 	for _, r := range networks {
	// 		t.Log(r.String())
	// 	}
	// }
}

func TestBitString(t *testing.T) {
	is := is.New(t)
	ip, err := netip.ParseAddr("127.0.0.1")
	is.NoErr(err)
	bitStr := util.BitStr4(ip, ".")
	t.Log(bitStr)
}

func TestIPString(t *testing.T) {
	is := is.New(t)
	start := "127.0.0.1"
	ip, err := netip.ParseAddr(start)
	is.NoErr(err)

	bitStr := util.BitStr4(ip, ".")
	t.Log("bitStr for 127.0.0.1", bitStr)

	bytes, err := util.BinaryIP4StrToBytes(bitStr)
	is.NoErr(err)

	list := make([]byte, 0, 0)

	for _, b := range bytes {
		list = append(list, b)
	}
	t.Logf("Started with %s got bytes %v", start, list)

	start = "99.236.32.0"
	bitStr = "01100011.11101100.00100000.00000000"
	t.Log("bitStr for 99.236.32.0", bitStr)
	bytes, err = util.BinaryIP4StrToBytes(bitStr)
	is.NoErr(err)

	list = make([]byte, 0, 0)

	for _, b := range bytes {
		list = append(list, b)
	}
	t.Logf("Started with %s got bytes %v", start, list)
}

// TestIPStringSplit test split of IP4 binary string to 4 byte slice
func TestIPStringSplit(t *testing.T) {
	is := is.New(t)

	list := []string{"01100011.11101100.00100000.00000000", "01100011111011000010000000000000"}

	for _, bitStr := range list {
		// bitStr := "01100011.11101100.00100000.00000000"
		bytes, err := util.BinaryIP4StrToBytes(bitStr)
		is.NoErr(err)
		t.Log(bytes)
	}
}

// go test -bench=. -benchmem
func BenchmarkNewSubnet(b *testing.B) {
	// is := is.New(b)
	// s, err := NewFromIPAndBits("10.32.0.0", 28)
	// is.NoErr(err)
	// s.IPRanges()
}

func BenchmarkSubnetSplit(b *testing.B) {
	// is := is.New(b)
	// s, err := NewFromIPAndBits("10.32.0.0", 28)
	// is.NoErr(err)
	// subnets, err := s.subnets(s)
	// is.NoErr(err)
	// is.True(len(subnets) == 16)
	// s.IPRanges()
}
