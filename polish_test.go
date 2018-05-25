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

func TestSetProgConstants(t *testing.T) {
	var (
		prog = Program{
			Tree: tree.MustParseCode("sum(0, 1)"),
		}
		progPolish = newProgramPolish(prog)
	)
	// Set new Constants
	for i, c := range progPolish.ConstValues {
		progPolish.ConstValues[i] = c + 1
	}
	progPolish.setProgConstants()
	var expected = tree.MustParseCode("sum(1, 2)")
	if progPolish.Program.Tree.String() != expected.String() {
		t.Errorf("Expected %v, got %v", expected, progPolish.Program.Tree)
	}
}
