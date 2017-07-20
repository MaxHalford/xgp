package tree

import "math/rand"

type Crossover interface {
	Apply(tree1, tree2 Tree, rng *rand.Rand)
}

type SubtreeCrossover struct {
	Picker Picker
}

func (sc SubtreeCrossover) Apply(tree1, tree2 Tree, rng *rand.Rand) {
	var (
		subtree1 = sc.Picker.Apply(tree1, rng)
		subtree2 = sc.Picker.Apply(tree2, rng)
	)
	subtree1.Swap(subtree2)
}
