package op

import (
	"strconv"

	"gonum.org/v1/gonum/floats"
)

// A Constant holds a float64 value.
type Constant struct {
	Value float64
}

// ApplyRow Constant.
func (c Constant) ApplyRow(x []float64) float64 {
	return c.Value
}

// ApplyCols Constant.
func (c Constant) ApplyCols(X [][]float64) []float64 {
	var C = make([]float64, len(X[0]))
	floats.AddConst(c.Value, C)
	return C
}

// Arity of a Constant is 0 because it is a terminal operator.
func (c Constant) Arity() int {
	return 0
}

// String representation of a Constant.
func (c Constant) String() string {
	return strconv.FormatFloat(c.Value, 'f', -1, 64)
}
