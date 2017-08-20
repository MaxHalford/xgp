package xgp

import (
	"testing"
)

func TestNewConstantSetter(t *testing.T) {
	var (
		node        = &Node{Operator: Constant{1}}
		constSetter = node.newConstantSetter()
	)
	constSetter(2)
	if c, ok := node.Operator.(Constant); (!ok || c != Constant{2}) {
		t.Errorf("Expected %v, got %v", Constant{2}, c)
	}
}
