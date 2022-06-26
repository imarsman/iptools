package subnet

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewSubnet(t *testing.T) {
	is := is.New(t)
	s, err := NewSubnet("200.200.200.200/28", true)
	is.NoErr(err)
	t.Log(s.Prefix)
	t.Log("subnet max bits", s.MaxBitsForClass())
	t.Log("block size", s.BlockSize())
	t.Log("subnets", s.Subnets())
	t.Log("hosts", s.Hosts())
	t.Log("total hosts", s.TotalHosts())
	t.Log("usable hosts", s.UsableHosts())

	r, err := s.Range()
	is.NoErr(err)
	t.Log("subnet range", r)

	ips, err := s.IPs()
	t.Log("ips", ips)
	is.NoErr(err)

	usableIps, err := s.UsableIPs()
	is.NoErr(err)
	t.Log("usable ips", usableIps)
	s.Blocks()
}

// go test -bench=. -benchmem
func BenchmarkBlocks(b *testing.B) {
	is := is.New(b)
	s, err := NewSubnet("200.200.200.200/28", true)
	is.NoErr(err)
	s.Blocks()
}
