package xgp

import (
	"testing"

	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
)

func TestNewConstantSetter(t *testing.T) {
	var (
		tr          = tree.NewTree(op.Constant{1})
		constSetter = newConstantSetter(&tr)
	)
	constSetter(2)
	if c, ok := tr.Op.(op.Constant); (!ok || c != op.Constant{2}) {
		t.Errorf("Expected %v, got %v", op.Constant{2}, c)
	}
}
