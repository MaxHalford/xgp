package xgp

import (
	"testing"

	"github.com/MaxHalford/xgp/tree"
)

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
