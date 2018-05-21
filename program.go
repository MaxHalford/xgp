package xgp

import (
	"errors"
	"math"

	"github.com/MaxHalford/xgp/tree"
	"github.com/gonum/floats"
)

// A Program is simply an abstraction of top of a Tree.
type Program struct {
	Tree      tree.Tree  `json:"tree"`
	Estimator *Estimator `json:"-"`
}

// String representation of a Program.
func (prog Program) String() string {
	return prog.Tree.String()
}

// Clone a Program.
func (prog Program) clone() Program {
	return Program{
		Tree:      prog.Tree.Clone(),
		Estimator: prog.Estimator,
	}
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
	if prog.Estimator != nil && prog.Estimator.LossMetric.Classification() {
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
	if prog.Estimator.LossMetric.Classification() {
		if predictProba {
			return sigmoid(yPred), nil
		}
		return binary(yPred), nil
	}
	// Regression
	return yPred, nil
}
