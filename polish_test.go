package xgp

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
)

func TestPolish(t *testing.T) {
	var (
		rng       = rand.New(rand.NewSource(42))
		testCases = []struct {
			prog   Program
			consts []float64
		}{
			{
				prog: Program{
					Op: op.Add{op.Const{2}, op.Mul{op.Const{3}, op.Var{0}}},
					Estimator: &Estimator{
						XTrain: [][]float64{
							[]float64{1, 2, 3, 4, 5},
						},
						YTrain:     []float64{8, 13, 18, 23, 28},
						LossMetric: metrics.MeanSquaredError{},
					},
				},
				consts: []float64{3, 5},
			},
			{
				prog: Program{
					Op: op.Add{op.Var{2}, op.Mul{op.Var{1}, op.Var{0}}},
					Estimator: &Estimator{
						XTrain: [][]float64{
							[]float64{1, 2, 3, 4, 5},
						},
						YTrain:     []float64{8, 13, 18, 23, 28},
						LossMetric: metrics.MeanSquaredError{},
					},
				},
				consts: []float64{},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			prog, err := polishProgram(tc.prog, rng)
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
			consts := op.GetConsts(prog.Op)
			for i, c := range consts {
				if math.Abs(c-tc.consts[i]) > 0.00001 {
					t.Errorf("Expected %v, got %v", tc.consts, consts)
				}
			}
		})
	}
}
