package xgp

import (
	"testing"

	"github.com/MaxHalford/xgp/tree"
)

func TestSetProgConstants(t *testing.T) {
	var (
		prog = Program{
			Tree: &tree.Tree{
				Operator: tree.Sum{},
				Branches: []*tree.Tree{
					&tree.Tree{Operator: tree.Constant{1}},
					&tree.Tree{Operator: tree.Constant{2}},
				},
			},
		}
		progTuner = newProgramTuner(&prog)
	)
	// Set new Constants
	for i, c := range progTuner.ConstValues {
		progTuner.ConstValues[i] = c + 1
	}
	progTuner.setProgConstants()
	// Check with the Program's Constants
	for i, branch := range progTuner.Program.Tree.Branches {
		if branch.Operator.(tree.Constant).Value != prog.Tree.Branches[i].Operator.(tree.Constant).Value+1 {
			t.Errorf("Expected %v, got %v", prog.Tree.Branches[i], branch.Operator)
		}
	}
}

func TestJitterConstants(t *testing.T) {
	var (
		prog = Program{
			Tree: &tree.Tree{
				Operator: tree.Sum{},
				Branches: []*tree.Tree{
					&tree.Tree{Operator: tree.Constant{1}},
					&tree.Tree{Operator: tree.Constant{2}},
				},
			},
		}
		progTuner = newProgramTuner(&prog)
	)
	// Jitter Constants
	progTuner.jitterConstants(newRand())
	// Compare with the Program's Constants
	for i, c := range progTuner.ConstValues {
		if c == prog.Tree.Branches[i].Operator.(tree.Constant).Value {
			t.Errorf("Expected %v and %v to be different", prog.Tree.Branches[i], c)
		}
	}
}
