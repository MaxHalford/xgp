package op

import (
	"fmt"

	"github.com/gonum/floats"
)

// The Mul operator.
type Mul struct {
	Left, Right Operator
}

// Eval multiplies aligned values.
func (mul Mul) Eval(X [][]float64) []float64 {
	x := mul.Left.Eval(X)
	floats.Mul(x, mul.Right.Eval(X))
	return x
}

// Arity of Mul is 2.
func (mul Mul) Arity() uint {
	return 2
}

// Operand returns one of Mul's operands, or nil.
func (mul Mul) Operand(i uint) Operator {
	switch i {
	case 0:
		return mul.Left
	case 1:
		return mul.Right
	default:
		return nil
	}
}

// SetOperand replaces one of Mul's operands if i is equal to 0 or 1.
func (mul Mul) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		mul.Left = op
	} else if i == 1 {
		mul.Right = op
	}
	return mul
}

// The Mul operator is symmetric so we have to check u*v as well as v*u.
func (mul Mul) simplify(left, right Operator) (Operator, bool) {
	switch left := left.(type) {
	case Const:
		switch left.Value {
		// 0 * a = 0
		case 0:
			return Const{0}, true
		// 1 * a = a
		case 1:
			return right, true
		// -1 * a = -a
		case -1:
			return Neg{right}, true
		}
		switch right := right.(type) {
		// a * b = c
		case Const:
			return Const{left.Value * right.Value}, true
		}
	}
	return mul, false
}

// Simplify Mul.
func (mul Mul) Simplify() Operator {

	// Simplify branches
	mul.Left = mul.Left.Simplify()
	mul.Right = mul.Right.Simplify()

	// Try to simplify left/right
	simpl, ok := mul.simplify(mul.Left, mul.Right)
	if ok {
		return simpl.Simplify()
	}

	// Try to simplify right/left
	simpl, ok = mul.simplify(mul.Right, mul.Left)
	if ok {
		return simpl.Simplify()
	}

	return simpl
}

// Diff computes the following derivative: (uv)' = u'v + uv'
func (mul Mul) Diff(i uint) Operator {
	return Add{Mul{mul.Left.Diff(i), mul.Right}, Mul{mul.Left, mul.Right.Diff(i)}}
}

// Name of Mul is "mul".
func (mul Mul) Name() string {
	return "mul"
}

// String formatting.
func (mul Mul) String() string {
	return fmt.Sprintf("%s*%s", parenthesize(mul.Left), parenthesize(mul.Right))
}
