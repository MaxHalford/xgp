package op

import (
	"fmt"

	"github.com/gonum/floats"
)

// The Sub operator.
type Sub struct {
	Left, Right Operator
}

// Eval computes the difference between Left and Right.
func (sub Sub) Eval(X [][]float64) []float64 {
	x := sub.Left.Eval(X)
	floats.Sub(x, sub.Right.Eval(X))
	return x
}

// Arity of Sub is 2.
func (sub Sub) Arity() uint {
	return 2
}

// Operand returns one of Sub's operands, or nil.
func (sub Sub) Operand(i uint) Operator {
	switch i {
	case 0:
		return sub.Left
	case 1:
		return sub.Right
	default:
		return nil
	}
}

// SetOperand replaces one of Sub's operands if i is equal to 0 or 1.
func (sub Sub) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		sub.Left = op
	} else if i == 1 {
		sub.Right = op
	}
	return sub
}

// Simplify Sub.
func (sub Sub) Simplify() Operator {
	switch left := sub.Left.(type) {
	case Const:
		switch left.Value {
		// 0 - x = -x
		case 0:
			return Neg{}
		default:
			break
		}
		switch right := sub.Right.(type) {
		// a - b = c
		case Const:
			return Const{left.Value - right.Value}
		default:
			break
		}
	case Var:
		if right, ok := sub.Right.(Var); ok {
			// x - x = 1
			if left.Index == right.Index {
				return Const{0}
			}
		}
	default:
		break
	}
	switch right := sub.Right.(type) {
	case Const:
		// x - 0 = x
		if right.Value == 0 {
			return sub.Left
		}
	default:
		break
	}
	return sub
}

// Diff computes the following derivative: (u - v)' = u' - v'
func (sub Sub) Diff(i uint) Operator {
	return Sub{sub.Left.Diff(i), sub.Right.Diff(i)}
}

// Name of Sub is "sub".
func (sub Sub) Name() string {
	return "sub"
}

// String formatting.
func (sub Sub) String() string {
	return fmt.Sprintf("%s-%s", parenthesize(sub.Left), parenthesize(sub.Right))
}
