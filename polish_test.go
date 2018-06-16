package xgp

import (
	"fmt"
	"testing"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
)

func TestPolish(t *testing.T) {
	var prog = Program{
		Op: op.Add{op.Const{2}, op.Mul{op.Const{3}, op.Var{0}}},
		Estimator: &Estimator{
			XTrain: [][]float64{
				[]float64{1, 2, 3, 4, 5},
			},
			YTrain:     []float64{8, 13, 18, 23, 28},
			LossMetric: metrics.MeanSquaredError{},
		},
	}
	prog, err := polishProgram(prog)
	fmt.Println(prog, err)
}
