package koza

import (
	"testing"

	"github.com/MaxHalford/koza/tree"
)

func TestNewConstantSetter(t *testing.T) {
	var (
		tr          = &tree.Tree{Operator: tree.Constant{1}}
		constSetter = newConstantSetter(tr)
	)
	constSetter(2)
	if c, ok := tr.Operator.(tree.Constant); (!ok || c != tree.Constant{2}) {
		t.Errorf("Expected %v, got %v", tree.Constant{2}, c)
	}
}
