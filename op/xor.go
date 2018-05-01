package op

// XOR boolean operator.
type XOR struct{}

// ApplyRow XOR.
func (op XOR) ApplyRow(x []float64) float64 {
	if ((x[0] == 1) && (x[1] == 0)) || ((x[0] == 0) && (x[1] == 1)) {
		return 1
	}
	return 0
}

// ApplyCols XOR.
func (op XOR) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of XOR.
func (op XOR) Arity() int {
	return 2
}

// String representation of XOR.
func (op XOR) String() string {
	return "XOR"
}
