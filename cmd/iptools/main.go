package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/alexflint/go-arg"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/cmd/handler"
)

// the new comparable package is not there in 1.18
type comparable interface {
	~string | ~int | ~int64 | float64
}

// dedup de-duplicate a generic slice
func dedup[T comparable](slice []T) (result []T) {
	m := make(map[T]bool)
	for i := 0; i < len(slice); i++ {
		if !m[slice[i]] {
			result = append(result, slice[i])
			m[slice[i]] = true
		}
	}

	return
}

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
			if args.CLIArgs.IP6Subnet.IP6SubnetDescribe.IP == "" && !args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Random {
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Random = true
			}
			handler.IP6SubnetDescribe(
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.IP,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Bits,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Random,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.Type,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.JSON,
				args.CLIArgs.IP6Subnet.IP6SubnetDescribe.YAML,
			)
		}
		if args.CLIArgs.IP6Subnet.IP6RandomIPs != nil {
			handler.IP6RandomIPs(
				args.CLIArgs.IP6Subnet.IP6RandomIPs.Type,
				args.CLIArgs.IP6Subnet.IP6RandomIPs.Number,
			)
		}
	}
	if args.CLIArgs.Utilities != nil {
		if len(args.CLIArgs.Utilities.Lookup.Domains) != 0 {
			domains := args.CLIArgs.Utilities.Lookup.Domains
			sort.Strings(domains)
			domains = dedup(domains)

			args.CLIArgs.Utilities.Lookup.Domains = domains
			handler.LookupDomain(
				args.CLIArgs.Utilities.Lookup.Domains, args.CLIArgs.Utilities.Lookup.MXLookup,
				args.CLIArgs.Utilities.Lookup.JSON, args.CLIArgs.Utilities.Lookup.YAML,
			)
		} else {
			fmt.Println("No valid utilities option selected")
			os.Exit(1)
		}
	}

}
