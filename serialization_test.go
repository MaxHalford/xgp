package xgp

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
	"testing"
)

func TestJSONEncodeDecode(t *testing.T) {
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

func TestJSONPersistence(t *testing.T) {
	var initialNode = &Node{
		Operator: Sum{},
		Children: []*Node{
			&Node{Operator: Constant{42}},
			&Node{Operator: Variable{1}},
		},
	}

	const path = "node_test_json_persistence.json"

	// Persist the Node to the disk
	var err = SaveNodeToJSON(initialNode, path)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Load the Node from the disk
	newNode, err := LoadNodeFromJSON(path)
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

	// Delete the JSON file
	os.Remove(path)
}

func TestGobEncodeDecode(t *testing.T) {
	var initialNode = &Node{
		Operator: Sum{},
		Children: []*Node{
			&Node{Operator: Constant{42}},
			&Node{Operator: Variable{1}},
		},
	}

	// Initialize variables for gob
	var (
		buffer  bytes.Buffer
		encoder = gob.NewEncoder(&buffer)
		decoder = gob.NewDecoder(&buffer)
	)

	// Serialize the initial Node
	err := encoder.Encode(initialNode)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Parse the bytes into a new Node
	var newNode *Node
	err = decoder.Decode(&newNode)
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

func TestGobPersistence(t *testing.T) {
	var initialNode = &Node{
		Operator: Sum{},
		Children: []*Node{
			&Node{Operator: Constant{42}},
			&Node{Operator: Variable{1}},
		},
	}

	const path = "node_test_gob_persistence.gob"

	// Persist the Node to the disk
	var err = SaveNodeToGob(initialNode, path)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Load the Node from the disk
	newNode, err := LoadNodeFromGob(path)
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

	// Delete the JSON file
	os.Remove(path)
}
