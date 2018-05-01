package op

// Div returns the protected division of two operands. The left operand is
// the numerator and the right operand is the denominator.
type Div struct{}

// ApplyRow Div.
func (op Div) ApplyRow(x []float64) float64 {
	if x[1] == 0 {
		return 1
	}
	return x[0] / x[1]
}

// ApplyCols Div.
func (op Div) ApplyCols(X [][]float64) []float64 {
	for i, x := range X[1] {
		if x == 1 {
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
