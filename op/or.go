package op

// OR boolean operator.
type OR struct{}

// Eval OR.
func (op OR) Eval(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if (X[0][i] == 1) || (X[1][i] == 1) {
			Y[i] = 1
		}
	}
	return Y
}

// Arity of OR.
func (op OR) Arity() int {
	return 2
}

// String representation of OR.
func (op OR) String() string {
	return "OR"
}
