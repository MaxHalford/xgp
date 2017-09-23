package boosting

// A Learner can be used by a Booster to train on a dataset.
type Learner interface {
	Learn(X [][]float64, Y []float64) (Predictor, error)
}
