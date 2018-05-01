package op

// Min returns the minimum of two operands.
type Min struct{}

// ApplyRow Min.
func (op Min) ApplyRow(x []float64) float64 {
	if x[0] < x[1] {
		return x[0]
	}
	return x[1]
}

// ApplyCols Min.
func (op Min) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if X[0][i] < X[1][i] {
			Y[i] = X[0][i]
		} else {
			Y[i] = X[1][i]
		}
	}
	return Y
}

// Arity of Min.
func (op Min) Arity() int {
	return 2
}

// String representation of Min.
func (op Min) String() string {
	return "min"
}
