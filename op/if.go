package op

import "fmt"

// The If operator.
type If struct {
	Condition Operator
	Lower     Operator
	Upper     Operator
}

// Eval computes
func (iff If) Eval(X [][]float64) []float64 {
	var (
		lower = iff.Lower.Eval(X)
		upper = iff.Upper.Eval(X)
	)
	for i, c := range iff.Condition.Eval(X) {
		if c > 0 {
			lower[i] = upper[i]
		}
	}
	return lower
}

// Arity of If is 3.
func (iff If) Arity() uint {
	return 3
}

// Operand returns If's operand or nil.
func (iff If) Operand(i uint) Operator {
	switch i {
	case 0:
		return iff.Condition
	case 1:
		return iff.Lower
	case 2:
		return iff.Upper
	}
	return nil
}

// SetOperand replaces If's operand if i is equal to 0.
func (iff If) SetOperand(i uint, op Operator) Operator {
	switch i {
	case 0:
		iff.Condition = op
	case 1:
		iff.Lower = op
	case 2:
		iff.Upper = op
	}
	return iff
}

// Simplify If.
func (iff If) Simplify() Operator {
	iff.Condition = iff.Condition.Simplify()
	iff.Lower = iff.Lower.Simplify()
	iff.Upper = iff.Upper.Simplify()
	switch operand := iff.Condition.(type) {
	case Const:
		if operand.Value < 0 {
			return iff.Lower
		}
		return iff.Upper
	}
	return iff
}

// Diff compute the following derivative: (1 / u)' = -u' / u².
func (iff If) Diff(i uint) Operator {
	return iff
}

// Name of If is "if".
func (iff If) Name() string {
	return "if"
}

// String formatting.
func (iff If) String() string {
	return fmt.Sprintf(
		"if(%s < 0 then %s else %s)",
		parenthesize(iff.Condition),
		parenthesize(iff.Lower),
		parenthesize(iff.Upper),
	)
}
