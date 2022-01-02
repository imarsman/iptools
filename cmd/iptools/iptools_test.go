package main

import (
	"testing"

	"github.com/imarsman/iptools/pkg/subnet"
	"github.com/matryer/is"
	"inet.af/netaddr"
)

//                Tests and benchmarks
// -----------------------------------------------------
// benchmark
//   go test -run=XXX -bench=. -benchmem
// Get allocation information and pipe to less
//   go build -gcflags '-m -m' ./*.go 2>&1 |less
// Run all tests
//   go test -v
// Run one test and do allocation profiling
//   go test -run=XXX -bench=IterativeISOTimestampLong -gcflags '-m' 2>&1 |less
// Run a specific test by function name pattern
//  go test -run=TestParsISOTimestamp
//
//  go test -run=XXX -bench=.
//  go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out
//  go tool pprof -http=:8080 memprofile.out
//  go tool pprof -http=:8080 cpuprofile.out

const (
	bechmarkBytesPerOp int64 = 10
)

func TestRange(t *testing.T) {
	is := is.New(t)

	pfx, err := netaddr.ParseIPPrefix("192.168.0.1/16")
	is.NoErr(err)

	subnet, err := subnet.NewSubnet(pfx)
	is.NoErr(err)
	subnet.Prefix = pfx

	t.Log("hosts", subnet.Hosts)

	t.Log("ip", pfx.IP())
	t.Log("range", pfx.Range())
	t.Log("bits", pfx.Bits())
	t.Log("ipnet", pfx.IPNet())
	t.Log("bitlen", pfx.IP().BitLen())
	t.Log("mask", pfx.IPNet().Mask)
	t.Log("single IP", pfx.IsSingleIP())
	t.Log("hosts", subnet.Hosts)
	t.Log("subnetsize", subnet.SubnetSize)
	t.Log("equal subnets", subnet.EqualSubnets())

	var b netaddr.IPSetBuilder
	b.AddPrefix(pfx)
	s, _ := b.IPSet()
	t.Log(s.Ranges())

	is.True(true == true)
}

func TestBits(t *testing.T) {
	is := is.New(t)

	pfx, err := netaddr.ParseIPPrefix("99.236.0.0/21")
	subnet, err := subnet.NewSubnet(pfx)
	is.NoErr(err)

	t.Log("prefix", pfx.String())
	t.Log("active byte", subnet.ClassByte())
	t.Log("ip range", pfx.Range())
	t.Log("subnet usable ip range", UsableRange(subnet.Prefix.Range()))
	t.Log("partial bits", subnet.ClassPartialBits())
	t.Log("partial remainder bits", subnet.ClassHostBits())
	t.Log("prefix bits", pfx.Bits())
	t.Log("hosts", subnet.Hosts)
	t.Log("subnetsize", subnet.SubnetSize)
	t.Log("equal subnets", subnet.EqualSubnets())
	t.Log("subnet 3")
	t.Log(subnet.SubnetDivisions())
	t.Log(subnet.SubnetDivisions())
}

func BenchmarkPathParts(b *testing.B) {
	is := is.New(b)

	b.SetBytes(bechmarkBytesPerOp)
	b.ReportAllocs()
	b.SetParallelism(30)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			is.True(1 == 1)
		}
	})

}
