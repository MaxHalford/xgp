package metric

import (
	"math"
)

// A Metric evaluates the performance of a predictive model.
type Metric interface {
	Apply(yTrue []float64, yPred []float64) (float64, error)
}

// Accuracy measures the fraction of matches between true classes and predicted
// classes.
type Accuracy struct{}

// Apply Accuracy.
func (eval Accuracy) Apply(yTrue []float64, yPred []float64) (float64, error) {

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

// InverseAccuracy measures the inverse accuracy.
type InverseAccuracy struct{}

// Apply InverseAccuracy.
func (eval InverseAccuracy) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var accuracy, err = Accuracy{}.Apply(yTrue, yPred)
	return -accuracy, err
}

// Precision measures the fraction of times a class was correctly predicted.
type Precision struct {
	Class float64
}

// Apply Precision.
func (eval Precision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	TP, err := cm.TruePositives(eval.Class)
	// Check class exists
	if err != nil {
		return 0, err
	}
	FP, err := cm.FalsePositives(eval.Class)
	// If the class has never been predicted return 0
	if TP+FP == 0 {
		return 0, nil
	}
	return TP / (TP + FP), nil
}

// InversePrecision measures the inverse accuracy.
type InversePrecision struct {
	Class float64
}

// Apply InversePrecision.
func (eval InversePrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var precision, err = Precision{Class: eval.Class}.Apply(yTrue, yPred)
	return -precision, err
}

// Recall measures the fraction of times a true class was predicted.
type Recall struct {
	Class float64
}

// Apply Recall.
func (eval Recall) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	TP, err := cm.TruePositives(eval.Class)
	// Check class exists
	if err != nil {
		return 0, err
	}
	FN, err := cm.FalseNegatives(eval.Class)
	// If the class has never been predicted return 0
	if TP+FN == 0 {
		return 0, nil
	}
	return TP / (TP + FN), nil
}

// InverseRecall measures the inverse accuracy.
type InverseRecall struct {
	Class float64
}

// Apply InverseRecall.
func (eval InverseRecall) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var recall, err = Recall{Class: eval.Class}.Apply(yTrue, yPred)
	return -recall, err
}

// F1Score measures the F1-score.
type F1Score struct {
	Class float64
}

// Apply F1Score.
func (eval F1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	TP, err := cm.TruePositives(eval.Class)
	// Check class exists
	if err != nil {
		return 0, err
	}
	FP, err := cm.FalsePositives(eval.Class)
	FN, err := cm.FalseNegatives(eval.Class)
	// If the class has never been predicted return 0
	if TP+FP == 0 || TP+FN == 0 {
		return 0, nil
	}
	var (
		precision = TP / (TP + FP)
		recall    = TP / (TP + FN)
		f1Score   = 2 * (precision * recall) / (precision + recall)
	)
	return f1Score, nil
}

// InverseRecall measures the inverse accuracy.
type InverseF1Score struct {
	Class float64
}

// Apply InverseF1Score.
func (eval InverseF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var f1Score, err = F1Score{Class: eval.Class}.Apply(yTrue, yPred)
	return -f1Score, err
}

// MinkowskiDistance measures the Minkowski distance.
type MinkowskiDistance struct {
	P float64
}

// Apply MinkowskiDistance.
func (eval MinkowskiDistance) Apply(yTrue []float64, yPred []float64) (float64, error) {

	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}

	var dist float64
	for i, y := range yTrue {
		dist += math.Pow(math.Abs(y-yPred[i]), eval.P)
	}
	return math.Pow(dist, 1/eval.P), nil
}
