package metrics

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

// AbsolutePearson measures the ROC AUC score.
type AbsolutePearson struct{}

// Apply AbsolutePearson.
func (ap AbsolutePearson) Apply(yTrue, yPred, weights []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return 0, &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return 0, &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var rho = stat.Correlation(yTrue, yPred, weights)

	return math.Abs(rho), nil
}

// Classification method of AbsolutePearson.
func (ap AbsolutePearson) Classification() bool {
	return false
}

// BiggerIsBetter method of AbsolutePearson.
func (ap AbsolutePearson) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of AbsolutePearson.
func (ap AbsolutePearson) NeedsProbabilities() bool {
	return false
}

// String method of AbsolutePearson.
func (ap AbsolutePearson) String() string {
	return "pearson"
}
