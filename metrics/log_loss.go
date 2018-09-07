package metrics

import (
	"math"
)

// LogLoss implementes logistic loss.
type LogLoss struct{}

// Apply LogLoss.
func (ll LogLoss) Apply(yTrue, yPred, weights []float64) (float64, error) {

	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var score float64
	if weights != nil {
		var ws float64
		for i, yt := range yTrue {
			var yp = clip(yPred[i], 0.00001, 0.99999)
			score += weights[i]*yt*math.Log(yp) + (1-yt)*math.Log(1-yp)
			ws += weights[i]
		}
		return -score / ws, nil
	}

	for i, yt := range yTrue {
		var yp = clip(yPred[i], 0.00001, 0.99999)
		score += yt*math.Log(yp) + (1-yt)*math.Log(1-yp)
	}
	return -score / float64(len(yTrue)), nil
}

// Classification method of LogLoss.
func (ll LogLoss) Classification() bool {
	return true
}

// BiggerIsBetter method of LogLoss.
func (ll LogLoss) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of LogLoss.
func (ll LogLoss) NeedsProbabilities() bool {
	return true
}

// String method of LogLoss.
func (ll LogLoss) String() string {
	return "logloss"
}

// Gradients computes yPred[i] - yTrue[i].
func (ll LogLoss) Gradients(yTrue, yPred []float64) ([]float64, error) {
	var grad = make([]float64, len(yTrue))
	for i, y := range yTrue {
		grad[i] = yPred[i] - y
	}
	return grad, nil
}
