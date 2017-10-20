package xgp

import (
	"testing"

	"github.com/MaxHalford/xgp/tree"
)

func TestNewProgramTuner(t *testing.T) {
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
	// Evaluate the Program
	var y, _ = progTuner.Program.PredictRow([]float64{})
	if y != 3.0 {
		t.Errorf("Expected %f, got %f", 3.0, y)
	}
	// Check the ProgramTuner's fields
	if len(progTuner.ConstValues) != 2 {
		t.Errorf("Expected %d, got %d", 2, len(progTuner.ConstValues))
	}
	// Check the ProgramTuner's fields
	if len(progTuner.ConstSetters) != 2 {
		t.Errorf("Expected %d, got %d", 2, len(progTuner.ConstSetters))
	}
	// Apply the ConstSetters and evaluate the Program
	progTuner.ConstSetters[0](2)
	y, _ = progTuner.Program.PredictRow([]float64{})
	if y != 4.0 {
		t.Errorf("Expected %f, got %f", 4.0, y)
	}
	progTuner.ConstSetters[1](3)
	y, _ = progTuner.Program.PredictRow([]float64{})
	if y != 5.0 {
		t.Errorf("Expected %f, got %f", 5.0, y)
	}
	// Check the initial Program hasn't been modified
	y, _ = prog.PredictRow([]float64{})
	if y != 3.0 {
		t.Errorf("Expected %f, got %f", 3.0, y)
	}
}

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
