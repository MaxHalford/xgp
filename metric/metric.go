package metric

// A Metric metricuates the performance of a predictive model.
type Metric interface {
	Apply(yTrue []float64, yPred []float64) (float64, error)
}
