package tree

import (
	"math/rand"
)

// A Crossover mixes two Programs in-place.
type Crossover interface {
	Apply(left, right *Tree, rng *rand.Rand)
}

// SubTreeCrossover applies sub-tree crossover to two Tree.
type SubTreeCrossover struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

// Apply SubTreeCrossover.
func (cross SubTreeCrossover) Apply(left, right *Tree, rng *rand.Rand) {
	var (
		weight = func(tree Tree) float64 {
			switch tree.Operator.(type) {
			case Constant:
				return cross.PConstant
			case Variable:
				return cross.PVariable
			default:
				return cross.PFunction
			}
		}
		subTree1, _ = pickSubTree(*left, weight, 0, -1, rng)
		subTree2, _ = pickSubTree(*right, weight, 0, -1, rng)
	)
	*subTree1 = *subTree2
}
