package koza

import (
	"math/rand"

	"github.com/MaxHalford/koza/tree"
)

// A Crossover mixes two Programs in-place.
type Crossover interface {
	Apply(tree1, tree2 *tree.Tree, rng *rand.Rand)
}

// SubTreeCrossover applies sub-tree crossover to two Tree.
type SubTreeCrossover struct {
	Picker Picker
}

// Apply SubTreeCrossover.
func (cross SubTreeCrossover) Apply(tree1, tree2 *tree.Tree, rng *rand.Rand) {
	var (
		subTree1 = cross.Picker.Apply(tree1, 0, tree1.Height()-1, rng)
		subTree2 = cross.Picker.Apply(tree2, 0, tree2.Height(), rng)
	)
	*subTree1, *subTree2 = *subTree2, *subTree1
}
