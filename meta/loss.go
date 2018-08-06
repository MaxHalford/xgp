package meta

import (
	"math"

	"github.com/MaxHalford/xgp/op"
)

type Loss interface {
	Eval(yTrue, yPred []float64) []float64
	GradEval(yTrue, yPred []float64) []float64
	Classification() bool
}

type SquareLoss struct{}

// Eval computes 0.5 * (yTrue - yPred) ** 2
func (sq SquareLoss) Eval(yTrue, yPred []float64) []float64 {
	var ys = [][]float64{yTrue, yPred}
	return op.Mul{op.Const{0.5}, op.Square{op.Sub{op.Var{0}, op.Var{1}}}}.Eval(ys)
}

// GradEval computes yPred - yTrue
func (sq SquareLoss) GradEval(yTrue, yPred []float64) []float64 {
	var ys = [][]float64{yPred, yTrue}
	return op.Sub{op.Var{0}, op.Var{1}}.Eval(ys)
}

// Classification return false.
func (sq SquareLoss) Classification() bool {
	return false
}

type LogLoss struct{}

// Eval computes -(yTrue * log(yPred) + (1 - yTrue) * log(1 - yPred))
func (ll LogLoss) Eval(yTrue, yPred []float64) []float64 {
	var loss = make([]float64, len(yTrue))
	for i, y := range yTrue {
		loss[i] = -(y*math.Log(yPred[i]) + (1-y)*math.Log(1-yPred[i]))
	}
	return loss
}

// GradEval computes yTrue - yPred.
func (ll LogLoss) GradEval(yTrue, yPred []float64) []float64 {
	return op.Sub{op.Var{0}, op.Var{1}}.Eval([][]float64{yTrue, yPred})
}

// Classification return true.
func (ll LogLoss) Classification() bool {
	return true
}
