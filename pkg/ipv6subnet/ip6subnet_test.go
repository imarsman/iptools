package ipv6subnet

import (
	"testing"

	"github.com/imarsman/iptools/pkg/ipv6subnet/ip6util"
	"github.com/matryer/is"
)

// Sample IP addresses
// 2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b

func TestNewSubnet(t *testing.T) {
	is := is.New(t)

	addr, err := ip6util.RandomAddrGlobalUnicast()
	is.NoErr(err)
	s, err := NewFromIPAndBits(addr.StringExpanded(), 64)
	is.NoErr(err)
	t.Log("First in subnet", s.First().StringExpanded())
	t.Log("Last in subnet", s.Last().StringExpanded())
	t.Log(s.prefix.Masked())
	t.Log(s.prefix)
	t.Log(s.prefix.Addr().StringExpanded())

	t.Log("subnet", s.SubnetString())
	t.Log("interface", s.InterfaceString())
	t.Log("is global unicast", s.Addr().IsGlobalUnicast())
	t.Log("Address type", ip6util.AddressTypeName(s.Addr()))
	t.Log("Address prefix", s.TypePrefix().Masked().String())
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
