package args

import (
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// IP6ips list of subnet IP6ips to make life easier
var IP6ips = []string{"99.236.32.0", "10.32.0.0", "192.168.1.1"}

// IP6PrefixBits number of bits for prefix
var IP6PrefixBits = []string{"64", "48", "10", "8"}

// IP6Types IP6 address types
var IP6Types = []string{
	"global-unicast",
	"link-local",
	"unique-local",
}

// Define command structure to enable completion
var cmd = &complete.Command{
	Sub: map[string]*complete.Command{
		"subnetip4": {
			Sub: map[string]*complete.Command{
				// Scheduler health for an environment
				"ranges": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(IP6ips),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
						"pretty":         predict.Nothing,
					},
				},
				"divide": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(IP6ips),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
						"pretty":         predict.Nothing,
					},
				},
				"describe": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(IP6ips),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
					},
				},
			},
		},
		"subnetip6": {
			Sub: map[string]*complete.Command{
				// Describe an IP
				"describe": {
					Flags: map[string]complete.Predictor{
						"ip":     predict.Nothing,
						"bits":   predict.Set(IP6PrefixBits),
						"random": predict.Nothing,
						"type":   predict.Set(IP6Types),
					},
				},
				"random-ips": {
					Flags: map[string]complete.Predictor{
						"number": predict.Nothing,
						"type":   predict.Set(IP6Types),
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
