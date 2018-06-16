package xgp

import (
	"github.com/MaxHalford/xgp/op"
	"gonum.org/v1/gonum/optimize"
)

func polishProgram(prog Program) (Program, error) {
	// Extract the Program's Consts
	var consts = op.GetConsts(prog.Op)

	// If there are no Consts then nothing can be done
	if len(consts) == 0 {
		return prog, nil
	}

	var (
		problem = optimize.Problem{
			Func: func(x []float64) float64 {
				fitness, _ := Program{
					Op:        op.SetConsts(prog.Op, x),
					Estimator: prog.Estimator,
				}.Evaluate()
				return fitness
			},
		}
		method = &optimize.CmaEsChol{InitMean: consts}
	)

	// Run the optimisation
	result, err := optimize.Global(problem, len(consts), nil, method)
	if err != nil {
		return prog, err
	}

	prog.Op = op.SetConsts(prog.Op, result.X)
	return prog, nil
}
