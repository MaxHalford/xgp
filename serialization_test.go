package xgp

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/tree"
)

func TestDRSJSONEncodeDecode(t *testing.T) {
	var initialDRS = &DynamicRangeSelection{
		cutPoints: []float64{0, 1, 2},
		rangeMap:  map[float64]float64{0: 1, 1: 1, 2: 1},
	}

	// Serialize the initial DRS
	var bytes, err = json.Marshal(initialDRS)
	if err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
	}

	// Parse the bytes into a new tree
	var newDRS *DynamicRangeSelection
	err = json.Unmarshal(bytes, &newDRS)
	if err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
	}

	// Compare the DRSs
	if !reflect.DeepEqual(newDRS, initialDRS) {
		t.Error("Initial and new DRS do not match")
	}
}

func TestProgramJSONPersistence(t *testing.T) {
	var initialProgram = Program{
		Tree: &tree.Tree{
			Operator: tree.Sum{},
			Branches: []*tree.Tree{
				&tree.Tree{Operator: tree.Constant{42}},
				&tree.Tree{Operator: tree.Variable{1}},
			},
		},
		DRS: &DynamicRangeSelection{
			cutPoints: []float64{0, 1, 2},
			rangeMap:  map[float64]float64{0: -1, 1: 1, 2: -1},
		},
		Task: Task{
			Metric: metrics.MeanAbsoluteError{},
		},
	}

	const path = "test_program_json_persistence.json"

	// Persist the Program to the disk
	var err = SaveProgramToJSON(initialProgram, path)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Load the Program from the disk
	newProgram, err := LoadProgramFromJSON(path)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Compare the new Program with the initial Program
	var check func(n1, n2 *tree.Tree)
	check = func(n1, n2 *tree.Tree) {
		if n1.Operator.String() != n2.Operator.String() {
			t.Errorf("Operator mismatch: %s != %s", n1.Operator.String(), n2.Operator.String())
		}
		for i := range n1.Branches {
			check(n1.Branches[i], n2.Branches[i])
		}
	}
	check(newProgram.Tree, initialProgram.Tree)

	// Compare the DRSs
	if !reflect.DeepEqual(*newProgram.DRS, *initialProgram.DRS) {
		t.Error("Initial and new DRS do not match")
	}

	// Delete the JSON file
	os.Remove(path)
}
