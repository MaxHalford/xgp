package xgp

import (
	"encoding/json"
	"errors"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
	"github.com/gonum/floats"
)

// A Program is a thin layer on top of an Operator.
type Program struct {
	*GP
	Op op.Operator
}

// String formatting.
func (prog Program) String() string {
	return prog.Op.String()
}

// Classification determines if the Program has to perform classification or
// not. It does so by looking at the GP's LossMetric.
func (prog Program) classification() bool {
	if prog.GP != nil {
		if prog.GP.LossMetric != nil {
			return prog.GP.LossMetric.Classification()
		}
	}
	return false
}

// Predict predicts the output of a slice of features.
func (prog Program) Predict(X [][]float64, proba bool) ([]float64, error) {
	// Make predictions
	yPred := prog.Op.Eval(X)
	// Check the predictions don't contain any NaNs
	if floats.HasNaN(yPred) {
		return nil, errors.New("yPred contains NaNs")
	}
	// Regression
	if !prog.classification() {
		return yPred, nil
	}
	// Classification
	var transform = map[bool]func(float64) float64{true: sigmoid, false: binary}[proba]
	for i, y := range yPred {
		yPred[i] = transform(y)
	}
	return yPred, nil
}

// PredictPartial is a convenience function on top of Predict to make
// predictions on a single instance.
func (prog Program) PredictPartial(x []float64, proba bool) (float64, error) {
	var X = make([][]float64, len(x))
	for i, xi := range x {
		X[i] = []float64{xi}
	}
	yPred, err := prog.Predict(X, proba)
	if err != nil {
		return 0, err
	}
	return yPred[0], nil
}

type serialProgram struct {
	Op         op.SerialOp `json:"op"`
	LossMetric string      `json:"loss_metric"`
}

// MarshalJSON serializes a Program.
func (prog Program) MarshalJSON() ([]byte, error) {
	return json.Marshal(&serialProgram{
		Op:         op.SerializeOp(prog.Op),
		LossMetric: prog.GP.LossMetric.String(),
	})
}

// UnmarshalJSON parses a Program.
func (prog *Program) UnmarshalJSON(bytes []byte) error {
	var serial = &serialProgram{}
	if err := json.Unmarshal(bytes, serial); err != nil {
		return err
	}
	loss, err := metrics.ParseMetric(serial.LossMetric, 1)
	if err != nil {
		return err
	}
	operator, err := op.ParseOp(serial.Op)
	if err != nil {
		return err
	}
	prog.Op = operator
	prog.GP = &GP{LossMetric: loss}
	return nil
}
