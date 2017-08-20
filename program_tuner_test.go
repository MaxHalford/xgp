package xgp

import (
	"testing"
)

func TestNewProgramTuner(t *testing.T) {
	var (
		prog = Program{
			Root: &Node{
				Operator: Sum{},
				Children: []*Node{
					&Node{Operator: Constant{1}},
					&Node{Operator: Constant{2}},
				},
			},
		}
		progTuner = newProgramTuner(prog)
	)
	// Evaluate the Program
	var y = progTuner.Program.PredictRow([]float64{})
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
	y = progTuner.Program.PredictRow([]float64{})
	if y != 4.0 {
		t.Errorf("Expected %f, got %f", 4.0, y)
	}
	progTuner.ConstSetters[1](3)
	y = progTuner.Program.PredictRow([]float64{})
	if y != 5.0 {
		t.Errorf("Expected %f, got %f", 5.0, y)
	}
	// Check the initial Program hasn't been modified
	y = prog.PredictRow([]float64{})
	if y != 3.0 {
		t.Errorf("Expected %f, got %f", 3.0, y)
	}
}
