package args

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
