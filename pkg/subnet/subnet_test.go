package subnet

import (
	"testing"

	"github.com/matryer/is"
	"inet.af/netaddr"
)

func TestNewSubnet(t *testing.T) {
	for i := 1; i <= 32; i++ {
		// subnetMask := fmt.Sprintf("255.255.255.0/%d", i)
		is := is.New(t)
		s, err := NewDefaultFromMask(uint8(i))
		is.NoErr(err)
		t.Log("masked", s.Prefix.Masked())
		t.Log("class mask", s.classMask())
		t.Log("subnet max bits", s.maxBitsForClass())
		t.Log("subnet start bits", s.startBitsForClass())
		// t.Log("block size", s.BlockSize())
		t.Log("networks", s.NetworkCount())
		// hosts per network
		// t.Log("count", float64(s.PrefixBits()), int64((math.Exp2(float64(32 - s.Prefix.Bits())))))
		t.Log("hosts per network", s.Hosts())
		t.Log("usable hosts per network", s.UsableHosts())
		t.Log("total hosts per subnet", s.Hosts()*s.NetworkCount())
		t.Log("class network prefix bits", s.ClassNetworkPrefixBits())
		t.Log("class host identifier bits", s.ClassHostItentifierBits())
		t.Log()
	}
}

func TestNetworks(t *testing.T) {
	is := is.New(t)
	p, err := netaddr.ParseIPPrefix("192.24.12.0/18")
	is.NoErr(err)
	t.Log(p.IP())
	t.Log(p.Masked())
	s, err := NewFromPrefix(p.Masked().String())
	networks, err := s.Networks()
	is.NoErr(err)
	t.Log(networks)
}

func TestDifferingSubnets(t *testing.T) {
	is := is.New(t)
	s, err := NewFromPrefix("10.0.0.0/23")
	is.NoErr(err)
	t.Log("starting from", s)
	t.Log("hosts per network", s.Hosts())
	t.Log("network count", s.NetworkCount())
	networks, err := s.Networks()
	is.NoErr(err)
	t.Log(networks)
	s, err = NewFromPrefix("10.0.0.0/24")
	t.Log("starting from", s)
	t.Log("hosts per network", s.Hosts())
	t.Log("network count", s.NetworkCount())
	networks, err = s.Networks()
	t.Log(networks)
	is.NoErr(err)
}

func TestChildSubnets(t *testing.T) {
	is := is.New(t)
	s, err := NewFromPrefix("10.0.0.0/23")
	is.NoErr(err)
	childSubnet, err := NewFromPrefix("10.0.0.0/24")
	t.Log("starting from", s)
	t.Log("child subnet", childSubnet)
	t.Log("hosts per network", s.Hosts())
	t.Log("child subnet hosts per network", childSubnet.Hosts())
	t.Log("network count", s.NetworkCount())
	networks, err := s.NetworksInSubnets(childSubnet)
	is.NoErr(err)
	t.Log("total networks", len(networks))
	t.Log(networks)
}

// go test -bench=. -benchmem
func BenchmarkBlocks(b *testing.B) {
	is := is.New(b)
	s, err := NewDefaultFromMask(28)
	is.NoErr(err)
	s.NetworkRanges()
}
