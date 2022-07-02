package args
import (
"github.com/posener/complete/v2"
 	"github.com/posener/complete/v2/predict"
)


cmd := &complete.Command{
	 	Flags: map[string]complete.Predictor{
 			"name":      predict.Set{"foo", "bar", "foo bar"},
 			"something": predict.Something,
 			"nothing":   predict.Nothing,
 		},
 	}
 	// Run the completion - provide it with the binary name.
 	cmd.Complete("my-program")