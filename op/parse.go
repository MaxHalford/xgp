package op

import (
	"fmt"
	"strings"
)

// ParseFunc parses a name and returns the corresponding Operator.
func ParseFunc(name string) (Operator, error) {
	var f, ok = map[string]Operator{
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
