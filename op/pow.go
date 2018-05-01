package op

import "math"

// Pow computes the exponent of a first value by a second one.
type Pow struct{}

// ApplyRow Pow.
func (op Pow) ApplyRow(X []float64) float64 {
	return math.Pow(X[0], X[1])
}

// ApplyCols Pow.
func (op Pow) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = math.Pow(X[0][i], X[1][i])
	}
	return Y
}

// Arity of Pow.
func (op Pow) Arity() int {
	return 2
}

// String representation of Pow.
func (op Pow) String() string {
	return "pow"
}
