package tree

import "math/rand"

// An OperatorFactory produces new Operators.
type OperatorFactory struct {
	PVariable   float64
	NewConstant func(rng *rand.Rand) Constant
	NewVariable func(rng *rand.Rand) Variable
	NewFunction func(rng *rand.Rand) Operator
}

// New returns an Operator.
func (of OperatorFactory) New(terminal bool, rng *rand.Rand) Operator {
	if terminal {
		if rng.Float64() < of.PVariable {
			return of.NewVariable(rng)
		}
		return of.NewConstant(rng)
	}
	return of.NewFunction(rng)
}

// A Initializer generates a random Tree.
type Initializer interface {
	Apply(of OperatorFactory, rng *rand.Rand) *Tree
}

// FullInitializer generates Trees with node depths equal to Height.
type FullInitializer struct {
	Height int
}

// Apply FullInitializer.
func (init FullInitializer) Apply(of OperatorFactory, rng *rand.Rand) *Tree {
	var (
		op   = of.New(init.Height == 0, rng)
		tree = &Tree{
			Operator: op,
			Branches: make([]*Tree, op.Arity()),
		}
	)
	for i := range tree.Branches {
		tree.Branches[i] = FullInitializer{Height: init.Height - 1}.Apply(of, rng)
	}
	return tree
}

// GrowInitializer generates Trees with node depths in [MinHeight, MaxHeight].
type GrowInitializer struct {
	MinHeight int
	MaxHeight int
	PLeaf     float64
}

// Apply GrowInitializer.
func (init GrowInitializer) Apply(of OperatorFactory, rng *rand.Rand) *Tree {
	var (
		leaf = init.MinHeight <= 0 && (init.MaxHeight == 0 || rng.Float64() < init.PLeaf)
		op   = of.New(leaf, rng)
		tree = &Tree{
			Operator: op,
			Branches: make([]*Tree, op.Arity()),
		}
	)
	for i := range tree.Branches {
		tree.Branches[i] = GrowInitializer{
			MinHeight: init.MinHeight - 1,
			MaxHeight: init.MaxHeight - 1,
			PLeaf:     init.PLeaf,
		}.Apply(of, rng)
	}
	return tree
}

// RampedHaldAndHalfInitializer randomly applies GrowTreeInitializer or
// FullTreeInitializer.
type RampedHaldAndHalfInitializer struct {
	MinHeight int
	MaxHeight int
	PLeaf     float64 // Probability of producing a leaf for GrowtreeInitializer
}

// Apply RampedHaldAndHalfInitializer.
func (init RampedHaldAndHalfInitializer) Apply(of OperatorFactory, rng *rand.Rand) *Tree {
	if rng.Float64() < 0.5 {
		return FullInitializer{
			Height: randInt(init.MinHeight, init.MaxHeight, rng),
		}.Apply(of, rng)
	}
	return GrowInitializer{
		MinHeight: init.MinHeight,
		MaxHeight: init.MaxHeight,
		PLeaf:     init.PLeaf,
	}.Apply(of, rng)
}
