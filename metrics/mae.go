package metrics

import "math"

// MAE measures the mean absolute error (MAE).
type MAE struct{}

// Apply MAE.
func (mae MAE) Apply(yTrue, yPred, weights []float64) (float64, error) {
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

// Classification method of MAE.
func (mae MAE) Classification() bool {
	return false
}

// BiggerIsBetter method of MAE.
func (mae MAE) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of MAE.
func (mae MAE) NeedsProbabilities() bool {
	return false
}

// String method of MAE.
func (mae MAE) String() string {
	return "mae"
}
