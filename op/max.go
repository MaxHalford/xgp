package op

import "fmt"

// The Max operator.
type Max struct {
	Left, Right Operator
}

// Eval takes the maximum along aligned values.
func (max Max) Eval(X [][]float64) []float64 {
	x := max.Left.Eval(X)
	y := max.Right.Eval(X)
	for i, yi := range y {
		if yi > x[i] {
			x[i] = yi
		}
	}
	return x
}

// Arity of Max is 2.
func (max Max) Arity() uint {
	return 2
}

// Operand returns one of Max's operands, or nil.
func (max Max) Operand(i uint) Operator {
	switch i {
	case 0:
		return max.Left
	case 1:
		return max.Right
	default:
		return nil
	}
}

// SetOperand replaces one of Max's operands if i is equal to 0 or 1.
func (max Max) SetOperand(i uint, op Operator) Operator {
	switch i {
	case 0:
		return Max{op, max.Right}
	case 1:
		return Max{max.Left, op}
	default:
		return max
	}
}

// Simplify Max.
func (max Max) Simplify() Operator {
	max.Left = max.Left.Simplify()
	max.Right = max.Right.Simplify()
	// max(a, b) = c
	left, ok := max.Left.(Const)
	if !ok {
		return max
	}
	right, ok := max.Right.(Const)
	if !ok {
		return max
	}
	if left.Value > right.Value {
		return left
	}
	return right
}

// Diff computes the following derivative: max(u, v)' = ((u + v + |u - v|) / 2)'
func (max Max) Diff(i uint) Operator {
	return Div{
		Add{Add{max.Left, max.Right}, Abs{Sub{max.Left, max.Right}}},
		Const{2},
	}.Diff(i)
}

// Name of Max is "max".
func (max Max) Name() string {
	return "max"
}

// String formatting.
func (max Max) String() string {
	return fmt.Sprintf("max(%s, %s)", max.Left, max.Right)
}
