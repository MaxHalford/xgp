package op

import (
	"fmt"
	"math"
)

// The Sin operator.
type Sin struct {
	Op Operator
}

// Eval computes the sine of each value.
func (sin Sin) Eval(X [][]float64) []float64 {
	x := sin.Op.Eval(X)
	for i, xi := range x {
		x[i] = math.Sin(xi)
	}
	return x
}

// Arity of Sin is 1.
func (sin Sin) Arity() uint {
	return 1
}

// Operand returns Sin's operand or nil.
func (sin Sin) Operand(i uint) Operator {
	if i == 0 {
		return sin.Op
	}
	return nil
}

// SetOperand replaces Sin's operand if i is equal to 0.
func (sin Sin) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		sin.Op = op
	}
	return sin
}

// Simplify Sin.
func (sin Sin) Simplify() Operator {
	sin.Op = sin.Op.Simplify()
	switch op := sin.Op.(type) {
	case Const:
		return Const{math.Sin(op.Value)}
	}
	return sin
}

// Diff computes the following derivative: sin(u)' = u' * cos(u)
func (sin Sin) Diff(i uint) Operator {
	return Mul{sin.Op.Diff(i), Cos{sin.Op}}
}

// Name of Sin is "sin".
func (sin Sin) Name() string {
	return "sin"
}

// String formatting.
func (sin Sin) String() string {
	return fmt.Sprintf("sin(%s)", sin.Op)
}
