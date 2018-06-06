package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/op"
)

// A Mutator takes an Operator and returns a modified version of it.
type Mutator interface {
	Apply(operator op.Operator, rng *rand.Rand) op.Operator
}

// PointMutation randomly replaces Operators.
type PointMutation struct {
	Rate   float64
	Mutate func(op op.Operator, rng *rand.Rand) op.Operator
}

// Apply PointMutation.
func (pm PointMutation) Apply(operator op.Operator, rng *rand.Rand) op.Operator {
	for i := uint(0); i < operator.Arity(); i++ {
		operator = operator.SetOperand(i, pm.Apply(operator.Operand(i), rng))
	}
	if rng.Float64() < pm.Rate {
		operator = pm.Mutate(operator, rng)
	}
	return operator
}

// HoistMutation replaces an Operator by of it's operands.
type HoistMutation struct {
	Weight1, Weight2 func(operator op.Operator, depth uint, rng *rand.Rand) float64
}

// Apply HoistMutation.
func (hm HoistMutation) Apply(operator op.Operator, rng *rand.Rand) op.Operator {
	var (
		subOp, pos  = op.Sample(operator, hm.Weight1, rng)
		subSubOp, _ = op.Sample(subOp, hm.Weight2, rng)
	)
	return op.Replace(operator, pos, subSubOp)
}

// SubtreeMutation selects a suboperator at random and replaces it with a new
// Operator.
type SubtreeMutation struct {
	Weight      func(operator op.Operator, depth uint, rng *rand.Rand) float64
	NewOperator func(rng *rand.Rand) op.Operator
}

// Apply SubtreeMutation.
func (sm SubtreeMutation) Apply(operator op.Operator, rng *rand.Rand) op.Operator {
	var _, pos = op.Sample(operator, sm.Weight, rng)
	return op.Replace(operator, pos, sm.NewOperator(rng))
}
