package boosting

// A Booster can train on a dataset given a "weak" learner.
type Booster interface {
	Fit(learner Learner, X [][]float64, Y []float64, rounds int) error
	Predict(X [][]float64) ([]float64, error)
}
