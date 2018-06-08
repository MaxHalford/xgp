package op

import "fmt"

// The Neg operator.
type Neg struct {
	Op Operator
}

// Eval computes the opposite of each value.
func (neg Neg) Eval(X [][]float64) []float64 {
	x := neg.Op.Eval(X)
	for i, xi := range x {
		x[i] = -xi
	}
	return x
}

// Arity of Neg is 1.
func (neg Neg) Arity() uint {
	return 1
}

// Operand returns Neg's operand or nil.
func (neg Neg) Operand(i uint) Operator {
	if i == 0 {
		return neg.Op
	}
	return nil
}

// SetOperand replaces Neg's operand if i is equal to 0.
func (neg Neg) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		neg.Op = op
	}
	return neg
}

// Simplify Neg.
func (neg Neg) Simplify() Operator {
	switch operand := neg.Op.(type) {
	case Neg:
		return operand.Op
	default:
		return neg
	}
}

// Diff compute the following derivative: (-u)' = -u'.
func (neg Neg) Diff(i uint) Operator {
	return Neg{neg.Op.Diff(i)}
}

// Name of Neg is "neg".
func (neg Neg) Name() string {
	return "neg"
}

// String formatting.
func (neg Neg) String() string {
	return fmt.Sprintf("-%s", parenthesize(neg.Op))
}
