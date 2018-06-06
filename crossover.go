package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/op"
)

// A Crossover takes two Operators and combines them in order to produce two
// new Operators.
type Crossover interface {
	Apply(op1, op2 op.Operator, rng *rand.Rand) (op.Operator, op.Operator)
}

// SubtreeCrossover applies subtree crossover to two Operators.
type SubtreeCrossover struct {
	Weight func(operator op.Operator, depth uint, rng *rand.Rand) float64
}

// Apply SubtreeCrossover.
func (sc SubtreeCrossover) Apply(op1, op2 op.Operator, rng *rand.Rand) (op.Operator, op.Operator) {
	var (
		sub1, i = op.Sample(op1, sc.Weight, rng)
		sub2, j = op.Sample(op2, sc.Weight, rng)
	)
	return op.Replace(op1, i, sub2), op.Replace(op2, j, sub1)
}
