package op

import (
	"fmt"
	"math"
)

// The Abs operator.
type Abs struct {
	Op Operator
}

// Eval computes the absolute value of each value.
func (abs Abs) Eval(X [][]float64) []float64 {
	x := abs.Op.Eval(X)
	for i, xi := range x {
		x[i] = math.Abs(xi)
	}
	return x
}

// Arity of Abs is 1.
func (abs Abs) Arity() uint {
	return 1
}

// Operand returns Abs's operand or nil.
func (abs Abs) Operand(i uint) Operator {
	if i == 0 {
		return abs.Op
	}
	return nil
}

// SetOperand replaces Abs's operand if i is equal to 0.
func (abs Abs) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		abs.Op = op
	}
	return abs
}

// Simplify Abs.
func (abs Abs) Simplify() Operator {
	abs.Op = abs.Op.Simplify()
	switch operand := abs.Op.(type) {
	case Abs:
		// ||x|| = |x|
		return operand
	case Neg:
		// |-x| = x
		return operand.Op
	case Const:
		// |a| = b
		return Const{math.Abs(operand.Value)}
	}
	return abs
}

// Diff compute the following derivative: |u|' = uu' / |u|.
func (abs Abs) Diff(i uint) Operator {
	return Div{Mul{abs.Op, abs.Op.Diff(i)}, abs}
}

// Name of Abs is "abs".
func (abs Abs) Name() string {
	return "abs"
}

// String formatting.
func (abs Abs) String() string {
	return fmt.Sprintf("|%s|", parenthesize(abs.Op))
}
