package args

import (
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// ip4ips list of subnet ip4ips to make life easier
var ip4ips = []string{"99.236.32.0", "10.32.0.0", "192.168.1.1"}

// ip6PrefixBits number of bits for prefix
var ip6PrefixBits = []string{"64", "48", "10", "8"}

var domains = []string{"cisco.com", "workday.cisco.com", "ibm.com", "java.com"}

// ip6Types IP6 address types
var ip6Types = []string{
	"global-unicast",
	"link-local",
	"private",
	"multicast",
	"interface-local-multicast",
	"link-local-multicast",
}

// Define command structure to enable completion
var cmd = &complete.Command{
	Sub: map[string]*complete.Command{
		"subnetip4": {
			Sub: map[string]*complete.Command{
				// Scheduler health for an environment
				"ranges": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(ip4ips),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
						"pretty":         predict.Nothing,
					},
				},
				"divide": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(ip4ips),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
						"pretty":         predict.Nothing,
					},
				},
				"describe": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(ip4ips),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
					},
				},
			},
		},
		"ip6": {
			Sub: map[string]*complete.Command{
				// Describe an IP
				"describe": {
					Flags: map[string]complete.Predictor{
						"ip":     predict.Nothing,
						"bits":   predict.Set(ip6PrefixBits),
						"random": predict.Nothing,
						"type":   predict.Set(ip6Types),
					},
				},
				"random-ips": {
					Flags: map[string]complete.Predictor{
						"number": predict.Nothing,
						"type":   predict.Set(ip6Types),
					},
				},
			},
		},
		"utilities": {
			Sub: map[string]*complete.Command{
				"lookup-domains": {
					Flags: map[string]complete.Predictor{
						"domain": predict.Set(domains),
					},
				},
			},
		},
	},
}

// InitializeCompletion initialize cmdline completion
func InitializeCompletion() {
	// Run the completion - provide it with the binary name.
	cmd.Complete("iptools")
}
