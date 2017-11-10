package tree

import "math/rand"

// A Mutator modifies a Program in-place.
type Mutator interface {
	Apply(tree *Tree, rng *rand.Rand)
}

// PointMutation applies point mutation to a Program.
type PointMutation struct {
	NewOperator func(op Operator, rng *rand.Rand) Operator
	P           float64
}

// Apply PointMutation.
func (mut PointMutation) Apply(tree *Tree, rng *rand.Rand) {
	var f = func(tree *Tree, depth int) (stop bool) {
		if rng.Float64() < mut.P {
			tree.Operator = mut.NewOperator(tree.Operator, rng)
		}
		return
	}
	tree.rApply(f)
}

// HoistMutation applies hoist mutation to a Program.
type HoistMutation struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

// Apply HoistMutation.
func (mut HoistMutation) Apply(tree *Tree, rng *rand.Rand) {
	var (
		weight = func(tree *Tree) float64 {
			switch tree.Operator.(type) {
			case Constant:
				return mut.PConstant
			case Variable:
				return mut.PVariable
			default:
				return mut.PFunction
			}
		}
		sub, _    = pickSubTree(*tree, weight, 1, -1, rng)
		subsub, _ = pickSubTree(*sub, weight, 1, -1, rng)
	)
	*sub = *subsub
}
