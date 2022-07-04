package main

import (
	"bytes"
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/cmd/handler"
)

// GitCommit the git commit hash at compile time
var GitCommit string

// GitLastTag the last tag
var GitLastTag string

// GitExactTag extract tag
var GitExactTag string

// Date the compile date
var Date string

func printHelp(p *arg.Parser) {
	fmt.Println()
	var help bytes.Buffer
	p.WriteHelp(&help)
	fmt.Println(help.String())
}

func main() {
	args.InitializeCompletion()

	var args args.Args
	arg.MustParse(&args)

	if args.Subnet != nil {
		if args.Subnet.SubnetDivide != nil {
			handler.SubnetDivide(args.Subnet.SubnetDivide.IP, uint8(args.Subnet.SubnetDivide.Mask), uint8(args.Subnet.SubnetDivide.SubMask))
		}
	}
}
