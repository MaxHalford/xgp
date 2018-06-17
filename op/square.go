package op

import (
	"fmt"
	"math"
)

// The Square operator.
type Square struct {
	Op Operator
}

// Eval computes the square of each value.
func (square Square) Eval(X [][]float64) []float64 {
	x := square.Op.Eval(X)
	for i, xi := range x {
		x[i] = math.Pow(xi, 2)
	}
	return x
}

// Arity of Square is 1.
func (square Square) Arity() uint {
	return 1
}

// Operand returns Square's operand or nil.
func (square Square) Operand(i uint) Operator {
	if i == 0 {
		return square.Op
	}
	return nil
}

// SetOperand replaces Square's operand if i is equal to 0.
func (square Square) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		square.Op = op
	}
	return square
}

// Simplify Square.
func (square Square) Simplify() Operator {
	square.Op = square.Op.Simplify()
	switch op := square.Op.(type) {
	case Const:
		return Const{math.Pow(op.Value, 2)}
	case Neg:
		return Square{op.Op}.Simplify()
	}
	return square
}

// Diff computes the following derivative: (u²)' = 2u'u
func (square Square) Diff(i uint) Operator {
	return Mul{Const{2}, Mul{square.Op.Diff(i), square.Op}}
}

// Name of Square is "square".
func (square Square) Name() string {
	return "square"
}

// String formatting.
func (square Square) String() string {
	if square.Op.Arity() < 2 {
		return fmt.Sprintf("%s²", square.Op)
	}
	return fmt.Sprintf("%s²", parenthesize(square.Op))
}
