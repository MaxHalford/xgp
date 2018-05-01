package op

// AND boolean operator.
type AND struct{}

// ApplyRow AND.
func (op AND) ApplyRow(x []float64) float64 {
	if (x[0] == 1) && (x[1] == 1) {
		return 1
	}
	return 0
}

// ApplyCols AND.
func (op AND) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of AND.
func (op AND) Arity() int {
	return 2
}

// String representation of AND.
func (op AND) String() string {
	return "AND"
}
