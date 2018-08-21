package metrics

import "fmt"

// A NegativeMetric returns the negative output of a given Metric.
type NegativeMetric struct {
	Metric Metric
}

// Apply NegativeMetric.
func (neg NegativeMetric) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var output, err = neg.Metric.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	return -output, nil
}

// Classification method of NegativeMetric.
func (neg NegativeMetric) Classification() bool {
	return neg.Metric.Classification()
}

// BiggerIsBetter method of NegativeMetric.
func (neg NegativeMetric) BiggerIsBetter() bool {
	return !neg.Metric.BiggerIsBetter()
}

// NeedsProbabilities method of NegativeMetric.
func (neg NegativeMetric) NeedsProbabilities() bool {
	return neg.Metric.NeedsProbabilities()
}

// String method of NegativeMetric.
func (neg NegativeMetric) String() string {
	return fmt.Sprintf("neg_%s", neg.Metric.String())
}
