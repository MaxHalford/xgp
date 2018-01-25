package koza

import (
	"math/rand"

	"github.com/MaxHalford/koza/op"
	"github.com/MaxHalford/koza/tree"
)

// A Mutator modifies a Tree in-place.
type Mutator interface {
	Apply(tree *tree.Tree, rng *rand.Rand)
}

// PointMutation picks one sub-Tree at random and replaces it's Operator.
type PointMutation struct {
	Weighting      Weighting
	MutateOperator func(op op.Operator, rng *rand.Rand) op.Operator
}

// Apply PointMutation.
func (mut PointMutation) Apply(tr *tree.Tree, rng *rand.Rand) {
	var f = func(tr *tree.Tree, depth int) (stop bool) {
		if rng.Float64() < mut.Weighting.apply(tr.Operator()) {
			tr.SetOperator(mut.MutateOperator(tr.Operator(), rng))
		}
		return false
	}
	tr.Walk(f)
}

// HoistMutation selects a first sub-Tree from a Tree. It then selects a second
// sub-Tree from the first sub-Tree and replaces the first one with it. Hoist
// mutation is good for controlling bloat.
type HoistMutation struct {
	Picker Picker
}

// Apply HoistMutation.
func (mut HoistMutation) Apply(tr *tree.Tree, rng *rand.Rand) {
	// Hoist mutation only works if the height of Tree exceeds 1
	var height = tr.Height()
	if height < 1 {
		return
	}
	var (
		sub    = mut.Picker.Apply(tr, 1, tr.Height(), rng)
		subsub = mut.Picker.Apply(sub, 0, sub.Height()-1, rng)
	)
	*sub = *subsub
}

// SubTreeMutation selects a sub-Tree at random and replaces with a new Tree.
// The new Tree has at most the same height as the selected sub-Tree.
type SubTreeMutation struct {
	NewTree   func(rng *rand.Rand) tree.Tree
	Crossover Crossover
}

// Apply SubTreeMutation.
func (mut SubTreeMutation) Apply(tr *tree.Tree, rng *rand.Rand) {
	var mutant = mut.NewTree(rng)
	mut.Crossover.Apply(tr, &mutant, rng)
}
