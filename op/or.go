package op

// OR boolean operator.
type OR struct{}

// ApplyRow OR.
func (op OR) ApplyRow(x []float64) float64 {
	if (x[0] == 1) || (x[1] == 1) {
		return 1
	}
	return 0
}

// ApplyCols OR.
func (op OR) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
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
