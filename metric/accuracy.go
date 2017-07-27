package metric

// Accuracy measures the fraction of matches between true classes and predicted
// classes.
type Accuracy struct{}

// Apply Accuracy.
func (metric Accuracy) Apply(yTrue []float64, yPred []float64) (float64, error) {

	if len(yTrue) != len(yPred) {
		return 0, &errMismatchedLengths{len(yTrue), len(yPred)}
	}

	var accuracy float64
	for i, y := range yTrue {
		if y == yPred[i] {
			accuracy++
		}
	}

	return accuracy / float64(len(yTrue)), nil
}

// NegativeAccuracy measures the inverse accuracy.
type NegativeAccuracy struct{}

// Apply NegativeAccuracy.
func (metric NegativeAccuracy) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var accuracy, err = Accuracy{}.Apply(yTrue, yPred)
	return -accuracy, err
}
