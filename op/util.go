package op

import (
	"fmt"
	"strings"
)

// ParseFunc returns a functional Operator from it's String representation.
func ParseFunc(name string) (Operator, error) {
	var f, ok = map[string]Operator{
		Cos{}.String(): Cos{},
		Sin{}.String(): Sin{},
		Exp{}.String(): Exp{},
		Max{}.String(): Max{},
		Min{}.String(): Min{},
		Sum{}.String(): Sum{},
		Sub{}.String(): Sub{},
		Div{}.String(): Div{},
		Mul{}.String(): Mul{},
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
