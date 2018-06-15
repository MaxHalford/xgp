package op

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type serialOperator struct {
	Type     string           `json:"type"`
	Value    string           `json:"value"`
	Operands []serialOperator `json:"operands"`
}

func serializeOp(op Operator) serialOperator {
	var (
		arity  = op.Arity()
		serial = serialOperator{Operands: make([]serialOperator, arity)}
	)
	for i := uint(0); i < op.Arity(); i++ {
		serial.Operands[i] = serializeOp(op.Operand(i))
	}
	switch op := op.(type) {
	case Const:
		serial.Type = "const"
		serial.Value = strconv.FormatFloat(op.Value, 'f', -1, 64)
	case Var:
		serial.Type = "var"
		serial.Value = strconv.Itoa(int(op.Index))
	default:
		serial.Type = "func"
		serial.Value = op.Name()
	}
	return serial
}

func parseOp(serial serialOperator) (Operator, error) {
	var op Operator
	switch serial.Type {
	case "const":
		val, err := strconv.ParseFloat(serial.Value, 64)
		if err != nil {
			return nil, err
		}
		op = Const{val}
	case "var":
		idx, err := strconv.Atoi(serial.Value)
		if err != nil {
			return nil, err
		}
		op = Var{uint(idx)}
	default:
		function, err := ParseFunc(serial.Value)
		if err != nil {
			return nil, err
		}
		op = function
	}
	// Set the operands; this is where the recursion happens
	for i, serialOperand := range serial.Operands {
		operand, err := parseOp(serialOperand)
		if err != nil {
			return nil, err
		}
		op = op.SetOperand(uint(i), operand)
	}
	return op, nil
}

// MarshalJSON serializes an Operator into JSON.
func MarshalJSON(op Operator) ([]byte, error) {
	return json.Marshal(serializeOp(op))
}

// UnmarshalJSON parses JSON into an Operator.
func UnmarshalJSON(raw []byte) (Operator, error) {
	var serial = serialOperator{}
	if err := json.Unmarshal(raw, &serial); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return parseOp(serial)
}
