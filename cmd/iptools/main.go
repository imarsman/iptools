package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/imarsman/iptools/pkg/tools"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// GitCommit the git commit hash at compile time
var GitCommit string

// GitLastTag the last tag
var GitLastTag string

// GitExactTag extract tag
var GitExactTag string

// Date the compile date
var Date string

func printHelp(p *arg.Parser) {
	fmt.Println()
	var help bytes.Buffer
	p.WriteHelp(&help)
	fmt.Println(help.String())
}

func main() {
	var maskLengths = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32"}
	var args struct {
		Encode  int    `arg:"-e,--encode" help:"encode range to mask"`
		Decode  string `arg:"-d,--decode" help:"decode range from mask"`
		Mask    int    `arg:"-m,--mask" help:"get cidr mask"`
		Verbose bool   `arg:"-V,--verbose" help:"give more output"`
	}

	var completion = &complete.Command{
		Flags: map[string]complete.Predictor{
			"encode":  predict.Set(maskLengths),
			"decode":  predict.Nothing,
			"verbose": predict.Nothing,
			"mask":    predict.Set(maskLengths),
		},
	}
	// Tell completion to assess what has been defined
	completion.Complete("slvmmask")

	p, err := arg.NewParser(arg.Config{Program: "slvmmask"}, &args)
	if err != nil {

	}
	arg.MustParse(&args)

	if args.Encode != 0 {
		mask := net.CIDRMask(args.Encode, 32)

		encoded, err := json.Marshal([]byte(mask))
		if err != nil {
			fmt.Println(err)
			printHelp(p)

			os.Exit(1)
		}
		encoded = []byte(strings.ReplaceAll(string(encoded), `"`, ""))

		fmt.Printf("%s\n", encoded)
	} else if args.Decode != "" {
		masklen, err := tools.DecodeMaskBase64(string(args.Decode), args.Verbose)
		if err != nil {
			fmt.Println(err)
			printHelp(p)

			os.Exit(1)
		}

		fmt.Printf("%d\n", masklen)
	} else if args.Mask != 0 {
		cidrString, err := tools.CIDR(args.Mask)
		if err != nil {
			fmt.Println(err)
			printHelp(p)

			os.Exit(1)
		}
		fmt.Println(cidrString)
	} else {
		fmt.Println("-encode and -decode parameters missing")
		printHelp(p)

		os.Exit(1)
	}
}
