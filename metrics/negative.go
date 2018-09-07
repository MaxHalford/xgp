package metrics

import "fmt"

// A Negative returns the negative output of a given Metric.
type Negative struct {
	Metric Metric
}

// Apply Negative.
func (neg Negative) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var output, err = neg.Metric.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	return -output, nil
}

// Classification method of Negative.
func (neg Negative) Classification() bool {
	return neg.Metric.Classification()
}

// BiggerIsBetter method of Negative.
func (neg Negative) BiggerIsBetter() bool {
	return !neg.Metric.BiggerIsBetter()
}

// NeedsProbabilities method of Negative.
func (neg Negative) NeedsProbabilities() bool {
	return neg.Metric.NeedsProbabilities()
}

// String method of Negative.
func (neg Negative) String() string {
	return fmt.Sprintf("neg_%s", neg.Metric.String())
}
