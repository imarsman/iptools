package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"inet.af/netaddr"
)

// SubnetCmd arg to get subnet information
type SubnetCmd struct {
	Value string `arg:"positional"`
}

// args cli args
type args struct {
	Subnet  *SubnetCmd `arg:"subcommand:subnet"`
	Verbose bool       `arg:"-v"`
}

func (args) Description() string {
	return "a tool to deal with IP subnets and splitting subnets"
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
