package xgp

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// A serialNode can be serialized and holds information that can be used to
// initialize a Node.
type serialNode struct {
	OpType   string       `json:"op_type"`
	OpValue  string       `json:"op_value"`
	Children []serialNode `json:"children"`
}

// serializeNode recursively transforms a *Node into a serialNode.
func serializeNode(node *Node) (serialNode, error) {
	var serial = serialNode{
		Children: make([]serialNode, node.NBranches()),
	}
	switch node.Operator.(type) {
	case Constant:
		serial.OpType = "constant"
		serial.OpValue = strconv.FormatFloat(node.Operator.(Constant).Value, 'f', -1, 64)
	case Variable:
		serial.OpType = "variable"
		serial.OpValue = strconv.Itoa(node.Operator.(Variable).Index)
	default:
		serial.OpType = "function"
		serial.OpValue = node.Operator.String()
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
	switch serial.OpType {
	case "constant":
		var val, err = strconv.ParseFloat(serial.OpValue, 64)
		if err != nil {
			return nil, err
		}
		node.Operator = Constant{val}
	case "variable":
		var idx, err = strconv.Atoi(serial.OpValue)
		if err != nil {
			return nil, err
		}
		node.Operator = Variable{idx}
	default:
		var function, err = GetFunction(serial.OpValue)
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

// A serialDRS can be serialized and holds information that can be used to
// initialize a Node.
type serialDRS struct {
	CutPoints []float64         `json:"cut_points"`
	RangeMap  map[string]string `json:"range_map"`
}

// serializeDRS transforms a *DynamicRangeSelection into a serialDRS.
func serializeDRS(drs *DynamicRangeSelection) (serialDRS, error) {
	var serial = serialDRS{
		CutPoints: drs.cutPoints,
		RangeMap:  make(map[string]string),
	}
	for k, v := range drs.rangeMap {
		var (
			ks = strconv.FormatFloat(k, 'f', -1, 64)
			vs = strconv.FormatFloat(v, 'f', -1, 64)
		)
		serial.RangeMap[ks] = vs
	}
	return serial, nil
}

// parseSerialDRS recursively transforms a serialDRS into a *DynamicRangeSelection.
func parseSerialDRS(serial serialDRS) (*DynamicRangeSelection, error) {
	var drs = &DynamicRangeSelection{
		cutPoints: serial.CutPoints,
		rangeMap:  make(map[float64]float64),
	}
	for k, v := range serial.RangeMap {
		kf, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, err
		}
		vf, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		drs.rangeMap[kf] = vf
	}
	return drs, nil
}

// MarshalJSON serializes a *DynamicRangeSelection into JSON bytes. A serialDRS
// is used as an intermediary.
func (drs *DynamicRangeSelection) MarshalJSON() ([]byte, error) {
	var serial, err = serializeDRS(drs)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serial)
}

// UnmarshalJSON parses JSON bytes into a *DynamicRangeSelection. A serialDRS is
// used as an intermediary.
func (drs *DynamicRangeSelection) UnmarshalJSON(bytes []byte) error {
	var serial serialDRS
	if err := json.Unmarshal(bytes, &serial); err != nil {
		return err
	}
	var parsedDRS, err = parseSerialDRS(serial)
	if err != nil {
		return err
	}
	*drs = *parsedDRS
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

// ParseEquation parses an equation and returns a Program.
func ParseEquation(eq string) Program {
	return Program{}
}
