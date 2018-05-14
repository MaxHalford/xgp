package tree

import (
	"strconv"

	"github.com/MaxHalford/xgp/op"
)

// A serialTree can be serialized and holds information that can be used to
// initialize a tree.
type serialTree struct {
	OpType   string       `json:"op_type"`
	OpValue  string       `json:"op_value"`
	Branches []serialTree `json:"branches"`
}

// serializeTree recursively transforms a Tree into a serialTree.
func serializeTree(tree Tree) (serialTree, error) {
	var serial = serialTree{
		Branches: make([]serialTree, len(tree.branches)),
	}
	switch tree.op.(type) {
	case op.Constant:
		serial.OpType = "constant"
		serial.OpValue = strconv.FormatFloat(tree.op.(op.Constant).Value, 'f', -1, 64)
	case op.Variable:
		serial.OpType = "variable"
		serial.OpValue = strconv.Itoa(tree.op.(op.Variable).Index)
	default:
		serial.OpType = "function"
		serial.OpValue = tree.op.String()
	}
	for i, branch := range tree.branches {
		var serialBranch, err = serializeTree(*branch)
		if err != nil {
			return serial, err
		}
		serial.Branches[i] = serialBranch
	}
	return serial, nil
}

// parseSerialTree recursively transforms a serialTree into a *tree.
func parseSerialTree(serial serialTree) (Tree, error) {
	var tree = Tree{
		branches: make([]*Tree, len(serial.Branches)),
	}
	switch serial.OpType {
	case "constant":
		var val, err = strconv.ParseFloat(serial.OpValue, 64)
		if err != nil {
			return tree, err
		}
		tree.op = op.Constant{val}
	case "variable":
		var idx, err = strconv.Atoi(serial.OpValue)
		if err != nil {
			return tree, err
		}
		tree.op = op.Variable{idx}
	default:
		var function, err = op.ParseFuncName(serial.OpValue)
		if err != nil {
			return tree, err
		}
		tree.op = function
	}
	for i, branches := range serial.Branches {
		var treeChild, err = parseSerialTree(branches)
		if err != nil {
			return tree, err
		}
		tree.SetBranch(i, treeChild)
	}
	return tree, nil
}
