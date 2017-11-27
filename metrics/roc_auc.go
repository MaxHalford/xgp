package metrics

import (
	"errors"

	"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/gonum/stat"
)

// ROCAUC measures the ROC AUC score.
type ROCAUC struct{}

// Apply ROCAUC.
func (metric ROCAUC) Apply(yTrue, yPred, weights []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return 0, &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return 0, &errMismatchedLengths{len(yTrue), len(weights)}
	}

	// Convert yTrue to a bool slice
	var classes = make([]bool, len(yTrue))
	for i, y := range yTrue {
		switch y {
		case 0:
			classes[i] = false
		case 1:
			classes[i] = true
		default:
			return 0, errors.New("ROC AUC is restricted to binary classification")
		}
	}

	// Copy yPred and weights before in-place sorting
	var yPredCopy = make([]float64, len(yPred))
	copy(yPredCopy, yPred)
	var weightsCopy []float64
	if weights != nil {
		var weightsCopy = make([]float64, len(weights))
		copy(weightsCopy, weights)
	}

	// Sort the slices according to yPred
	stat.SortWeightedLabeled(yPredCopy, classes, weightsCopy)

	var (
		tpr, fpr = stat.ROC(0, yPredCopy, classes, weightsCopy)
		auc      = integrate.Trapezoidal(tpr, fpr)
	)

	return auc, nil
}

// Classification method of ROCAUC.
func (metric ROCAUC) Classification() bool {
	return true
}

// BiggerIsBetter method of ROCAUC.
func (metric ROCAUC) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of ROCAUC.
func (metric ROCAUC) NeedsProbabilities() bool {
	return true
}

// String method of ROCAUC.
func (metric ROCAUC) String() string {
	return "roc_auc"
}
