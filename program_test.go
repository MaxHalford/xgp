package xgp

import (
	"fmt"
	"math"
	"testing"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
)

func TestProgramPredict(t *testing.T) {
	var testCases = []struct {
		X           [][]float64
		program     Program
		proba       bool
		y           []float64
		raisesError bool
	}{
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program:     Program{nil, op.Add{op.Var{0}, op.Var{1}}},
			proba:       false,
			y:           []float64{-0.2, 0.1, 0.6, 3},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: nil,
			},
			proba:       true,
			y:           []float64{-0.2, 0.1, 0.6, 3},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{},
			},
			proba:       false,
			y:           []float64{-0.2, 0.1, 0.6, 3},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{},
			},
			proba:       true,
			y:           []float64{-0.2, 0.1, 0.6, 3},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{LossMetric: metrics.MSE{}},
			},
			proba:       false,
			y:           []float64{-0.2, 0.1, 0.6, 3},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{LossMetric: metrics.MSE{}},
			},
			proba:       true,
			y:           []float64{-0.2, 0.1, 0.6, 3},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{LossMetric: metrics.Accuracy{}},
			},
			proba:       false,
			y:           []float64{0, 1, 1, 1},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{LossMetric: metrics.Accuracy{}},
			},
			proba:       true,
			y:           []float64{0.45017, 0.52498, 0.64566, 0.95257},
			raisesError: false,
		},
		{
			X: [][]float64{
				[]float64{math.NaN(), -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				GP: &GP{LossMetric: metrics.Accuracy{}},
			},
			proba:       false,
			y:           nil,
			raisesError: true,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var y, err = tc.program.Predict(tc.X, tc.proba)
			if (err != nil) != tc.raisesError {
				t.Errorf("Expected nil error, got %s", err)
			}
			for j := range y {
				if math.Abs(y[j]-tc.y[j]) > 10e-5 {
					t.Errorf("Expected %.5f, got %.5f", tc.y[j], y[j])
				}
			}
		})
	}
}

func TestProgramMarshalJSON(t *testing.T) {
	var (
		prog = Program{
			Op: op.Add{op.Var{0}, op.Const{42}},
			GP: &GP{LossMetric: metrics.LogLoss{}},
		}
		bytes, err = prog.MarshalJSON()
	)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	var newProg = Program{}
	err = newProg.UnmarshalJSON(bytes)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	if newProg.String() != prog.String() {
		t.Errorf("Expected %s, got %s", prog, newProg)
		return
	}
}
