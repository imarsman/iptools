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

// SubnetDescribe for calls to describe a subnet
type SubnetDescribe struct {
	IP   string `arg:"-i,--ip" help:""`
	Bits int    `arg:"-b,--bits" help:""`
}

// SubnetRanges for calls to get list of subnet ranges
type SubnetRanges struct {
	IP            string `arg:"-i,--ip" help:""`
	Bits          int    `arg:"-b,--bits" help:""`
	SecondaryBits int    `arg:"-s,--secondary-bits" help:""`
	Pretty        bool   `arg:"-p,--pretty" help:""`
}

// SubnetDivide for calls to divide subnet into networks
type SubnetDivide struct {
	IP            string `arg:"-i,--ip" help:""`
	Bits          int    `arg:"-b,--bits" help:""`
	SecondaryBits int    `arg:"-s,--secondary-bits" help:""`
	Pretty        bool   `arg:"-p,--pretty" help:""`
}

// Subnet top level subnet arg
type Subnet struct {
	SubnetRanges   *SubnetRanges   `arg:"subcommand:ranges" help:"divide a subnet into ranges"`
	SubnetDivide   *SubnetDivide   `arg:"subcommand:divide" help:"divide a subnet into smaller subnets"`
	SubnetDescribe *SubnetDescribe `arg:"subcommand:describe" help:"describe a subnet"`
}

// Args container for cli pargs
type Args struct {
	Subnet *Subnet `arg:"subcommand:subnet" help:"Get networks for subnet"`
}

// CLIArgs the args structure to be filled at runtime
var CLIArgs Args
