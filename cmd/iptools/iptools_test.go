package main

import (
	"bytes"
	"encoding/binary"
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

	s, err := subnet.NewFromPrefix("192.168.0.1/16")
	is.NoErr(err)

	// t.Log("hosts", s.DivisionIncrement)

	t.Log("ip", s.Prefix().IP())
	t.Log("range", s.Prefix().Range())
	t.Log("bits", s.Prefix().Bits())
	t.Log("ipnet", s.Prefix().IPNet())
	t.Log("bitlen", s.Prefix().IP().BitLen())
	t.Log("mask", s.Prefix().IPNet().Mask)
	t.Log("single IP", s.Prefix().IsSingleIP())
	// t.Log("hosts", s.DivisionIncrement)
	// t.Log("subnetsize", s.SubnetHosts)
	// t.Log("equal subnets", s.TotalDivisions)

	var b netaddr.IPSetBuilder
	// b.AddPrefix(s.Prefix)
	ipSet, _ := b.IPSet()
	t.Log(ipSet.Ranges())

	is.True(true == true)
}

func TestBits(t *testing.T) {
	is := is.New(t)

	prefixes := []string{"99.236.32.255/21", "99.236.32.255/22", "99.236.32.255/16", "99.236.32.255/17", "223.255.89.0/24", "99.236.255.255/21"}

	for _, p := range prefixes {
		// pfx, err := netaddr.ParseIPPrefix(p)
		s, err := subnet.NewFromPrefix(p)
		is.NoErr(err)
		t.Log("valid", s.Prefix().Valid())

		// t.Log("hosts", s.DivisionIncrement)
		t.Log("prefix", s.Prefix().String())
		t.Log("active byte", s.PrefixBits())
		t.Log("ip range", s.Prefix().Range())
		// t.Log("subnet usable ip range", subnet.UsableRange(s.Prefix.Range()))
		t.Log("partial bits", s.ClassHostBits())
		// t.Log("host bits", s.ClassHots())
		t.Log("prefix bits", s.Prefix().Bits())
		// t.Log("hosts", s.DivisionIncrement)
		// t.Log("subnet hosts", s.SubnetHosts)
		// t.Log("division increment", s.DivisionIncrement)
		// t.Log("total divisions", s.TotalDivisions)
		// t.Log("division hosts", s.DivisionIncrement)
		// t.Log(s.Divisions)
	}
}

func readInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func TestShift(t *testing.T) {
	b := []byte{1}
	t.Log("bytes", b)
	// b = subnet.ShiftLeft(b, 7)
	t.Log("bytes", b)
	t.Log(0 << 0)
	t.Log(1 << 8)
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
