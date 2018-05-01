package op

// NAND boolean operator.
type NAND struct{}

// ApplyRow NAND.
func (op NAND) ApplyRow(x []float64) float64 {
	if (x[0] == 1) && (x[1] == 1) {
		return 0
	}
	return 1
}

// ApplyCols NAND.
func (op NAND) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of NAND.
func (op NAND) Arity() int {
	return 2
}

// String representation of NAND.
func (op NAND) String() string {
	return "NAND"
}
