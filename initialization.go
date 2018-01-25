package koza

import (
	"math/rand"

	"github.com/MaxHalford/koza/op"
	"github.com/MaxHalford/koza/tree"
)

// An OperatorFactory produces new Operators.
type OperatorFactory struct {
	PConstant   float64
	NewConstant func(rng *rand.Rand) op.Constant
	NewVariable func(rng *rand.Rand) op.Variable
	NewFunction func(rng *rand.Rand) op.Operator
}

// New returns an Operator.
func (of OperatorFactory) New(terminal bool, rng *rand.Rand) op.Operator {
	if terminal {
		if rng.Float64() < of.PConstant {
			return of.NewConstant(rng)
		}
		return of.NewVariable(rng)
	}
	return of.NewFunction(rng)
}

// A Initializer generates a random Tree.
type Initializer interface {
	Apply(minHeight, maxHeight int, of OperatorFactory, rng *rand.Rand) tree.Tree
}

// FullInitializer generates Trees with node depths equal to Height.
type FullInitializer struct{}

// Apply FullInitializer.
func (init FullInitializer) Apply(minHeight, maxHeight int, of OperatorFactory, rng *rand.Rand) tree.Tree {
	var (
		op = of.New(maxHeight == 0, rng)
		tr = tree.NewTree(op)
	)
	for i := 0; i < tr.NBranches(); i++ {
		tr.SetBranch(i, init.Apply(0, maxHeight-1, of, rng))
	}
	return tr
}

// GrowInitializer generates Trees with node depths in [MinHeight, MaxHeight].
type GrowInitializer struct {
	PTerminal float64
}

// Apply GrowInitializer.
func (init GrowInitializer) Apply(minHeight, maxHeight int, of OperatorFactory, rng *rand.Rand) tree.Tree {
	var (
		leaf = maxHeight == 0 || (minHeight <= 0 && rng.Float64() < init.PTerminal)
		op   = of.New(leaf, rng)
		tr   = tree.NewTree(op)
	)
	for i := 0; i < tr.NBranches(); i++ {
		tr.SetBranch(i, init.Apply(minHeight-1, maxHeight-1, of, rng))
	}
	return tr
}

// RampedHaldAndHalfInitializer randomly applies GrowTreeInitializer or
// FullTreeInitializer.
type RampedHaldAndHalfInitializer struct {
	PFull           float64
	FullInitializer FullInitializer
	GrowInitializer GrowInitializer
}

// Apply RampedHaldAndHalfInitializer.
func (init RampedHaldAndHalfInitializer) Apply(minHeight, maxHeight int, of OperatorFactory, rng *rand.Rand) tree.Tree {
	if rng.Float64() < init.PFull {
		var height = randInt(minHeight, maxHeight, rng)
		return init.FullInitializer.Apply(0, height, of, rng)
	}
	return init.GrowInitializer.Apply(minHeight, maxHeight, of, rng)
}
