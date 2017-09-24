package xgp

import (
	"encoding/json"
	"io/ioutil"
)

// A serialNode can be serialized and holds information that can be used to
// initialize a Node.
type serialNode struct {
	OperatorType  string       `json:"operator_type"`
	FunctionName  string       `json:"function_name"`
	ConstantValue float64      `json:"constant_value"`
	VariableIndex int          `json:"variable_index"`
	Children      []serialNode `json:"children"`
}

// serializeNode recursively transforms a *Node into a serialNode.
func serializeNode(node *Node) (serialNode, error) {
	var serial = serialNode{
		Children: make([]serialNode, node.NBranches()),
	}
	switch node.Operator.(type) {
	case Constant:
		serial.OperatorType = "constant"
		serial.ConstantValue = node.Operator.(Constant).Value
	case Variable:
		serial.OperatorType = "variable"
		serial.VariableIndex = node.Operator.(Variable).Index
	default:
		serial.OperatorType = "function"
		serial.FunctionName = node.Operator.String()
	}
	for i, child := range node.Children {
		var serialChild, err = serializeNode(child)
		if err != nil {
			return serial, err
		}
		serial.Children[i] = serialChild
	}
	return serial, nil
}

// parseSerialNode recursively transforms a serialNode into a *Node.
func parseSerialNode(serial serialNode) (*Node, error) {
	var node = &Node{
		Children: make([]*Node, len(serial.Children)),
	}
	switch serial.OperatorType {
	case "constant":
		node.Operator = Constant{serial.ConstantValue}
	case "variable":
		node.Operator = Variable{serial.VariableIndex}
	default:
		var function, err = GetFunction(serial.FunctionName)
		if err != nil {
			return nil, err
		}
		node.Operator = function
	}
	for i, child := range serial.Children {
		var nodeChild, err = parseSerialNode(child)
		if err != nil {
			return node, err
		}
		node.Children[i] = nodeChild
	}
	return node, nil
}

// MarshalJSON serializes a *Node into JSON bytes. A serialNode is used as an
// intermediary.
func (node *Node) MarshalJSON() ([]byte, error) {
	var serial, err = serializeNode(node)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *Node. A serialNode is used as an
// intermediary.
func (node *Node) UnmarshalJSON(bytes []byte) error {
	var serial serialNode
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedNode, err = parseSerialNode(serial)
	if err != nil {
		return err
	}
	*node = *parsedNode
	return nil
}

// A serialProgram can be serialized and holds information that can be used to
// initialize a Program.
type serialProgram struct {
	Root          serialNode `json:"root"`
	TransformName string     `json:"transform"`
}

// serializeProgram transforms a serialProgram into a Program.
func serializeProgram(prog Program) (serialProgram, error) {
	var root, err = serializeNode(prog.Root)
	if err != nil {
		return serialProgram{}, err
	}
	return serialProgram{
		Root:          root,
		TransformName: prog.Estimator.Transform.String(),
	}, nil
}

// parseserialProgram recursively transforms a serialProgram into a Program.
func parseserialProgram(serial serialProgram) (Program, error) {
	var root, err = parseSerialNode(serial.Root)
	if err != nil {
		return Program{}, err
	}
	transform, err := GetTransform(serial.TransformName)
	if err != nil {
		return Program{}, err
	}
	return Program{
		Root:      root,
		Estimator: &Estimator{Transform: transform},
	}, nil
}

// MarshalJSON serializes a Program into JSON. A serialProgram is used as an
// intermediary.
func (prog *Program) MarshalJSON() ([]byte, error) {
	var serial, err = serializeProgram(*prog)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON into a Program. A serialProgram is used as an
// intermediary.
func (prog *Program) UnmarshalJSON(bytes []byte) error {
	var serial serialProgram
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedProg, err = parseserialProgram(serial)
	if err != nil {
		return err
	}
	*prog = parsedProg
	return nil
}

// SaveProgramToJSON saves a Program to a JSON file.
func SaveProgramToJSON(program Program, path string) error {
	var bytes, err = json.Marshal(&program)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadProgramFromJSON loads a Program from a JSON file.
func LoadProgramFromJSON(path string) (Program, error) {
	var (
		program    Program
		bytes, err = ioutil.ReadFile(path)
	)
	if err != nil {
		return program, err
	}
	err = json.Unmarshal(bytes, &program)
	if err != nil {
		return program, err
	}
	return program, nil
}
