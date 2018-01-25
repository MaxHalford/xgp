package koza

import (
	"github.com/MaxHalford/koza/op"
	"github.com/MaxHalford/koza/tree"
)

// A ConstantSetter can replace a tree's Operator with a Constant.
type ConstantSetter func(value float64)

func newConstantSetter(tr *tree.Tree) ConstantSetter {
	return func(value float64) {
		tr.SetOperator(op.Constant{Value: value})
	}
}
