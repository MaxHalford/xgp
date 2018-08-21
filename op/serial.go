package op

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFunc parses a name and returns the corresponding Operator.
func ParseFunc(name string) (Operator, error) {
	var f, ok = map[string]Operator{
		If{}.Name():     If{},
		Abs{}.Name():    Abs{},
		Add{}.Name():    Add{},
		Cos{}.Name():    Cos{},
		Div{}.Name():    Div{},
		Inv{}.Name():    Inv{},
		Max{}.Name():    Max{},
		Min{}.Name():    Min{},
		Mul{}.Name():    Mul{},
		Neg{}.Name():    Neg{},
		Sin{}.Name():    Sin{},
		Square{}.Name(): Square{},
		Sub{}.Name():    Sub{},
	}[name]
	if !ok {
		return nil, fmt.Errorf("Unknown function name '%s'", name)
	}
	return f, nil
}

// ParseFuncs parses a string into a slice of Operators.
func ParseFuncs(names, sep string) ([]Operator, error) {
	var funcs = make([]Operator, strings.Count(names, sep)+1)
	for i, name := range strings.Split(names, sep) {
		var f, err = ParseFunc(name)
		if err != nil {
			return nil, err
		}
		funcs[i] = f
	}
	return funcs, nil
}

type SerialOp struct {
	Type     string     `json:"type"`
	Value    string     `json:"value"`
	Operands []SerialOp `json:"operands"`
}

func SerializeOp(op Operator) SerialOp {
	var (
		arity  = op.Arity()
		serial = SerialOp{Operands: make([]SerialOp, arity)}
	)
	for i := uint(0); i < op.Arity(); i++ {
		serial.Operands[i] = SerializeOp(op.Operand(i))
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

func ParseOp(serial SerialOp) (Operator, error) {
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
	for i, serialOp := range serial.Operands {
		operand, err := ParseOp(serialOp)
		if err != nil {
			return nil, err
		}
		op = op.SetOperand(uint(i), operand)
	}
	return op, nil
}
