package op

// AND boolean operator.
type AND struct{}

// Eval AND.
func (op AND) Eval(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if X[0][i] == 1 && X[1][i] == 1 {
			Y[i] = 1
		}
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
