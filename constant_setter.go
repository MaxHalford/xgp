package koza

import (
	"github.com/MaxHalford/koza/tree"
	"github.com/MaxHalford/koza/tree/op"
)

// A ConstantSetter can replace a tree's Operator with a Constant.
type ConstantSetter func(value float64)

func newConstantSetter(t *tree.Tree) ConstantSetter {
	return func(value float64) {
		t.Operator = op.Constant{Value: value}
	}
}
