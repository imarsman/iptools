package main

import (
	"github.com/alexflint/go-arg"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/cmd/handler"
)

func main() {
	args.InitializeCompletion()

	arg.MustParse(&args.CLIArgs)

	// Inspect cli args and make calls to handlers as apppropriate
	if args.CLIArgs.Subnet != nil {
		if args.CLIArgs.Subnet.SubnetRanges != nil {
			handler.IP4SubnetRanges(
				args.CLIArgs.Subnet.SubnetRanges.IP,
				uint8(args.CLIArgs.Subnet.SubnetRanges.Bits),
				uint8(args.CLIArgs.Subnet.SubnetRanges.SecondaryBits),
			)
		}
		if args.CLIArgs.Subnet.SubnetDivide != nil {
			handler.IP4SubnetDivide(
				args.CLIArgs.Subnet.SubnetDivide.IP,
				uint8(args.CLIArgs.Subnet.SubnetDivide.Bits),
				uint8(args.CLIArgs.Subnet.SubnetDivide.SecondaryBits),
			)
		}
		if args.CLIArgs.Subnet.SubnetDescribe != nil {
			handler.IP4SubnetDescribe(
				args.CLIArgs.Subnet.SubnetDescribe.IP,
				uint8(args.CLIArgs.Subnet.SubnetDescribe.Bits),
			)
		}
	}
}
