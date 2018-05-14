package op

// Div returns the protected division of two operands. The left operand is
// the numerator and the right operand is the denominator.
type Div struct{}

// Eval Div.
func (op Div) Eval(X [][]float64) []float64 {
	for i, x := range X[1] {
		if x == 0 {
			X[0][i] = 1
		} else {
			X[0][i] /= x
		}
	}
	return X[0]
}

// Arity of Div.
func (op Div) Arity() int {
	return 2
}

// String representation of Div.
func (op Div) String() string {
	return "div"
}
