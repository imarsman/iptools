package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/alexflint/go-arg"
	"inet.af/netaddr"
)

// https://www.calculator.net/ip-subnet-calculator.html?cclass=any&csubnet=16&cip=10.0.0.0&ctype=ipv4&printit=0&x=43&y=21
// https://www.calculator.net/ip-subnet-calculator.html
type subnet struct {
	// ipSet []netaddr.IPSet

	prefix netaddr.IPPrefix
	hosts  uint64
	ipAddr netaddr.IP
	ranges []netaddr.IPRange

	// ipAddr netaddr.IP
	// networkAddr   netaddr.IP
	// broadcastAddr netaddr.IP
	// hosts         int
	// subnetMask    netaddr.IP
	// wildcardMask  netaddr.IP
	// ipClass       string
	// cidrNotation  string
	// isPublic      bool
	// inAddrArpa    string
}

func newSubnet(pfx netaddr.IPPrefix) (*subnet, error) {
	subnet := new(subnet)
	subnet.prefix = pfx

	if subnet.prefix.IP().BitLen() > 64 {
		return nil, errors.New("subnet too large")
	}

	subnet.hosts = uint64(math.Pow(2, float64(subnet.prefix.Bits())))

	return subnet, nil
}

func (s *subnet) partialBits() uint8 {
	r := s.prefix.Bits() % 8

	return r
}

func (s *subnet) partialRemainderBits() uint8 {
	r := s.partialBits()
	if r == 0 {
		return 0
	}

	if s.partialBits() == 0 {
		return 0
	}

	return 8 - s.partialBits()
}

// How many equal sized subnets can prefix be split into?
func (s *subnet) subnetSplit() uint8 {
	return uint8(math.Pow(2, float64(s.partialBits())))
}

// SubnetCmd arg to get subnet information
type SubnetCmd struct {
	Value string `arg:"positional"`
}

type args struct {
	Subnet  *SubnetCmd `arg:"subcommand:subnet"`
	Verbose bool       `arg:"-v"`
}

func (args) Description() string {
	return "this program does this and that"
}

func (args) Version() string {
	return "iptools 0.0.0.1"
}

func main() {
	var args args

	p := arg.MustParse(&args)
	fmt.Printf("%+v\n", args.Subnet)

	if args.Subnet.Value == "" {
		p.Fail("No subnet specified")
	}

	var b netaddr.IPSetBuilder
	pfx, err := netaddr.ParseIPPrefix(args.Subnet.Value)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.AddPrefix(pfx)
	b.Remove(netaddr.MustParseIP("10.2.3.4"))
	s, _ := b.IPSet()
	fmt.Println(s.Ranges())
	// fmt.Println(s.Prefixes())
}
