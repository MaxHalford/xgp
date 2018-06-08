package op

import (
	"fmt"

	"github.com/gonum/floats"
)

// The Add operator.
type Add struct {
	Left, Right Operator
}

// Eval sums aligned values.
func (add Add) Eval(X [][]float64) []float64 {
	x := add.Left.Eval(X)
	floats.Add(x, add.Right.Eval(X))
	return x
}

// Arity of Add is 2.
func (add Add) Arity() uint {
	return 2
}

// Operand returns one of Add's operands, or nil.
func (add Add) Operand(i uint) Operator {
	switch i {
	case 0:
		return add.Left
	case 1:
		return add.Right
	default:
		return nil
	}
}

// SetOperand replaces one of Add's operands if i is equal to 0 or 1.
func (add Add) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		add.Left = op
	} else if i == 1 {
		add.Right = op
	}
	return add
}

// The Add operator is symmetric so we have to check u+v as well as v+u.
func (add Add) simplify(left, right Operator) (Operator, bool) {
	switch left := left.(type) {
	case Const:
		switch left.Value {
		case 0:
			return right, true
		default:
			break
		}
	default:
		break
	}
	return add, false
}

// Simplify Add.
func (add Add) Simplify() Operator {

	// Simplify branches
	add.Left = add.Left.Simplify()
	add.Right = add.Right.Simplify()

	// Try to simplify u+v
	simpl, ok := add.simplify(add.Left, add.Right)
	if ok {
		return simpl.Simplify()
	}

	// Try to simplify v+u
	simpl, ok = add.simplify(add.Right, add.Left)
	if ok {
		return simpl.Simplify()
	}

	return simpl
}

// Diff computes the following derivative: (u + v)' = u' + v'
func (add Add) Diff(i uint) Operator {
	return Add{add.Left.Diff(i), add.Right.Diff(i)}
}

// Name of Add is "add".
func (add Add) Name() string {
	return "add"
}

// String formatting.
func (add Add) String() string {
	return fmt.Sprintf("%s+%s", parenthesize(add.Left), parenthesize(add.Right))
}
