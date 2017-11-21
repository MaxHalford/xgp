package tree

import (
	"math/rand"
)

// A Crossover mixes two Programs in-place.
type Crossover interface {
	Apply(tree1, tree2 *Tree, rng *rand.Rand)
}

// SubTreeCrossover applies sub-tree crossover to two Tree.
type SubTreeCrossover struct {
	Picker Picker
}

// Apply SubTreeCrossover.
func (cross SubTreeCrossover) Apply(tree1, tree2 *Tree, rng *rand.Rand) {
	// var (
	// 	subTree1, _ = cross.Picker.Apply(tree1, 0, -1, rng)
	// 	subTree2, _ = cross.Picker.Apply(tree2, 0, -1, rng)
	// )
	// *subTree1, *subTree2 = *subTree2, *subTree1
}
