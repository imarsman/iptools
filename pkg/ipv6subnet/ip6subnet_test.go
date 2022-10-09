package ipv6subnet

import (
	"net/netip"
	"testing"

	"github.com/imarsman/iptools/pkg/ipv6subnet/ipv6util"
	ip6util "github.com/imarsman/iptools/pkg/ipv6subnet/ipv6util"
	"github.com/matryer/is"
)

// Sample IP addresses
// 2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b

func TestNewSubnet(t *testing.T) {
	is := is.New(t)

	addr, err := ip6util.RandomAddrGlobalUnicast()
	prefix := netip.PrefixFrom(addr, 64)
	is.NoErr(err)
	// s, err := NewFromIPAndBits(addr.StringExpanded(), 64)
	// is.NoErr(err)
	t.Log("First in subnet", ipv6util.First(addr).StringExpanded())
	t.Log("Last in subnet", ipv6util.Last(addr).StringExpanded())
	t.Log(prefix.Masked())
	t.Log(prefix)
	t.Log(prefix.Addr().StringExpanded())

	t.Log("subnet", ipv6util.SubnetString(addr))
	t.Log("interface", ipv6util.InterfaceString(addr))
	t.Log("is global unicast", addr.IsGlobalUnicast())
	t.Log("Address type", ip6util.AddressTypeName(addr))
	t.Log("Address prefix", ipv6util.TypePrefix(addr).Masked().String())
}

func TestRandomGlobalUnicast(t *testing.T) {
	is := is.New(t)
	addr, err := ip6util.RandomAddrGlobalUnicast()
	is.NoErr(err)
	t.Log(ip6util.AddressTypeName(addr))
}

func TestRandomLinkLocal(t *testing.T) {
	is := is.New(t)
	addr, err := ip6util.RandomAddrLinkLocal()
	is.NoErr(err)
	t.Log(ip6util.AddressTypeName(addr))
}

type parentChild struct {
	subnetIP string
	parent   int
	child    int
}
