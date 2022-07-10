package args

type SubnetDescribe struct {
	IP   string `arg:"-i,--ip" help:""`
	Mask int    `arg:"-m,--mask" help:""`
}

type SubnetRanges struct {
	IP      string `arg:"-i,--ip" help:""`
	Mask    int    `arg:"-m,--mask" help:""`
	SubMask int    `arg:"-s,--sub-mask" help:""`
	Pretty  bool   `arg:"-p,--pretty" help:""`
}

type SubnetDivide struct {
	IP      string `arg:"-i,--ip" help:""`
	Mask    int    `arg:"-m,--mask" help:""`
	SubMask int    `arg:"-s,--sub-mask" help:""`
	Pretty  bool   `arg:"-p,--pretty" help:""`
}

type Subnet struct {
	SubnetRanges   *SubnetRanges   `arg:"subcommand:ranges" help:"divide a subnet into ranges"`
	SubnetDivide   *SubnetDivide   `arg:"subcommand:divide" help:"divide a subnet into smaller subnets"`
	SubnetDescribe *SubnetDescribe `arg:"subcommand:describe" help:"describe a subnet"`
}

type Args struct {
	Subnet *Subnet `arg:"subcommand:subnet" help:"Get networks for subnet"`
}

var CLIArgs Args
