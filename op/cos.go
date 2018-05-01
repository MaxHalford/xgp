package op

import "math"

// Cos computes the cosine of an operand.
type Cos struct{}

// ApplyRow Cos.
func (op Cos) ApplyRow(x []float64) float64 {
	return math.Cos(x[0])
}

// ApplyCols Cos.
func (op Cos) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i, x := range X[0] {
		Y[i] = math.Cos(x)
	}
	return Y
}

// Arity of Cos.
func (op Cos) Arity() int {
	return 1
}

// String representation of Cos.
func (op Cos) String() string {
	return "cos"
}
