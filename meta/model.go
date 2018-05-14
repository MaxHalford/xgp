package meta

import "github.com/MaxHalford/xgp"

// A Model can train on a dataset given a "weak" learner and return an
// Ensemble.
type Model interface {
	Fit(
		estimator *xgp.Estimator,
		XTrain [][]float64,
		YTrain []float64,
		XVal [][]float64,
		YVal []float64,
		verbose bool,
	) (Ensemble, error)
}
