package metrics

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

// AbsolutePearsonCorrelation measures the ROC AUC score.
type AbsolutePearsonCorrelation struct{}

// Apply AbsolutePearsonCorrelation.
func (metric AbsolutePearsonCorrelation) Apply(yTrue, yPred, weights []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return 0, &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return 0, &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var rho = stat.Correlation(yTrue, yPred, weights)

	return math.Abs(rho), nil
}

// Classification method of AbsolutePearsonCorrelation.
func (metric AbsolutePearsonCorrelation) Classification() bool {
	return false
}

// BiggerIsBetter method of AbsolutePearsonCorrelation.
func (metric AbsolutePearsonCorrelation) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of AbsolutePearsonCorrelation.
func (metric AbsolutePearsonCorrelation) NeedsProbabilities() bool {
	return false
}

// String method of AbsolutePearsonCorrelation.
func (metric AbsolutePearsonCorrelation) String() string {
	return "pearson"
}
