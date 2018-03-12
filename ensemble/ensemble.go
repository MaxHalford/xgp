package ensemble

// An Ensemble can train on a dataset given a "weak" learner.
type Ensemble interface {
	Fit(learner Learner, X [][]float64, Y []float64, W []float64, verbose bool) error
	Predict(X [][]float64, predictProba bool) ([]float64, error)
}
