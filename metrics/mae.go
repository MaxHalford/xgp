package metrics

import "math"

// MeanAbsoluteError measures the mean absolute error (MAE).
type MeanAbsoluteError struct{}

// Apply MeanAbsoluteError.
func (metric MeanAbsoluteError) Apply(yTrue, yPred, weights []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var sum float64
	if weights != nil {
		var ws float64
		for i := range yTrue {
			sum += math.Abs(yTrue[i]-yPred[i]) * weights[i]
			ws += weights[i]
		}
		return sum / ws, nil
	}
	for i := range yTrue {
		sum += math.Abs(yTrue[i] - yPred[i])
	}
	return sum / float64(len(yTrue)), nil
}

// Classification method of MeanAbsoluteError.
func (metric MeanAbsoluteError) Classification() bool {
	return false
}

// BiggerIsBetter method of MeanAbsoluteError.
func (metric MeanAbsoluteError) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of MeanAbsoluteError.
func (metric MeanAbsoluteError) NeedsProbabilities() bool {
	return false
}

// String method of MeanAbsoluteError.
func (metric MeanAbsoluteError) String() string {
	return "mae"
}
