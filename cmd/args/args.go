package args

type SubnetSplit struct {
	IP     string
	Mask   int
	Pretty bool
}

type Subnet struct {
	SubnetSplit *SubnetSplit
}
