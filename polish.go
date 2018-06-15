package xgp

import (
	"errors"
	"math"

	"github.com/MaxHalford/xgp/op"
	"gonum.org/v1/gonum/optimize"
)

func meanFloat64s(fs []float64) (mean float64) {
	for _, f := range fs {
		mean += f
	}
	return mean / float64(len(fs))
}

func polishProgram(prog Program) (Program, error) {
	// Extract the Program's Consts
	var consts = op.GetConsts(prog.Op)

	// If there are no Consts then nothing can be done
	if len(consts) == 0 {
		return prog, nil
	}

	// Determine the loss function
	var (
		nVars = op.CountVars(prog.Op)
		loss  = op.Square{op.Sub{prog.Op, op.Var{nVars}}}
	)
	nVars++

	// Replace the Consts by Vars in the loss
	var (
		when = func(operator op.Operator) bool {
			_, ok := operator.(op.Const)
			return ok
		}
		nConsts uint
		how     = func(operator op.Operator) op.Operator {
			defer func() { nConsts++ }()
			return op.Var{nVars + nConsts}
		}
	)
	loss, ok := op.Replace(loss, when, how, false).Simplify().(op.Square)
	if !ok {
		return prog, errors.New("Bad type conversion")
	}

	// Build the loss derivatives
	var derivatives = make([]op.Operator, nConsts)
	for i := uint(0); i < nConsts; i++ {
		derivatives[i] = loss.Diff(nVars + i).(op.Mul).Right.Simplify()
	}

	// Replace the Consts by Vars in each derivative
	when = func(operator op.Operator) bool {
		v, ok := operator.(op.Var)
		return ok && v.Index >= nVars
	}
	how = func(operator op.Operator) op.Operator {
		defer func() { nConsts++ }()
		return op.Const{consts[nConsts]}
	}
	for i := range derivatives {
		nConsts = 0
		derivatives[i] = op.Replace(derivatives[i], when, how, false)
	}

	// Replace the Consts by Vars in the loss
	nConsts = 0
	loss = op.Replace(loss, when, how, false).(op.Square)

	// Define the function to minimize and the associated gradient
	var m = make([][]float64, len(prog.Estimator.XTrain)+1)
	for i, xi := range prog.Estimator.XTrain {
		m[i] = xi
	}
	m[len(m)-1] = prog.Estimator.YTrain
	var (
		p = optimize.Problem{
			Func: func(x []float64) float64 {
				return math.Pow(meanFloat64s(op.SetConsts(loss, x).Eval(m)), 2)
			},
			Grad: func(grad []float64, x []float64) {
				for i, d := range derivatives {
					grad[i] = meanFloat64s(op.SetConsts(d, x).Eval(m))
				}
			},
		}
	)

	// Run the optimisation
	result, err := optimize.Local(p, consts, nil, nil)
	if err != nil {
		return prog, err
	}

	// Return the Program with the best Consts
	return Program{Op: op.SetConsts(prog.Op, result.X)}, nil
}
