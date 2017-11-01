package metrics

// A Metric metricuates the performance of a predictive model. yTrue, yPred, and
// weights should all have the same length. If weights is nil then uniform
// weights are used.
type Metric interface {
	Apply(yTrue, yPred, weights []float64) (float64, error)
	Classification() bool
	BiggerIsBetter() bool
	NeedsProbabilities() bool
	String() string
}
