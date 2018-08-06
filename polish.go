package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/xgp/op"
	xrand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/optimize"
)

func polishProgram(prog Program, rng *rand.Rand) (Program, error) {
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
					Op: op.SetConsts(prog.Op, x),
					GP: prog.GP,
				}.Evaluate()
				return fitness
			},
		}
		method = &optimize.CmaEsChol{
			Population: int(15 + math.Floor(3*math.Log(float64(len(consts))))), // MAGIC
			Src:        xrand.NewSource(rng.Uint64()),
		}
	)

	// Run the optimisation
	result, err := optimize.Minimize(problem, consts, nil, method)
	if err != nil {
		return prog, err
	}

	prog.Op = op.SetConsts(prog.Op, result.X)
	return prog, nil
}
