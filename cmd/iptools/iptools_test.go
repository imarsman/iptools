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

	s, err := subnet.NewSubnet(pfx)
	is.NoErr(err)
	s.Prefix = pfx

	t.Log("hosts", s.Hosts)

	t.Log("ip", pfx.IP())
	t.Log("range", pfx.Range())
	t.Log("bits", pfx.Bits())
	t.Log("ipnet", pfx.IPNet())
	t.Log("bitlen", pfx.IP().BitLen())
	t.Log("mask", pfx.IPNet().Mask)
	t.Log("single IP", pfx.IsSingleIP())
	t.Log("hosts", s.Hosts)
	t.Log("subnetsize", s.SubnetSize)
	t.Log("equal subnets", s.EqualSubnets())

	var b netaddr.IPSetBuilder
	b.AddPrefix(pfx)
	ipSet, _ := b.IPSet()
	t.Log(ipSet.Ranges())

	is.True(true == true)
}

func TestBits(t *testing.T) {
	is := is.New(t)

	prefixes := []string{"99.236.0.0/21", "223.255.89.0/24"}

	for _, p := range prefixes {
		pfx, err := netaddr.ParseIPPrefix(p)
		s, err := subnet.NewSubnet(pfx)
		is.NoErr(err)

		t.Log("prefix", pfx.String())
		t.Log("active byte", s.ClassByte())
		t.Log("ip range", pfx.Range())
		t.Log("subnet usable ip range", subnet.UsableRange(s.Prefix.Range()))
		t.Log("partial bits", s.ClassPartialBits())
		t.Log("partial remainder bits", s.ClassHostBits())
		t.Log("prefix bits", pfx.Bits())
		t.Log("hosts", s.Hosts)
		t.Log("subnetsize", s.SubnetSize)
		t.Log("equal subnets", s.EqualSubnets())
		t.Log("subnet 3")
		t.Log(s.SubnetDivisions())
		t.Log(s.SubnetDivisions())
	}
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
