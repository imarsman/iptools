package args

type SubnetSplit struct {
	IP      string
	Mask    int
	SubMask int
	Pretty  bool
}

type Subnet struct {
	SubnetSplit *SubnetSplit
}
