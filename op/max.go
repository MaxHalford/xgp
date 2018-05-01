package op

// Max returns the maximum of two operands.
type Max struct{}

// ApplyRow Max.
func (op Max) ApplyRow(x []float64) float64 {
	if x[0] > x[1] {
		return x[0]
	}
	return x[1]
}

// ApplyCols Max.
func (op Max) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if X[0][i] > X[1][i] {
			Y[i] = X[0][i]
		} else {
			Y[i] = X[1][i]
		}
	}
	return Y
}

// Arity of Max.
func (op Max) Arity() int {
	return 2
}

// String representation of Max.
func (op Max) String() string {
	return "max"
}
