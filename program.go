package xgp

import (
	"errors"
	"math"

	"github.com/MaxHalford/xgp/tree"
	"github.com/gonum/floats"
)

// A Program is simply an abstraction of top of a Tree.
type Program struct {
	Tree      tree.Tree              `json:"tree"`
	Task      Task                   `json:"task"`
	DRS       *DynamicRangeSelection `json:"drs"`
	Estimator *Estimator             `json:"-"`
}

// String representation of a Program.
func (prog Program) String() string {
	return prog.Tree.String()
}

// Clone a Program.
func (prog Program) clone() Program {
	var clone = Program{
		Tree:      prog.Tree.Clone(),
		Task:      prog.Task,
		Estimator: prog.Estimator,
	}
	if prog.DRS != nil {
		clone.DRS = prog.DRS.clone()
	}
	return clone
}

// Predict predicts the output of a slice of features.
func (prog Program) Predict(X [][]float64, predictProba bool) ([]float64, error) {
	// Make predictions
	yPred := prog.Tree.Eval(X)
	// Check the predictions don't contain any NaNs
	if floats.HasNaN(yPred) {
		return nil, errors.New("yPred contains NaNs")
	}
	// Binary classification
	if prog.Task.binaryClassification() {
		if predictProba {
			for i, y := range yPred {
				yPred[i] = sigmoid(y)
			}
			return yPred, nil
		}
		for i, y := range yPred {
			yPred[i] = binary(y)
		}
		return yPred, nil
	}
	// Multi-class classification
	if prog.Task.multiClassification() {
		return prog.DRS.Predict(yPred), nil
	}
	// Regression
	return yPred, nil
}

// PredictPartial predicts the output of a slice of features.
func (prog Program) PredictPartial(x []float64, predictProba bool) (float64, error) {
	// Make predictions
	yPred := prog.Tree.EvalRow(x)
	// Check the predictions don't contain any NaNs
	if math.IsNaN(yPred) {
		return -1, errors.New("yPred is NaN")
	}
	// Binary classification
	if prog.Task.binaryClassification() {
		if predictProba {
			return sigmoid(yPred), nil
		}
		return binary(yPred), nil
	}
	// Multi-class classification
	if prog.Task.multiClassification() {
		return prog.DRS.PredictPartial(yPred), nil
	}
	// Regression
	return yPred, nil
}
