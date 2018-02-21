package metrics

// Accuracy measures the fraction of matches between true classes and predicted
// classes.
type Accuracy struct{}

// Apply Accuracy.
func (metric Accuracy) Apply(yTrue, yPred, weights []float64) (float64, error) {

	if len(yTrue) != len(yPred) {
		return 0, &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return 0, &errMismatchedLengths{len(yTrue), len(weights)}
	}

	var accuracy float64
	if weights != nil {
		var ws float64
		for i, y := range yTrue {
			if y == yPred[i] {
				accuracy += weights[i]
			}
			ws += weights[i]
		}
		return accuracy / ws, nil
	}
	for i, y := range yTrue {
		if y == yPred[i] {
			accuracy++
		}
	}
	return accuracy / float64(len(yTrue)), nil
}

// Classification method of Accuracy.
func (metric Accuracy) Classification() bool {
	return true
}

// BiggerIsBetter method of Accuracy.
func (metric Accuracy) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of Accuracy.
func (metric Accuracy) NeedsProbabilities() bool {
	return false
}

// String method of Accuracy.
func (metric Accuracy) String() string {
	return "accuracy"
}
