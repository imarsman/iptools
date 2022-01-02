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

	// pfx, err := netaddr.ParseIPPrefix("192.168.0.1/16")
	// is.NoErr(err)

	s, err := subnet.NewSubnet("192.168.0.1/16")
	is.NoErr(err)

	t.Log("hosts", s.DivisionHosts)

	t.Log("ip", s.Prefix.IP())
	t.Log("range", s.Prefix.Range())
	t.Log("bits", s.Prefix.Bits())
	t.Log("ipnet", s.Prefix.IPNet())
	t.Log("bitlen", s.Prefix.IP().BitLen())
	t.Log("mask", s.Prefix.IPNet().Mask)
	t.Log("single IP", s.Prefix.IsSingleIP())
	t.Log("hosts", s.DivisionHosts)
	t.Log("subnetsize", s.TotalHosts)
	t.Log("equal subnets", s.EqualRanges())

	var b netaddr.IPSetBuilder
	b.AddPrefix(s.Prefix)
	ipSet, _ := b.IPSet()
	t.Log(ipSet.Ranges())

	is.True(true == true)
}

func TestBits(t *testing.T) {
	is := is.New(t)

	prefixes := []string{"99.236.32.255/21", "223.255.89.0/24", "99.236.255.255/21"}

	for _, p := range prefixes {
		// pfx, err := netaddr.ParseIPPrefix(p)
		s, err := subnet.NewSubnet(p)
		is.NoErr(err)
		t.Log("valid", s.Prefix.Valid())

		t.Log("hosts", s.DivisionHosts)
		t.Log("prefix", s.Prefix.String())
		t.Log("active byte", s.ClassByte())
		t.Log("ip range", s.Prefix.Range())
		t.Log("subnet usable ip range", subnet.UsableRange(s.Prefix.Range()))
		t.Log("partial bits", s.ClassPartialBits())
		t.Log("partial remainder bits", s.ClassHostBits())
		t.Log("prefix bits", s.Prefix.Bits())
		t.Log("hosts", s.DivisionHosts)
		t.Log("subnetsize", s.TotalHosts)
		t.Log("equal subnets", s.EqualRanges())
		t.Log("subnet 3")
		t.Log(s.Divisions)
		bytes, err := s.YAML()
		is.NoErr(err)
		t.Log(string(bytes))
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
