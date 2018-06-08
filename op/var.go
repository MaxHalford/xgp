package op

import "fmt"

// The Var operator.
type Var struct {
	Index uint
}

// Eval returns the ith column of a matrix.
func (v Var) Eval(X [][]float64) []float64 {
	var x = make([]float64, len(X[v.Index]))
	copy(x, X[v.Index])
	return x
}

// Arity of Var is 0.
func (v Var) Arity() uint {
	return 0
}

// Operand returns nil.
func (v Var) Operand(i uint) Operator {
	return nil
}

// SetOperand returns the Const without any modifications.
func (v Var) SetOperand(i uint, op Operator) Operator {
	return v
}

// Simplify returns the Const without any modifications.
func (v Var) Simplify() Operator {
	return v
}

// Diff computes the following derivative: x' = 1. However if i not equal to
// the Var's Index then a 0 Const is returned.
func (v Var) Diff(i uint) Operator {
	if i == v.Index {
		return Const{1}
	}
	return Const{0}
}

// Name of Var is "xi" where i is the Var's Index.
func (v Var) Name() string {
	return fmt.Sprintf("x%d", v.Index)
}

// String formatting.
func (v Var) String() string {
	return v.Name()
}
