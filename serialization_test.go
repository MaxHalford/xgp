package xgp

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestNodeJSONEncodeDecode(t *testing.T) {
	var initialNode = &Node{
		Operator: Sum{},
		Children: []*Node{
			&Node{Operator: Constant{42}},
			&Node{Operator: Variable{1}},
		},
	}

	// Serialize the initial Node
	var bytes, err = json.Marshal(initialNode)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Parse the bytes into a new Node
	var newNode *Node
	err = json.Unmarshal(bytes, &newNode)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Compare the new Node with the initial Node
	var check func(n1, n2 *Node)
	check = func(n1, n2 *Node) {
		if n1.Operator.String() != n2.Operator.String() {
			t.Errorf("Operator mismatch: %s != %s", n1.Operator.String(), n2.Operator.String())
		}
		for i := range n1.Children {
			check(n1.Children[i], n2.Children[i])
		}
	}
	check(newNode, initialNode)
}

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

	// Parse the bytes into a new Node
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
		Root: &Node{
			Operator: Sum{},
			Children: []*Node{
				&Node{Operator: Constant{42}},
				&Node{Operator: Variable{1}},
			},
		},
		DRS: &DynamicRangeSelection{
			cutPoints: []float64{0, 1, 2},
			rangeMap:  map[float64]float64{0: -1, 1: 1, 2: -1},
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
	var check func(n1, n2 *Node)
	check = func(n1, n2 *Node) {
		if n1.Operator.String() != n2.Operator.String() {
			t.Errorf("Operator mismatch: %s != %s", n1.Operator.String(), n2.Operator.String())
		}
		for i := range n1.Children {
			check(n1.Children[i], n2.Children[i])
		}
	}
	check(newProgram.Root, initialProgram.Root)

	// Compare the DRSs
	if !reflect.DeepEqual(*newProgram.DRS, *initialProgram.DRS) {
		t.Error("Initial and new DRS do not match")
	}

	// Delete the JSON file
	os.Remove(path)
}
