package op

import (
	"strings"
)

// An Operator takes float64s as input and outputs a float64.
type Operator interface {
	Eval(X [][]float64) []float64
	Arity() int
	String() string
}

// ParseStringFuncs parses a string into a slice of Operators.
func ParseStringFuncs(str string) ([]Operator, error) {
	var (
		strs  = strings.Split(str, ",")
		funcs = make([]Operator, len(strs))
	)
	for i, s := range strs {
		var f, err = ParseFuncName(s)
		if err != nil {
			return nil, err
		}
		funcs[i] = f
	}
	return funcs, nil
}
