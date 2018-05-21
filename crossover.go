package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/tree"
)

// A Crossover mixes two Programs in-place.
type Crossover interface {
	Apply(tree1, tree2 *tree.Tree, rng *rand.Rand)
}

// SubtreeCrossover applies subtree crossover to two Tree.
type SubtreeCrossover struct {
	Picker Picker
}

// Apply SubtreeCrossover.
func (cross SubtreeCrossover) Apply(tree1, tree2 *tree.Tree, rng *rand.Rand) {
	var (
		subTree1 = cross.Picker.Apply(tree1, 0, tree1.Height()-1, rng)
		subTree2 = cross.Picker.Apply(tree2, 0, tree2.Height(), rng)
	)
	*subTree1, *subTree2 = *subTree2, *subTree1
}
