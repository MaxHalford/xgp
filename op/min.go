package op

import "fmt"

// The Min operator.
type Min struct {
	Left, Right Operator
}

// Eval takes the minimum along aligned values.
func (min Min) Eval(X [][]float64) []float64 {
	x := min.Left.Eval(X)
	y := min.Right.Eval(X)
	for i, yi := range y {
		if yi < x[i] {
			x[i] = yi
		}
	}
	return x
}

// Arity of Min is 2.
func (min Min) Arity() uint {
	return 2
}

// Operand returns one of Min's operands, or nil.
func (min Min) Operand(i uint) Operator {
	switch i {
	case 0:
		return min.Left
	case 1:
		return min.Right
	default:
		return nil
	}
}

// SetOperand replaces one of Min's operands if i is equal to 0 or 1.
func (min Min) SetOperand(i uint, op Operator) Operator {
	switch i {
	case 0:
		return Min{op, min.Right}
	case 1:
		return Min{min.Left, op}
	default:
		return min
	}
}

// Simplify Min.
func (min Min) Simplify() Operator {
	min.Left = min.Left.Simplify()
	min.Right = min.Right.Simplify()
	// min(a, b) = c
	left, ok := min.Left.(Const)
	if !ok {
		return min
	}
	right, ok := min.Right.(Const)
	if !ok {
		return min
	}
	if left.Value > right.Value {
		return right
	}
	return left
}

// Diff computes the following derivative: min(u, v)' = ((u + v - |u - v|) / 2)'
func (min Min) Diff(i uint) Operator {
	return Div{
		Sub{Add{min.Left, min.Right}, Abs{Sub{min.Left, min.Right}}},
		Const{2},
	}.Diff(i)
}

// Name of Min is "min".
func (min Min) Name() string {
	return "min"
}

// String formatting.
func (min Min) String() string {
	return fmt.Sprintf("min(%s, %s)", min.Left, min.Right)
}
