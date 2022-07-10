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

	// var Args args.Args
	arg.MustParse(&args.CLIArgs)

	if args.CLIArgs.Subnet != nil {
		if args.CLIArgs.Subnet.SubnetRanges != nil {
			handler.SubnetRanges(
				args.CLIArgs.Subnet.SubnetRanges.IP,
				uint8(args.CLIArgs.Subnet.SubnetRanges.Mask),
				uint8(args.CLIArgs.Subnet.SubnetRanges.SubMask),
			)
		}
		if args.CLIArgs.Subnet.SubnetDivide != nil {
			handler.SubnetDivide(
				args.CLIArgs.Subnet.SubnetDivide.IP,
				uint8(args.CLIArgs.Subnet.SubnetDivide.Mask),
				uint8(args.CLIArgs.Subnet.SubnetDivide.SubMask),
			)
		}
		if args.CLIArgs.Subnet.SubnetDescribe != nil {
			handler.SubnetDescribe(
				args.CLIArgs.Subnet.SubnetDescribe.IP,
				uint8(args.CLIArgs.Subnet.SubnetDescribe.Mask),
			)
		}
	}
}
