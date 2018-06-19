package metrics

import "math"

// MeanSquaredError measures the mean squared error (MSE).
type MeanSquaredError struct{}

// Apply MeanSquaredError.
func (metric MeanSquaredError) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
			sum += math.Pow(yTrue[i]-yPred[i], 2) * weights[i]
			ws += weights[i]
		}
		return sum / ws, nil
	}
	for i := range yTrue {
		sum += math.Pow(yTrue[i]-yPred[i], 2)
	}
	return sum / float64(len(yTrue)), nil
}

// Classification method of MeanSquaredError.
func (metric MeanSquaredError) Classification() bool {
	return false
}

// BiggerIsBetter method of MeanSquaredError.
func (metric MeanSquaredError) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of MeanSquaredError.
func (metric MeanSquaredError) NeedsProbabilities() bool {
	return false
}

// String method of MeanSquaredError.
func (metric MeanSquaredError) String() string {
	return "mse"
}
