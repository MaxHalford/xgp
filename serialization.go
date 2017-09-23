package xgp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// FUNCTIONS maps string representations of Functions to their respective
// Function for serialization purposes.
var FUNCTIONS = map[string]Operator{
	Cos{}.String():        Cos{},
	Sin{}.String():        Sin{},
	Log{}.String():        Log{},
	Exp{}.String():        Exp{},
	Max{}.String():        Max{},
	Min{}.String():        Min{},
	Sum{}.String():        Sum{},
	Difference{}.String(): Difference{},
	Division{}.String():   Division{},
	Product{}.String():    Product{},
	Power{}.String():      Power{},
}

// TRANSFORMS maps string representations of Transforms to their respective
// Transform for serialization purposes.
var TRANSFORMS = map[string]Transform{
	Identity{}.String(): Identity{},
	Binary{}.String():   Binary{},
	Sigmoid{}.String():  Sigmoid{},
}

// A SerialNode can be serialized and holds information that can be used to
// initialize a Node.
type SerialNode struct {
	OperatorType  string       `json:"operator_type"`
	FunctionName  string       `json:"function_name"`
	ConstantValue float64      `json:"constant_value"`
	VariableIndex int          `json:"variable_index"`
	Children      []SerialNode `json:"children"`
}

// SerializeNode recursively transforms a *Node into a SerialNode.
func SerializeNode(node *Node) (SerialNode, error) {
	var serial = SerialNode{
		Children: make([]SerialNode, node.NBranches()),
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
		var serialChild, err = SerializeNode(child)
		if err != nil {
			return serial, err
		}
		serial.Children[i] = serialChild
	}
	return serial, nil
}

// ParseSerialNode recursively transforms a SerialNode into a *Node.
func ParseSerialNode(serial SerialNode) (*Node, error) {
	var node = &Node{
		Children: make([]*Node, len(serial.Children)),
	}
	switch serial.OperatorType {
	case "constant":
		node.Operator = Constant{serial.ConstantValue}
	case "variable":
		node.Operator = Variable{serial.VariableIndex}
	default:
		var function, ok = FUNCTIONS[serial.FunctionName]
		if !ok {
			return nil, fmt.Errorf("Unknown function name '%s'", serial.FunctionName)
		}
		node.Operator = function
	}
	for i, child := range serial.Children {
		var nodeChild, err = ParseSerialNode(child)
		if err != nil {
			return node, err
		}
		node.Children[i] = nodeChild
	}
	return node, nil
}

// MarshalJSON serializes a *Node into JSON bytes. A SerialNode is used as an
// intermediary.
func (node *Node) MarshalJSON() ([]byte, error) {
	var serial, err = SerializeNode(node)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *Node. A SerialNode is used as an
// intermediary.
func (node *Node) UnmarshalJSON(bytes []byte) error {
	var serial SerialNode
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedNode, err = ParseSerialNode(serial)
	if err != nil {
		return err
	}
	*node = *parsedNode
	return nil
}

// A SerialProgram can be serialized and holds information that can be used to
// initialize a Program.
type SerialProgram struct {
	Root          SerialNode `json:"root"`
	TransformName string     `json:"transform"`
}

// SerializeProgram transforms a SerialProgram into a Program.
func SerializeProgram(prog Program) (SerialProgram, error) {
	var root, err = SerializeNode(prog.Root)
	if err != nil {
		return SerialProgram{}, err
	}
	return SerialProgram{
		Root:          root,
		TransformName: prog.Estimator.Transform.String(),
	}, nil
}

// ParseSerialProgram recursively transforms a SerialProgram into a Program.
func ParseSerialProgram(serial SerialProgram) (Program, error) {
	var root, err = ParseSerialNode(serial.Root)
	if err != nil {
		return Program{}, err
	}
	var transform, ok = TRANSFORMS[serial.TransformName]
	if !ok {
		return Program{}, fmt.Errorf("Unknown transform name '%s'", serial.TransformName)
	}
	return Program{
		Root:      root,
		Estimator: &Estimator{Transform: transform},
	}, nil
}

// MarshalJSON serializes a Program into JSON bytes. A SerialProgram is used as
// an intermediary.
func (prog *Program) MarshalJSON() ([]byte, error) {
	var serial, err = SerializeProgram(*prog)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a Program. A SerialProgram is used as an
// intermediary.
func (prog *Program) UnmarshalJSON(bytes []byte) error {
	var serial SerialProgram
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedProg, err = ParseSerialProgram(serial)
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

// LoadProgramFromJSON loadataset a Program from a JSON file.
func LoadProgramFromJSON(path string) (Program, error) {
	var bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return Program{}, err
	}
	var program Program
	err = json.Unmarshal(bytes, &program)
	if err != nil {
		return Program{}, err
	}
	return program, nil
}
