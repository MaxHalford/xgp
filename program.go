package xgp

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/MaxHalford/xgp/metrics"
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

// Classification determines if the Program has to perform classification or
// not.
func (prog Program) classification() bool {
	if prog.Estimator != nil {
		if prog.Estimator.LossMetric != nil {
			return prog.Estimator.LossMetric.Classification()
		}
	}
	return false
}

// Sigmoid applies the sigmoid transform.
func sigmoid(y float64) float64 {
	return 1 / (1 + math.Exp(-y))
}

// Binary converts a float64 to 0 or 1.
func binary(y float64) float64 {
	if y > 0.5 {
		return 1
	}
	return 0
}

// Predict predicts the output of a slice of features.
func (prog Program) Predict(X [][]float64, proba bool) ([]float64, error) {
	// Make predictions
	yPred := prog.Tree.Eval(X)
	// Check the predictions don't contain any NaNs
	if floats.HasNaN(yPred) {
		return nil, errors.New("yPred contains NaNs")
	}
	// Classification
	if prog.classification() {
		// Binary classification with probabilities
		if proba {
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
func (prog Program) PredictPartial(x []float64, proba bool) (float64, error) {
	// Make predictions
	yPred := prog.Tree.EvalRow(x)
	// Check the predictions don't contain any NaNs
	if math.IsNaN(yPred) {
		return 0, errors.New("yPred is NaN")
	}
	// Classification
	if prog.classification() {
		// Binary classification with probabilities
		if proba {
			return sigmoid(yPred), nil
		}
		// Binary classification
		return binary(yPred), nil
	}
	// Regression
	return yPred, nil
}

// MarshalJSON serializes a Program.
func (prog *Program) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Tree       tree.Tree `json:"tree"`
		LossMetric string    `json:"loss_metric"`
	}{
		Tree:       prog.Tree,
		LossMetric: prog.Estimator.LossMetric.String(),
	})
}

// UnmarshalJSON parses a Program.
func (prog *Program) UnmarshalJSON(bytes []byte) error {
	var serial = &struct {
		Tree       tree.Tree `json:"tree"`
		LossMetric string    `json:"loss_metric"`
	}{}
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	loss, err := metrics.ParseMetric(serial.LossMetric, 1)
	if err != nil {
		return err
	}
	prog.Tree = serial.Tree
	prog.Estimator = &Estimator{LossMetric: loss}
	return nil
}
