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
	Picker         Picker
	MutateOperator func(op Operator, rng *rand.Rand) Operator
}

// Apply PointMutation.
func (mut PointMutation) Apply(tree *Tree, rng *rand.Rand) {
	var subTree, _ = mut.Picker.Apply(tree, 0, -1, rng)
	subTree.Operator = mut.MutateOperator(subTree.Operator, rng)
}

// HoistMutation selects a first sub-Tree from a Tree. It then selects a second
// sub-Tree from the first sub-Tree and replaces the first one with it. Hoist
// mutation is good for controlling bloat.
type HoistMutation struct {
	Picker Picker
}

// Apply HoistMutation.
func (mut HoistMutation) Apply(tree *Tree, rng *rand.Rand) {
	// The Tree has to have a depth of at least 2 for hoist mutation to make
	// sense
	var height = tree.Height()
	if height < 2 {
		return
	}
	var (
		sub, _    = mut.Picker.Apply(tree, 1, height-1, rng)
		subsub, _ = mut.Picker.Apply(sub, 1, -1, rng)
	)
	*sub = *subsub
}

// SubTreeMutation selects a sub-Tree at random and replaces with a new Tree.
// The new Tree has at most the same height as the selected sub-Tree.
type SubTreeMutation struct {
	Picker  Picker
	NewTree func(minHeight, maxHeight int, rng *rand.Rand) *Tree
}

// Apply SubTreeMutation.
func (mut SubTreeMutation) Apply(tree *Tree, rng *rand.Rand) {
	var (
		sub, _  = mut.Picker.Apply(tree, 1, tree.Height(), rng)
		newTree = mut.NewTree(0, sub.Height(), rng)
	)
	*sub = *newTree
}
