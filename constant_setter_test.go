package koza

import (
	"testing"

	"github.com/MaxHalford/koza/tree"
	"github.com/MaxHalford/koza/tree/op"
)

func TestNewConstantSetter(t *testing.T) {
	var (
		tr          = &tree.Tree{Operator: op.Constant{1}}
		constSetter = newConstantSetter(tr)
	)
	constSetter(2)
	if c, ok := tr.Operator.(op.Constant); (!ok || c != op.Constant{2}) {
		t.Errorf("Expected %v, got %v", op.Constant{2}, c)
	}
}
