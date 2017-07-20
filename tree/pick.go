package tree

import "math/rand"

type Picker interface {
	Apply(tree Tree, rng *rand.Rand) Tree
}

// If the Node has children then it is picked with probability P. If it isn't
// picked then pickNode is called on one of it's children.
type BernoulliPicker struct {
	P float64
}

// Apply BernoulliPicker.
func (picker BernoulliPicker) Apply(tree Tree, rng *rand.Rand) Tree {
	switch nBranches := tree.NBranches(); nBranches {
	case 0:
		return tree
	default:
		if rng.Float64() < picker.P {
			return tree
		}
		return picker.Apply(tree.GetBranch(rng.Intn(nBranches)), rng)
	}
}
