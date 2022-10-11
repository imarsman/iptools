package ipv6

import (
	"net/netip"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func TestSubnet(t *testing.T) {
	is := is.New(t)

	addr, err := RandomAddrGlobalUnicast()
	prefix := netip.PrefixFrom(addr, 64)
	is.NoErr(err)
	t.Log("First in subnet", First(addr).StringExpanded())
	t.Log("Last in subnet", Last(addr).StringExpanded())
	t.Log(prefix.Masked())
	t.Log(prefix)
	t.Log(prefix.Addr().StringExpanded())

	t.Log("subnet", SubnetString(addr))
	t.Log("interface", InterfaceString(addr))
	t.Log("is global unicast", addr.IsGlobalUnicast())
	t.Log("Address type", AddrTypeName(addr))
	t.Log("Address prefix", AddrTypePrefix(addr).Masked().String())
}

func TestRandomGlobalUnicast(t *testing.T) {
	is := is.New(t)
	addr, err := RandomAddrGlobalUnicast()
	is.NoErr(err)
	t.Log(AddrTypeName(addr))
}

func TestRandomLinkLocal(t *testing.T) {
	is := is.New(t)
	addr, err := RandomAddrLinkLocal()
	is.NoErr(err)
	t.Log(AddrTypeName(addr))
}

func TestPrivate(t *testing.T) {
	is := is.New(t)
	addr, err := RandomAddrPrivate()
	is.NoErr(err)
	t.Log(AddrTypeName(addr))
}

func TestMulticast(t *testing.T) {
	is := is.New(t)
	addr, err := RandomAddrMulticast()
	is.NoErr(err)
	t.Log(AddrTypeName(addr))
}

func TestInterfaceLocalMulticast(t *testing.T) {
	is := is.New(t)
	addr, err := RandomAddrInterfaceLocalMulticast()
	is.NoErr(err)
	t.Log(AddrTypeName(addr))
}

func TestLinkLocalMulticast(t *testing.T) {
	is := is.New(t)
	addr, err := RandomAddrLinkLocalMulticast()
	is.NoErr(err)
	t.Log(AddrTypeName(addr))
}

func TestMakeMacAddress(t *testing.T) {
	// bytes, err := randomMacAddress()
	is := is.New(t)
	// is.NoErr(err)

	bytes, err := randomMacBytesForInterface()
	is.NoErr(err)
	macAddress := bytes2MacAddr(bytes)
	// macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])
	t.Log("MAC address", macAddress)
	t.Log(bytes2MacAddr(bytes))
}

func TestRandomSubnet(t *testing.T) {
	randSubnet := RandomSubnetID()

	t.Log(strconv.FormatInt(int64(randSubnet), 16))
}

func TestBits(t *testing.T) {
	is := is.New(t)
	ip, err := netip.ParseAddr("3501:db8:cafe:dcb2:f945:2aff:feee:f0d6")
	t.Log(ip.StringExpanded())
	is.NoErr(err)
	bytes := ip.As16()
	var toConvert [8]byte
	copy(toConvert[:], bytes[0:7])

	converted := GlobalID(ip)

	t.Log(converted)
}
