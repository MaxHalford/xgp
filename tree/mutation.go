package tree

import (
	"math/rand"
)

// A Mutator modifies a Tree in-place.
type Mutator interface {
	Apply(tree *Tree, rng *rand.Rand)
}

// PointMutation picks one sub-Tree at random and replaces it's Operator.
type PointMutation struct {
	Weighting      Weighting
	MutateOperator func(op Operator, rng *rand.Rand) Operator
}

// Apply PointMutation.
func (mut PointMutation) Apply(tree *Tree, rng *rand.Rand) {
	var f = func(tree *Tree, depth int) (stop bool) {
		if rng.Float64() < mut.Weighting.apply(tree.Operator) {
			tree.Operator = mut.MutateOperator(tree.Operator, rng)
		}
		return false
	}
	tree.Walk(f)
}

// HoistMutation selects a first sub-Tree from a Tree. It then selects a second
// sub-Tree from the first sub-Tree and replaces the first one with it. Hoist
// mutation is good for controlling bloat.
type HoistMutation struct {
	Picker Picker
}

// Apply HoistMutation.
func (mut HoistMutation) Apply(tree *Tree, rng *rand.Rand) {
	// Hoist mutation only works if the height of Tree exceeds 1
	var height = tree.Height()
	if height < 1 {
		return
	}
	var (
		sub    = mut.Picker.Apply(tree, 1, tree.Height(), rng)
		subsub = mut.Picker.Apply(sub, 0, sub.Height()-1, rng)
	)
	*sub = *subsub
}

// SubTreeMutation selects a sub-Tree at random and replaces with a new Tree.
// The new Tree has at most the same height as the selected sub-Tree.
type SubTreeMutation struct {
	NewTree   func(rng *rand.Rand) *Tree
	Crossover Crossover
}

// Apply SubTreeMutation.
func (mut SubTreeMutation) Apply(tree *Tree, rng *rand.Rand) {
	mut.Crossover.Apply(tree, mut.NewTree(rng), rng)
}
