package ipv6util

import (
	"fmt"
	"net/netip"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func TestMakeMacAddress(t *testing.T) {
	// bytes, err := randomMacAddress()
	is := is.New(t)
	// is.NoErr(err)

	bytes, err := randomMacBytesForInterface()
	is.NoErr(err)
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])
	t.Log("MAC address", macAddress)
	t.Log(bytes2MacAddr(bytes))
	addr, err := randomLinkLocal()
	is.NoErr(err)
	t.Log("link local address", addr)
	addr, err = randomGlobalUnicast()
	is.NoErr(err)
	t.Log("global unicast address", addr)
}

func TestRandomSubnet(t *testing.T) {
	randSubnet := RandomSubnet()

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
