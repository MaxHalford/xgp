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
		progTuner = newProgramTuner(prog)
	)
	// Set new Constants
	for i, c := range progTuner.ConstValues {
		progTuner.ConstValues[i] = c + 1
	}
	progTuner.setProgConstants()
	var expected = tree.MustParseCode("sum(1, 2)")
	if progTuner.Program.Tree.String() != expected.String() {
		t.Errorf("Expected %v, got %v", expected, progTuner.Program.Tree)
	}
}
