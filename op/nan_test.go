package op

import (
	"math"
	"testing"
)

func TestNaNUnary(t *testing.T) {
	var (
		ops = []Operator{
			Cos{},
			Exp{},
			Sin{},
		}
		values = [][]float64{
			[]float64{-42.42},
			[]float64{-1},
			[]float64{0},
			[]float64{1},
		}
	)
	for _, op := range ops {
		for _, x := range values {
			if math.IsNaN(op.Eval([][]float64{x})[0]) {
				t.Errorf("%s(%f) is NaN", op.String(), x)
			}
			if math.IsInf(op.Eval([][]float64{x})[0], -1) {
				t.Errorf("%s(%f) is -∞", op.String(), x)
			}
			if math.IsInf(op.Eval([][]float64{x})[0], 1) {
				t.Errorf("%s(%f) is +∞", op.String(), x)
			}
		}
	}
}

func TestNaNBinary(t *testing.T) {
	var (
		ops = []Operator{
			Div{},
			Max{},
			Min{},
			Mul{},
			Sub{},
			Sum{},
		}
		values = [][]float64{
			[]float64{-42.42, -42.42},
			[]float64{-1, -1},
			[]float64{-1, 0},
			[]float64{-1, 1},
			[]float64{0, -1},
			[]float64{0, 0},
			[]float64{0, 1},
			[]float64{1, -1},
			[]float64{1, 0},
			[]float64{1, 1},
		}
	)
	for _, op := range ops {
		for _, x := range values {
			if math.IsNaN(op.Eval([][]float64{[]float64{x[0]}, []float64{x[1]}})[0]) {
				t.Errorf("%s(%f, %f) is NaN", op.String(), x[0], x[1])
			}
			if math.IsInf(op.Eval([][]float64{[]float64{x[0]}, []float64{x[1]}})[0], -1) {
				t.Errorf("%s(%f, %f) is -∞", op.String(), x[0], x[1])
			}
			if math.IsInf(op.Eval([][]float64{[]float64{x[0]}, []float64{x[1]}})[0], 1) {
				t.Errorf("%s(%f, %f) is +∞", op.String(), x[0], x[1])
			}
		}
	}
}
