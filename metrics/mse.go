package metrics

import "math"

// MSE measures the mean squared error (MSE).
type MSE struct{}

// Apply MSE.
func (mse MSE) Apply(yTrue, yPred, weights []float64) (float64, error) {
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

// Classification method of MSE.
func (mse MSE) Classification() bool {
	return false
}

// BiggerIsBetter method of MSE.
func (mse MSE) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of MSE.
func (mse MSE) NeedsProbabilities() bool {
	return false
}

// String method of MSE.
func (mse MSE) String() string {
	return "mse"
}

// Gradient computes yPred - yTrue.
func (mse MSE) Gradient(yTrue, yPred []float64) ([]float64, error) {
	var grad = make([]float64, len(yTrue))
	for i, y := range yTrue {
		grad[i] = yPred[i] - y
	}
	return grad, nil
}
