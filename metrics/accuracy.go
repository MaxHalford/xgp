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

	var (
		accuracy float64
		ws       float64
	)
	for i, y := range yTrue {
		if y == yPred[i] {
			if weights != nil {
				accuracy += weights[i]
			} else {
				accuracy++
			}
		}
		if weights != nil {
			ws += weights[i]
		}
	}

	if weights != nil {
		return accuracy / ws, nil
	}
	return accuracy / float64(len(yTrue)), nil
}

// NegativeAccuracy measures the inverse accuracy.
type NegativeAccuracy struct{}

// Apply NegativeAccuracy.
func (metric NegativeAccuracy) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var accuracy, err = Accuracy{}.Apply(yTrue, yPred, weights)
	return -accuracy, err
}
