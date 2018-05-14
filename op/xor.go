package op

// XOR boolean operator.
type XOR struct{}

// Eval XOR.
func (op XOR) Eval(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if ((X[0][i] == 1) && (X[1][i] != 1)) || ((X[0][i] != 1) && (X[1][i] == 1)) {
			Y[i] = 1
		}
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
