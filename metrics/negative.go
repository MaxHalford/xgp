package metrics

import "fmt"

// A NegativeMetric returns the negative output of a given Metric.
type NegativeMetric struct {
	Metric Metric
}

// Apply NegativeMetric.
func (metric NegativeMetric) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var output, err = metric.Metric.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	return -output, nil
}

// Classification method of NegativeMetric.
func (metric NegativeMetric) Classification() bool {
	return metric.Metric.Classification()
}

// BiggerIsBetter method of NegativeMetric.
func (metric NegativeMetric) BiggerIsBetter() bool {
	return !metric.Metric.BiggerIsBetter()
}

// NeedsProbabilities method of NegativeMetric.
func (metric NegativeMetric) NeedsProbabilities() bool {
	return metric.Metric.NeedsProbabilities()
}

// String method of NegativeMetric.
func (metric NegativeMetric) String() string {
	return fmt.Sprintf("neg_%s", metric.Metric.String())
}
