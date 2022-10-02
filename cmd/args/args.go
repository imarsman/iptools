package args

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// GitCommit the git commit hash at compile time
var GitCommit string

// GitLastTag the last tag
var GitLastTag string

// GitExactTag extract tag
var GitExactTag string

// Date the compile date
var Date string

// Version get version information
func (Args) Version() string {
	var buf = new(bytes.Buffer)

	msg := os.Args[0]
	buf.WriteString(fmt.Sprintln(msg))
	buf.WriteString(fmt.Sprintln(strings.Repeat("-", len(msg))))

	if GitCommit != "" {
		buf.WriteString(fmt.Sprintf("Commit: %8s\n", GitCommit))
	}
	if Date != "" {
		buf.WriteString(fmt.Sprintf("Date: %23s\n", Date))
	}
	if GitExactTag != "" {
		buf.WriteString(fmt.Sprintf("Tag: %11s\n", GitExactTag))
	}
	buf.WriteString(fmt.Sprintf("OS: %11s\n", runtime.GOOS))
	buf.WriteString(fmt.Sprintf("ARCH: %8s\n", runtime.GOARCH))

	return buf.String()
}

// IP4SubnetDescribe for calls to describe a subnet
type IP4SubnetDescribe struct {
	IP            string `arg:"-i,--ip" help:""`
	Bits          int    `arg:"-b,--bits" help:""`
	SecondaryBits int    `arg:"-s,--secondary-bits" help:""`
}

// IP4SubnetRanges for calls to get list of subnet ranges
type IP4SubnetRanges struct {
	IP            string `arg:"-i,--ip" help:""`
	Bits          int    `arg:"-b,--bits" help:""`
	SecondaryBits int    `arg:"-s,--secondary-bits" help:""`
	Pretty        bool   `arg:"-p,--pretty" help:""`
}

// IP4SubnetDivide for calls to divide subnet into networks
type IP4SubnetDivide struct {
	IP            string `arg:"-i,--ip" help:""`
	Bits          int    `arg:"-b,--bits" help:""`
	SecondaryBits int    `arg:"-s,--secondary-bits" help:""`
	Pretty        bool   `arg:"-p,--pretty" help:""`
}

// IP6SubnetGlobalUnicastDescribe for calls to describe a subnet
type IP6SubnetGlobalUnicastDescribe struct {
	IP     string `arg:"-i,--ip" help:"IP address"`
	Random bool   `arg:"-r,--random" help:"generate random IP"`
	Bits   int    `arg:"-b,--bits" help:"subnet bits"`
}

// IP6RandomIPs get random list of IPs of type
type IP6RandomIPs struct {
	Number int    `arg:"-n,--number" help:"generate random IP"`
	Type   string `arg:"-t,--type" help:"global-unicast, link-local, unique-local, multicast"`
}

// IP6SubnetDescribe for calls to describe a subnet
type IP6SubnetDescribe struct {
	IP     string `arg:"-i,--ip" help:"IP address"`
	Random bool   `arg:"-r,--random" help:"generate random IP"`
	Bits   int    `arg:"-b,--bits" help:"subnet bits"`
	Type   string `arg:"-t,--type" help:"global-unicast, link-local, unique-local, multicast, interface-local-multicast, link-local-multicast"`
}

// IP6Subnet IP6 calls
type IP6Subnet struct {
	// IP6SubnetGlobalUnicastDescribe *IP6SubnetGlobalUnicastDescribe `arg:"subcommand:global-unicast-describe"`
	IP6SubnetDescribe *IP6SubnetDescribe `arg:"subcommand:describe"`
	IP6RandomIPs      *IP6RandomIPs      `arg:"subcommand:random-ips"`
}

// IP4Subnet top level IP4 subnet arg
type IP4Subnet struct {
	SubnetRanges   *IP4SubnetRanges   `arg:"subcommand:ranges" help:"divide a subnet into ranges"`
	SubnetDivide   *IP4SubnetDivide   `arg:"subcommand:divide" help:"divide a subnet into smaller subnets"`
	SubnetDescribe *IP4SubnetDescribe `arg:"subcommand:describe" help:"describe a subnet"`
}

// Args container for cli pargs
type Args struct {
	IP4Subnet *IP4Subnet `arg:"subcommand:subnetip4" help:"Get networks for subnet"`
	IP6Subnet *IP6Subnet `arg:"subcommand:subnetip6" help:"Get IP6 address information"`
}

// CLIArgs the args structure to be filled at runtime
var CLIArgs Args
