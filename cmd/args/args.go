package args

type SubnetDivide struct {
	IP      string `arg:"-i,--ip" help:""`
	Mask    int    `arg:"-m,--mask" help:""`
	SubMask int    `arg:"-s,--sub-mask" help:""`
	Pretty  bool   `arg:"-p,--pretty" help:""`
}

type Subnet struct {
	SubnetDivide *SubnetDivide `arg:"subcommand:divide" help:""`
}

type Args struct {
	Subnet *Subnet `arg:"subcommand:subnet" help:"Get networks for subnet"`
}
