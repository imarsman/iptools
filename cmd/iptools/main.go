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
		if args.CLIArgs.IP6Subnet.IP6SubnetGlobalUnicastDescribe != nil {
			handler.IP6DescribeGlobalUnicast(
				args.CLIArgs.IP6Subnet.IP6SubnetGlobalUnicastDescribe.IP,
				args.CLIArgs.IP6Subnet.IP6SubnetGlobalUnicastDescribe.Bits,
				args.CLIArgs.IP6Subnet.IP6SubnetGlobalUnicastDescribe.Random,
			)
		}
		if args.CLIArgs.IP6Subnet.IP6SubnetLinkLocalDescribe != nil {
			handler.IP6DescribeLinkLocal(
				args.CLIArgs.IP6Subnet.IP6SubnetLinkLocalDescribe.IP,
				args.CLIArgs.IP6Subnet.IP6SubnetLinkLocalDescribe.Bits,
				args.CLIArgs.IP6Subnet.IP6SubnetLinkLocalDescribe.Random,
			)
		}
	}
}
