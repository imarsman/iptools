package subnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"gopkg.in/yaml.v2"
	"inet.af/netaddr"
)

const (
	octets   = 4
	octetMax = 255
)

// Subnet an IP subnet
// https://www.calculator.net/ip-Subnet-calculator.html?cclass=any&csubnet=16&cip=10.0.0.0&ctype=ipv4&printit=0&x=43&y=21
// https://www.calculator.net/ip-Subnet-calculator.html
type Subnet struct {
	Prefix netaddr.IPPrefix `json:"prefix" yaml:"prefix"`
	IP     netaddr.IP       `json:"ipaddr" yaml:"ipaddr"`
	// TotalHosts    uint              `json:"totalhosts" yaml:"totalhosts"`
	SubnetHosts       uint64            `json:"subnethosts" yaml:"subnethosts"`
	DivisionIncrement uint64            `json:"divisionhosts" yaml:"divisionhosts"`
	Divisions         []netaddr.IPRange `json:"divisions" yaml:"divisions"`
	TotalDivisions    uint
}

// NewSubnet create new subnet from prefix
func NewSubnet(prefix string) (subnet *Subnet, err error) {
	subnet = new(Subnet)

	pfx, err := netaddr.ParseIPPrefix(prefix)
	if err != nil {
		return
	}

	if !pfx.IsValid() {
		return nil, errors.New("invalid prefix")
	}

	if pfx.IP().Is6() {
		return nil, errors.New("subnet too large for current implementation")
	}

	subnet.Prefix = pfx
	subnet.IP = pfx.IP()

	// Hosts is 2^[non-mask bits]
	subnet.SubnetHosts = uint64(math.Pow(2,
		float64(
			pfx.IP().BitLen()-subnet.Prefix.Bits(),
		)))

	// subnet.TotalHosts = uint(subnet.DivisionHosts * uint64(subnet.EqualRanges()))
	// subnet.DivisionIncrement = uint64(subnet.SubnetHosts) / uint64(math.Pow(2, float64(subnet.ClassHostBits())))
	subnet.DivisionIncrement = uint64(uint64(math.Pow(2, float64(subnet.ClassHostBits()))))
	// subnet.TotalDivisions = uint(subnet.SubnetHosts) / uint(subnet.DivisionIncrement)
	subnet.TotalDivisions = (uint(256)-uint(subnet.ClassByte()))/uint(subnet.DivisionIncrement) + 1
	subnet.Divisions = subnet.getDivisions()

	return subnet, nil
}

// UsableRange omit first and last IPs
func UsableRange(r netaddr.IPRange) netaddr.IPRange {
	r2 := netaddr.IPRangeFrom(r.From().Next(), r.To().Prior())

	return r2
}

// JSON get JSON for subnet
func (s *Subnet) JSON() (bytes []byte, err error) {
	bytes, err = json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	return bytes, nil
}

// YAML get YAML for subnet
func (s *Subnet) YAML() (bytes []byte, err error) {
	bytes, err = yaml.Marshal(s)
	if err != nil {
		return
	}

	return bytes, nil
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

// ShiftLeft performs a left bit shift operation on the provided bytes.
// If the bits count is negative, a right bit shift is performed.
func ShiftLeft(data []byte, bits int) []byte {
	n := len(data)
	if bits < 0 {
		bits = -bits
		for i := n - 1; i > 0; i-- {
			data[i] = data[i]>>bits | data[i-1]<<(8-bits)
		}
		data[0] >>= bits
	} else {
		for i := 0; i < n-1; i++ {
			data[i] = data[i]<<bits | data[i+1]>>(8-bits)
		}
		data[n-1] <<= bits
	}

	return data
}

// SubnetDivisions the set of equally sized subnet divisions for subnet
func (s *Subnet) getDivisions() (r []netaddr.IPRange) {

	r = make([]netaddr.IPRange, 0, 0)

	// Whatever the IP used to create the subnet range, use the subnet's first IP
	subnetBaseIP := s.Prefix.Range().From().As4()
	// subnetRanges := s.EqualRanges()
	// fmt.Println("shift", s.shift(subnetBaseIP[s.ClassByte()], s.ClassHostBits()))

	getIPRange := func(number int) (r netaddr.IPRange, last bool) {
		classByte := s.ClassByte()
		allBytes := subnetBaseIP

		byteToUse := subnetBaseIP[classByte]
		currentByte := byteToUse

		// In decimal terms, increment each by 2^host bits
		decimalIncrement := math.Pow(2, float64(s.ClassHostBits()))
		fmt.Println("decimal increment", decimalIncrement, s.ClassByte())

		ipFirst, ipLast := allBytes, allBytes

		// Will be for first
		ipFirstNewByte := uint(currentByte) + ((uint(number)) * uint(decimalIncrement))
		ipLastNewByte := uint(currentByte) + (uint(number+1) * uint(decimalIncrement))
		// fmt.Println("last new byte", ipLastNewByte)

		if ipLastNewByte > octetMax {
			ipLast[classByte] = octetMax
			fmt.Println(ipLast[classByte])
			last = true
		} else {
			ipLast[classByte] = byte(ipLastNewByte)
		}

		ipFirst[classByte] = byte(ipFirstNewByte)

		if byteToUse < octets {
			for i := classByte + 1; i < octets; i++ {
				ipFirst[i] = 0
			}
			for i := classByte + 1; i < octets; i++ {
				ipLast[i] = octetMax
			}
		} else {
			ipLast[3] = octetMax
		}

		r = netaddr.IPRangeFrom(netaddr.IPFrom4(ipFirst), netaddr.IPFrom4(ipLast))

		return
	}

	var newRange netaddr.IPRange
	var last bool
	if s.TotalDivisions == 1 {
		newRange, last = getIPRange(0)
		r = append(r)
	} else {

		// produce all subnet IP ranges
		for subNetNo := 0; subNetNo < int(s.TotalDivisions-1); subNetNo++ {
			newRange, last = getIPRange(subNetNo)
			r = append(r, newRange)
			if last {
				break
			}
		}
	}
	s.Divisions = r

	return
}
