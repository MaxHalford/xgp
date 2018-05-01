package op

import "math"

// Exp computes the exponential of an operand.
type Exp struct{}

// ApplyRow Exp.
func (op Exp) ApplyRow(x []float64) float64 {
	return math.Exp(x[0])
}

// ApplyCols Exp.
func (op Exp) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i, x := range X[0] {
		Y[i] = math.Exp(x)
	}
	return Y
}

// Arity of Exp.
func (op Exp) Arity() int {
	return 1
}

// String representation of Exp.
func (op Exp) String() string {
	return "exp"
}
