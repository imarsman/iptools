package subnet

import (
	"errors"
	"math"

	"inet.af/netaddr"
)

// Subnet an IP subnet
// https://www.calculator.net/ip-Subnet-calculator.html?cclass=any&csubnet=16&cip=10.0.0.0&ctype=ipv4&printit=0&x=43&y=21
// https://www.calculator.net/ip-Subnet-calculator.html
type Subnet struct {
	Prefix          netaddr.IPPrefix
	Hosts           uint64
	IPAddr          netaddr.IP
	SubnetSize      uint64
	subnetDivisions []netaddr.IPRange
}

// NewSubnet create new subnet from prefix
func NewSubnet(pfx netaddr.IPPrefix) (*Subnet, error) {
	subnet := new(Subnet)

	if !pfx.IsValid() {
		return nil, errors.New("invalid prefix")
	}

	if pfx.IP().Is6() {
		return nil, errors.New("subnet too large for current implementation")
	}

	subnet.Prefix = pfx

	// Hosts is 2^[non-mask bits]
	subnet.Hosts = uint64(math.Pow(2,
		float64(
			pfx.IP().BitLen()-subnet.Prefix.Bits(),
		)))

	subnet.SubnetSize = uint64(
		math.Pow(2, float64(subnet.ClassHostBits())),
	)

	return subnet, nil
}

// UsableRange omit first and last IPs
func UsableRange(r netaddr.IPRange) netaddr.IPRange {
	r2 := netaddr.IPRangeFrom(r.From().Next(), r.To().Prior())

	return r2
}

// ClassByte byte used for subnet
func (s *Subnet) ClassByte() uint8 {
	bits := s.Prefix.Bits()

	if bits <= 8 {
		return 0
	} else if bits <= 16 {
		return 1
	} else if bits <= 24 {
		return 2
	} else {
		return 3
	}
}

// ClassPartialBits bits used in mask block
func (s *Subnet) ClassPartialBits() uint8 {
	r := s.Prefix.Bits() % 8

	return r
}

// ClassHostBits bits remaining in mask block
func (s *Subnet) ClassHostBits() uint8 {
	r := s.ClassPartialBits()
	if r == 0 {
		return 8
	}

	return 8 - s.ClassPartialBits()
}

// EqualSubnets how many equal sized subnets can prefix be split into?
func (s *Subnet) EqualSubnets() uint {
	// return s.EqualSubnets
	// fmt.Println("partial remainder bits", s.PartialRemainderBits())
	// fmt.Println(math.Pow(2, float64(s.PartialRemainderBits())))
	return uint(math.Pow(2, float64(s.ClassPartialBits())))
}

// SubnetDivisions the set of equally sized subnet divisions for subnet
func (s *Subnet) SubnetDivisions() (r []netaddr.IPRange) {
	r = make([]netaddr.IPRange, 0, s.EqualSubnets()-1)

	// Return generated ranges if they have already been calculated
	if len(s.subnetDivisions) > 0 {
		return s.subnetDivisions
	}

	subnetCount := s.EqualSubnets()

	getIPRange := func(number int) netaddr.IPRange {
		classByte := s.ClassByte()
		allBytes := s.Prefix.IP().As4()

		byteToUse := s.Prefix.IP().As4()[classByte]
		currentByte := byteToUse

		ipFirstNewByte := uint(currentByte) + ((uint(number)) * uint(s.SubnetSize))
		ipLastNewByte := uint(currentByte) + (uint(number+1) * uint(s.SubnetSize))

		ipFirst := allBytes
		ipFirst[classByte] = byte(ipFirstNewByte)

		ipLast := allBytes
		ipLast[classByte] = byte(ipLastNewByte)

		if byteToUse < 4 {
			for i := classByte + 1; i < 4; i++ {
				ipFirst[i] = 0
			}
			for i := classByte + 1; i < 4; i++ {
				ipLast[i] = 255
			}
		} else {
			ipLast[3] = 255
		}

		return netaddr.IPRangeFrom(netaddr.IPFrom4(ipFirst), netaddr.IPFrom4(ipLast))
	}

	if subnetCount == 1 {
		r = append(r, getIPRange(0))
	} else {

		// produce all subnet IP ranges
		for subNetNo := 0; subNetNo < int(subnetCount-1); subNetNo++ {
			r = append(r, getIPRange(subNetNo))
		}
	}
	s.subnetDivisions = r

	return
}
