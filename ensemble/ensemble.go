package ensemble

import (
	"math/rand"

	"github.com/MaxHalford/koza/metrics"
)

// An Ensemble can train on a dataset given a "weak" learner.
type Ensemble interface {
	Fit(
		learner Learner,
		XTrain [][]float64,
		YTrain []float64,
		WTrain []float64,
		XVal [][]float64,
		YVal []float64,
		WVal []float64,
		verbose bool,
		evalMetric metrics.Metric,
		rng *rand.Rand,
	) error
	Predict(X [][]float64, predictProba bool) ([]float64, error)
}
