package ensemble

// A Learner can be used by an Ensemble to train on a dataset.
type Learner interface {
	Fit(X [][]float64, Y []float64, W []float64, verbose bool) (Predictor, error)
}
