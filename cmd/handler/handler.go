package handler

import (
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/imarsman/iptools/cmd/args"
	"github.com/imarsman/iptools/pkg/subnet"
	"inet.af/netaddr"
)

// SubnetDivide divide a subnet into ranges
func SubnetDivide(ip string, mask uint8, secondaryMask uint8) {
	var err error
	var s *subnet.IPV4Subnet
	prefix := fmt.Sprintf("%s/%d", ip, mask)
	s, err = subnet.NewFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ranges := []netaddr.IPRange{}
	var s2 *subnet.IPV4Subnet
	if secondaryMask != 0 {
		prefix := fmt.Sprintf("%s/%d", ip, secondaryMask)
		s2, err = subnet.NewFromPrefix(prefix)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ranges, err = s.NetworkRangesInSubnets(s2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		s2 = s
		ranges, err = s.NetworkRanges()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if args.CLIArgs.Subnet.SubnetDivide.Pretty {
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Category"},
				{Align: simpletable.AlignCenter, Text: "Value"},
			},
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.Prefix.String())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		r = []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Subnet IP")},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s.IP.String())},
		}
		if mask != secondaryMask {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary subnet")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", s2.Prefix.String())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}
		if mask == secondaryMask {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.NetworkHosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		} else {
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Secondary networks")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s2.Networks())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Network hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s.NetworkHosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
			r = []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", "Sub Network hosts")},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", s2.NetworkHosts())},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}
		fmt.Println()
		table.SetStyle(simpletable.StyleCompactLite)
		fmt.Println(table.String())
	}

	if args.CLIArgs.Subnet.SubnetDivide.Pretty {
		fmt.Println()
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "Start"},
				{Align: simpletable.AlignCenter, Text: "End"},
			},
		}
		for _, r := range ranges {
			cell := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", r.From().String())},
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", r.To().String())},
			}
			table.Body.Cells = append(table.Body.Cells, cell)
		}
		table.SetStyle(simpletable.StyleCompactLite)
		fmt.Println(table.String())
	} else {
		for _, r := range ranges {
			fmt.Println(r.String())
		}
	}
}
