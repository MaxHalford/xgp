package op

import (
	"fmt"
	"math"
)

// The Cos operator.
type Cos struct {
	Op Operator
}

// Eval computes the cosine of each value.
func (cos Cos) Eval(X [][]float64) []float64 {
	x := cos.Op.Eval(X)
	for i, xi := range x {
		x[i] = math.Cos(xi)
	}
	return x
}

// Arity of Cos is 1.
func (cos Cos) Arity() uint {
	return 1
}

// Operand returns Cos's operand or nil.
func (cos Cos) Operand(i uint) Operator {
	if i == 0 {
		return cos.Op
	}
	return nil
}

// SetOperand replaces Cos's operand if i is equal to 0.
func (cos Cos) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		cos.Op = op
	}
	return cos
}

// Simplify Cos.
func (cos Cos) Simplify() Operator {
	cos.Op = cos.Op.Simplify()
	return cos
}

// Diff computes the following derivative: cos(u)' = u' * -sin(u)
func (cos Cos) Diff(i uint) Operator {
	return Mul{cos.Op.Diff(i), Neg{Sin{cos.Op}}}
}

// Name of Cos is "cos".
func (cos Cos) Name() string {
	return "cos"
}

// String formatting.
func (cos Cos) String() string {
	return fmt.Sprintf("cos(%s)", cos.Op)
}
