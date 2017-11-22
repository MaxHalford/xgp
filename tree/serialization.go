package tree

import (
	"encoding/json"
	"strconv"
)

// A serialtree can be serialized and holds information that can be used to
// initialize a tree.
type serialtree struct {
	OpType   string       `json:"op_type"`
	OpValue  string       `json:"op_value"`
	Branches []serialtree `json:"branches"`
}

// serializetree recursively transforms a *tree into a serialtree.
func serializetree(tree *Tree) (serialtree, error) {
	var serial = serialtree{
		Branches: make([]serialtree, len(tree.Branches)),
	}
	switch tree.Operator.(type) {
	case Constant:
		serial.OpType = "constant"
		serial.OpValue = strconv.FormatFloat(tree.Operator.(Constant).Value, 'f', -1, 64)
	case Variable:
		serial.OpType = "variable"
		serial.OpValue = strconv.Itoa(tree.Operator.(Variable).Index)
	default:
		serial.OpType = "function"
		serial.OpValue = tree.Operator.String()
	}
	for i, branches := range tree.Branches {
		var serialChild, err = serializetree(branches)
		if err != nil {
			return serial, err
		}
		serial.Branches[i] = serialChild
	}
	return serial, nil
}

// parseSerialtree recursively transforms a serialtree into a *tree.
func parseSerialtree(serial serialtree) (*Tree, error) {
	var tree = &Tree{
		Branches: make([]*Tree, len(serial.Branches)),
	}
	switch serial.OpType {
	case "constant":
		var val, err = strconv.ParseFloat(serial.OpValue, 64)
		if err != nil {
			return nil, err
		}
		tree.Operator = Constant{val}
	case "variable":
		var idx, err = strconv.Atoi(serial.OpValue)
		if err != nil {
			return nil, err
		}
		tree.Operator = Variable{idx}
	default:
		var function, err = parseFuncName(serial.OpValue)
		if err != nil {
			return nil, err
		}
		tree.Operator = function
	}
	for i, branches := range serial.Branches {
		var treeChild, err = parseSerialtree(branches)
		if err != nil {
			return tree, err
		}
		tree.Branches[i] = treeChild
	}
	return tree, nil
}

// MarshalJSON serializes a *tree into JSON bytes. A serialtree is used as an
// intermediary.
func (tree *Tree) MarshalJSON() ([]byte, error) {
	var serial, err = serializetree(tree)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *tree. A serialtree is used as an
// intermediary.
func (tree *Tree) UnmarshalJSON(bytes []byte) error {
	var serial serialtree
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedtree, err = parseSerialtree(serial)
	if err != nil {
		return err
	}
	*tree = *parsedtree
	return nil
}
