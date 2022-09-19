package util

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func TestMakeMacAddress(t *testing.T) {
	bytes, err := makeMacAddress()
	is := is.New(t)
	is.NoErr(err)

	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])
	t.Log("MAC address", macAddress)
	macAddrBytes, err := bytes2MacAddrBytes(bytes)
	is.NoErr(err)
	t.Log(bytesToMacAddr(macAddrBytes))
	addr, err := mac2LinkLocal(bytesToMacAddr(macAddrBytes))
	is.NoErr(err)
	t.Log("link local address", addr)
	addr, err = mac2GlobalUnicast(bytesToMacAddr(macAddrBytes))
	is.NoErr(err)
	t.Log("global unicast address", addr)
}

func TestRandomSubnet(t *testing.T) {
	randSubnet := RandomSubnet()

	t.Log(strconv.FormatInt(int64(randSubnet), 16))
}
