package args

import (
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// IPs list of subnet IPs to make life easier
var IPs = []string{"99.236.32.0", "10.32.0.0", "192.168.1.1"}

// Define command structure to enable completion
var cmd = &complete.Command{
	Sub: map[string]*complete.Command{
		"subnet": {
			Sub: map[string]*complete.Command{
				// Scheduler health for an environment
				"ranges": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(IPs),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
						"pretty":         predict.Nothing,
					},
				},
				"divide": {
					Flags: map[string]complete.Predictor{
						"ip":             predict.Set(IPs),
						"bits":           predict.Nothing,
						"secondary-bits": predict.Nothing,
						"pretty":         predict.Nothing,
					},
				},
				"describe": {
					Flags: map[string]complete.Predictor{
						"ip":   predict.Set(IPs),
						"bits": predict.Nothing,
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
