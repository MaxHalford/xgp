package op

import "math"

// Log computes the natural logarithm of an operand.
type Log struct{}

// ApplyRow Log.
func (op Log) ApplyRow(x []float64) float64 {
	return math.Log(x[0])
}

// ApplyCols Log.
func (op Log) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i, x := range X[0] {
		Y[i] = math.Log(x)
	}
	return Y
}

// Arity of Log.
func (op Log) Arity() int {
	return 1
}

// String representation of Log.
func (op Log) String() string {
	return "log"
}
