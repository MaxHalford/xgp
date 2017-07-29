package xgp

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"
)

// A SerialNode can be serialized. The SerialNode's information can be used to
// reconstruct a Node.
type SerialNode struct {
	OperatorType  string       `json:"operator_type"`
	FunctionName  string       `json:"function_name"`
	ConstantValue float64      `json:"constant_value"`
	VariableIndex int          `json:"variable_index"`
	Children      []SerialNode `json:"children"`
}

// SerializeNode transforms a *Node into a SerialNode.
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

// ParseSerialNode transforms a SerialNode into a *Node.
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
		node.Operator = Sum{}
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

// UnmarshalJSON serializes a *Node into JSON bytes. A SerialNode is used as an
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

// SaveNodeToJSON saves a Node into a JSON file.
func SaveNodeToJSON(node *Node, path string) error {
	var bytes, err = json.Marshal(node)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// SaveNodeToJSON loads a Node from a JSON file.
func LoadNodeFromJSON(path string) (*Node, error) {
	var bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var node *Node
	err = json.Unmarshal(bytes, &node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// GobEncode makes the *Node type implement gob.GobEncoder.
func (node *Node) GobEncode() ([]byte, error) {
	return node.MarshalJSON()
}

// GobEncode makes the *Node type implement gob.GobDecoder.
func (node *Node) GobDecode(bytes []byte) error {
	return node.UnmarshalJSON(bytes)
}

// SaveNodeToGob saves a Node into a gob file.
func SaveNodeToGob(node *Node, path string) error {
	var file, err = os.Create(path)
	if err != nil {
		return err
	}
	var encoder = gob.NewEncoder(file)
	err = encoder.Encode(node)
	if err != nil {
		return err
	}
	return nil
}

// LoadNodeFromGob loads a Node from a gob file.
func LoadNodeFromGob(path string) (*Node, error) {
	var file, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	var (
		decoder = gob.NewDecoder(file)
		node    *Node
	)
	err = decoder.Decode(&node)
	if err != nil {
		return nil, err
	}
	return node, nil
}
