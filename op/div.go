package op

import "fmt"

// The Div operator.
type Div struct {
	Left, Right Operator
}

// Eval divides aligned values.
func (div Div) Eval(X [][]float64) []float64 {
	x := div.Left.Eval(X)
	y := div.Right.Eval(X)
	for i, yi := range y {
		x[i] = safeDiv(x[i], yi)
	}
	return x
}

// Arity of Div is 2.
func (div Div) Arity() uint {
	return 2
}

// Operand returns one of Div's operands, or nil.
func (div Div) Operand(i uint) Operator {
	switch i {
	case 0:
		return div.Left
	case 1:
		return div.Right
	default:
		return nil
	}
}

// SetOperand replaces one of Div's operands if i is equal to 0 or 1.
func (div Div) SetOperand(i uint, op Operator) Operator {
	switch i {
	case 0:
		return Div{op, div.Right}
	case 1:
		return Div{div.Left, op}
	default:
		return div
	}
}

// Simplify Div.
func (div Div) Simplify() Operator {

	// Simplify branches
	div.Left = div.Left.Simplify()
	div.Right = div.Right.Simplify()

	switch right := div.Right.(type) {
	case Const:
		switch right.Value {
		// x / 0 = 1
		case 0:
			return Const{1}
		// x / 1 = x
		case 1:
			return div.Left
		// x / -1 = -x
		case -1:
			return Neg{div.Left}.Simplify()
		}
		switch left := div.Left.(type) {
		// a / b = c
		case Const:
			return Const{safeDiv(left.Value, right.Value)}
		}
	case Var:
		if left, ok := div.Left.(Var); ok {
			// x / x = 1
			if left.Index == right.Index {
				return Const{1}
			}
		}
	}
	switch left := div.Left.(type) {
	case Const:
		switch left.Value {
		// 0 / x = 0
		case 0:
			return Const{0}
		// 1 / x = inv(x)
		case 1:
			return Inv{div.Right}.Simplify()
		// -1 / x = inv(-x)
		case -1:
			return Inv{Neg{div.Right}}.Simplify()
		}
	}
	return div
}

// Diff computes the following derivative: (u / v)' = (u'v - uv') / v²
func (div Div) Diff(i uint) Operator {
	return Div{
		Sub{
			Mul{div.Left.Diff(i), div.Right},
			Mul{div.Left, div.Right.Diff(i)}},
		div.Right,
	}
}

// Name of Div is "div".
func (div Div) Name() string {
	return "div"
}

// String formatting.
func (div Div) String() string {
	return fmt.Sprintf("%s/%s", parenthesize(div.Left), parenthesize(div.Right))
}
