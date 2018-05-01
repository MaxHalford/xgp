package op

import "math"

// Sin computes the sine of an operand.
type Sin struct{}

// ApplyRow Sin.
func (op Sin) ApplyRow(x []float64) float64 {
	return math.Sin(x[0])
}

// ApplyCols Sin.
func (op Sin) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i, x := range X[0] {
		Y[i] = math.Sin(x)
	}
	return Y
}

// Arity of Sin.
func (op Sin) Arity() int {
	return 1
}

// String representation of Sin.
func (op Sin) String() string {
	return "sin"
}
