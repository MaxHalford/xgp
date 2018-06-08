package op

import (
	"strconv"

	"gonum.org/v1/gonum/floats"
)

// The Const operator always returns the same value.
type Const struct {
	Value float64
}

// Eval returns the same value.
func (c Const) Eval(X [][]float64) []float64 {
	var x = make([]float64, len(X[0]))
	floats.AddConst(c.Value, x)
	return x
}

// Arity of Const is 0.
func (c Const) Arity() uint {
	return 0
}

// Operand returns nil.
func (c Const) Operand(i uint) Operator {
	return nil
}

// SetOperand returns the Const without any modifications.
func (c Const) SetOperand(i uint, op Operator) Operator {
	return c
}

// Simplify returns the Const without any modifications.
func (c Const) Simplify() Operator {
	return c
}

// Diff computes the following derivative: c' = 0
func (c Const) Diff(i uint) Operator {
	return Const{0}
}

// Name of Const is it's value.
func (c Const) Name() string {
	return strconv.FormatFloat(c.Value, 'f', -1, 64)
}

// String formatting.
func (c Const) String() string {
	return c.Name()
}
