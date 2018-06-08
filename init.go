package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/op"
)

// An Initializer generates a random Operator with random operands.
type Initializer interface {
	Apply(
		minHeight uint,
		maxHeight uint,
		newOp func(leaf bool, rng *rand.Rand) op.Operator,
		rng *rand.Rand,
	) op.Operator
}

// FullInit can generate an Operator of height maxHeight.
type FullInit struct{}

// Apply FullInit returns an Operator of height maxHeight, hence minHeight is
// not taken into account.
func (full FullInit) Apply(
	minHeight uint,
	maxHeight uint,
	newOp func(leaf bool, rng *rand.Rand) op.Operator,
	rng *rand.Rand,
) op.Operator {
	var op = newOp(maxHeight == 0, rng)
	for i := uint(0); i < op.Arity(); i++ {
		op = op.SetOperand(i, full.Apply(0, maxHeight-1, newOp, rng))
	}
	return op
}

// GrowInit can generate an Operator who's depth lies in [MinHeight, MaxHeight].
type GrowInit struct {
	PLeaf float64
}

// Apply GrowInit.
func (grow GrowInit) Apply(
	minHeight uint,
	maxHeight uint,
	newOp func(leaf bool, rng *rand.Rand) op.Operator,
	rng *rand.Rand,
) op.Operator {
	var (
		leaf = maxHeight == 0 || (minHeight <= 0 && rng.Float64() < grow.PLeaf)
		op   = newOp(leaf, rng)
	)
	for i := uint(0); i < op.Arity(); i++ {
		op = op.SetOperand(i, grow.Apply(minHeight-1, maxHeight-1, newOp, rng))
	}
	return op
}

// RampedHaldAndHalfInit randomly uses GrowInit and FullInit. If it uses
// FullInit then it uses a random height.
type RampedHaldAndHalfInit struct {
	PFull    float64
	FullInit FullInit
	GrowInit GrowInit
}

// Apply RampedHaldAndHalfInit.
func (rhah RampedHaldAndHalfInit) Apply(
	minHeight uint,
	maxHeight uint,
	newOp func(leaf bool, rng *rand.Rand) op.Operator,
	rng *rand.Rand,
) op.Operator {
	if rng.Float64() < rhah.PFull {
		var height = randInt(int(minHeight), int(maxHeight), rng)
		return rhah.FullInit.Apply(0, uint(height), newOp, rng)
	}
	return rhah.GrowInit.Apply(minHeight, maxHeight, newOp, rng)
}
