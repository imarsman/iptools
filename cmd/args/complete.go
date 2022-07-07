package args

import (
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// IPs list of subnet IPs to make life easier
var IPs = []string{"99.236.32.0", "10.32.0.0"}

var cmd = &complete.Command{
	Sub: map[string]*complete.Command{
		"subnet": {
			Sub: map[string]*complete.Command{
				// Scheduler health for an environment
				"divide": {
					Flags: map[string]complete.Predictor{
						"ip":       predict.Set(IPs),
						"mask":     predict.Nothing,
						"sub-mask": predict.Nothing,
						"pretty":   predict.Nothing,
					},
				},
				"describe": {
					Flags: map[string]complete.Predictor{
						"ip":   predict.Set(IPs),
						"mask": predict.Nothing,
					},
				},
			},
		},
	},
}

func InitializeCompletion() {
	// Run the completion - provide it with the binary name.
	cmd.Complete("iptools")
}
