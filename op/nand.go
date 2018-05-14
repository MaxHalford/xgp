package op

// NAND boolean operator.
type NAND struct{}

// Eval NAND.
func (op NAND) Eval(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if (X[0][i] == 1) && (X[1][i] == 1) {
			Y[i] = 0
		} else {
			Y[i] = 1
		}
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
