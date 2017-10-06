package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/tree"
)

// A Mutator modifies a Program in-place.
type Mutator interface {
	Apply(prog *Program, rng *rand.Rand)
}

// PointMutation applies point mutation to a Program.
type PointMutation struct {
	P float64
}

// Apply PointMutation.
func (mut PointMutation) Apply(prog *Program, rng *rand.Rand) {
	var f = func(node *Node) {
		if rng.Float64() < mut.P {
			switch node.Operator.(type) {
			case Constant:
				node.setOperator(prog.Estimator.newConstant(rng), rng)
			case Variable:
				node.setOperator(prog.Estimator.newVariable(rng), rng)
			default:
				node.setOperator(prog.Estimator.newFunctionOfArity(node.Operator.Arity(), rng), rng)
			}
		}
	}
	prog.Root.RecApply(f)
}

// HoistMutation applies hoist mutation to a Program.
type HoistMutation struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

// Apply HoistMutation.
func (mut HoistMutation) Apply(prog *Program, rng *rand.Rand) {
	var (
		weight = func(tree tree.Tree) float64 {
			switch tree.(*Node).Operator.(type) {
			case Constant:
				return mut.PConstant
			case Variable:
				return mut.PVariable
			default:
				return mut.PFunction
			}
		}
		sub, _    = tree.PickSubTree(prog.Root, weight, 1, -1, rng)
		subsub, _ = tree.PickSubTree(sub, weight, 1, -1, rng)
	)
	*sub.(*Node) = *subsub.(*Node)
}
