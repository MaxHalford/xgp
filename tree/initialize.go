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

// FullInitializer generates trees where all the leaves are the same depth.
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

// GrowInitializer generates trees where all the leaves have at most a
// certain depth.
type GrowInitializer struct {
	MaxHeight int
	PLeaf     float64 // Probability of producing a leaf
}

// Apply GrowInitializer.
func (init GrowInitializer) Apply(of OperatorFactory, rng *rand.Rand) *Tree {
	var (
		op   = of.New(init.MaxHeight == 0 || rng.Float64() < init.PLeaf, rng)
		tree = &Tree{
			Operator: op,
			Branches: make([]*Tree, op.Arity()),
		}
	)
	for i := range tree.Branches {
		tree.Branches[i] = GrowInitializer{
			MaxHeight: init.MaxHeight - 1,
			PLeaf:     init.PLeaf,
		}.Apply(of, rng)
	}
	return tree
}

// RampedHaldAndHalfInitializer randomly chooses GrowtreeInitializer or
// FulltreeInitializer with a random height in [MinHeight, MaxHeight].
type RampedHaldAndHalfInitializer struct {
	MinHeight int
	MaxHeight int
	PLeaf     float64 // Probability of producing a leaf for GrowtreeInitializer
}

// Apply RampedHaldAndHalfInitializer.
func (init RampedHaldAndHalfInitializer) Apply(of OperatorFactory, rng *rand.Rand) *Tree {
	// Randomly pick a height
	var height = randInt(init.MinHeight, init.MaxHeight, rng)
	// Randomly apply full initialization or grow initialization
	if rng.Float64() < 0.5 {
		return FullInitializer{Height: height}.Apply(of, rng)
	}
	return GrowInitializer{
		MaxHeight: height,
		PLeaf:     init.PLeaf,
	}.Apply(of, rng)
}
