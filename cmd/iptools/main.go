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
	if args.CLIArgs.IP4Subnet != nil {
		if args.CLIArgs.IP4Subnet.SubnetRanges != nil {
			handler.IP4SubnetRanges(
				args.CLIArgs.IP4Subnet.SubnetRanges.IP,
				args.CLIArgs.IP4Subnet.SubnetRanges.Bits,
				args.CLIArgs.IP4Subnet.SubnetRanges.SecondaryBits,
			)
		}
		if args.CLIArgs.IP4Subnet.SubnetDivide != nil {
			handler.IP4SubnetDivide(
				args.CLIArgs.IP4Subnet.SubnetDivide.IP,
				args.CLIArgs.IP4Subnet.SubnetDivide.Bits,
				args.CLIArgs.IP4Subnet.SubnetDivide.SecondaryBits,
			)
		}
		if args.CLIArgs.IP4Subnet.SubnetDescribe != nil {
			handler.IP4SubnetDescribe(
				args.CLIArgs.IP4Subnet.SubnetDescribe.IP,
				args.CLIArgs.IP4Subnet.SubnetDescribe.Bits,
				args.CLIArgs.IP4Subnet.SubnetDescribe.SecondaryBits,
			)
		}
	}
	if args.CLIArgs.IP6Subnet != nil {
		if args.CLIArgs.IP6Subnet.IP6SubnetDescribe != nil {
			handler.IP6SubnetDescribe(
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.IP,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Bits,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Random,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Type,
			)
		}
		if args.CLIArgs.IP6Subnet.IP6RandomIPs != nil {
			handler.IP6RandomIPs(
				args.CLIArgs.IP6Subnet.IP6RandomIPs.Type,
				args.CLIArgs.IP6Subnet.IP6RandomIPs.Number,
			)
		}
	}
}
