package xgp

import "math/rand"

// randProg is a convenience method that produces a random Program for testing
// purposes.
func randProg(rng *rand.Rand) Program {
	var (
		nodeInit    = FullNodeInitializer{Height: randInt(3, 5, rng)}
		functions   = []Operator{Cos{}, Sin{}, Sum{}, Difference{}, Product{}, Division{}}
		newOperator = func(terminal bool, rng *rand.Rand) Operator {
			if terminal {
				if rng.Float64() < 0.5 {
					return newVariable(3, rng)
				}
				return newConstant(-10, 10, rng)
			}
			return newFunction(functions, rng)
		}
	)
	return Program{
		Root: nodeInit.Apply(newOperator, rng),
	}
}
