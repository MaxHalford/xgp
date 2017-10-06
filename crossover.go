package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/tree"
)

// A Crossover mixes two Programs in-place.
type Crossover interface {
	Apply(prog1, prog2 *Program, rng *rand.Rand)
}

// SubTreeCrossover applies sub-tree crossover to two Programs.
type SubTreeCrossover struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

// Apply SubTreeCrossover.
func (cross SubTreeCrossover) Apply(prog1, prog2 *Program, rng *rand.Rand) {
	var (
		weight = func(tree tree.Tree) float64 {
			switch tree.(*Node).Operator.(type) {
			case Constant:
				return cross.PConstant
			case Variable:
				return cross.PVariable
			default:
				return cross.PFunction
			}
		}
		subTree1, _ = tree.PickSubTree(prog1.Root, weight, 0, -1, rng)
		subTree2, _ = tree.PickSubTree(prog2.Root, weight, 0, -1, rng)
	)
	subTree1.(*Node).Swap(subTree2.(*Node))
}
