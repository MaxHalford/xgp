package metrics

import (
	"math"
)

// BinaryLogLoss implementes logistic loss.
type BinaryLogLoss struct{}

// Apply BinaryLogLoss.
func (metric BinaryLogLoss) Apply(yTrue, yPred, weights []float64) (float64, error) {

	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var (
		score float64
		ws    float64
	)
	for i, yt := range yTrue {
		var (
			yp   = clip(yPred[i], 0.00001, 0.99999)
			loss = yt*math.Log(yp) + (1-yt)*math.Log(1-yp)
		)
		if weights != nil {
			score += weights[i] * loss
			ws += weights[i]
		} else {
			score += loss
		}
	}
	score *= -1

	if weights != nil {
		return score / ws, nil
	}
	return score / float64(len(yTrue)), nil
}

// Classification method of BinaryLogLoss.
func (metric BinaryLogLoss) Classification() bool {
	return true
}

// BiggerIsBetter method of BinaryLogLoss.
func (metric BinaryLogLoss) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of BinaryLogLoss.
func (metric BinaryLogLoss) NeedsProbabilities() bool {
	return true
}

// String method of BinaryLogLoss.
func (metric BinaryLogLoss) String() string {
	return "logloss"
}
