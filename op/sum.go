package op

import "github.com/gonum/floats"

// Sum returns the sum of two operands.
type Sum struct{}

// Eval Sum.
func (op Sum) Eval(X [][]float64) []float64 {
	floats.Add(X[0], X[1])
	return X[0]
}

// Arity of Sum.
func (op Sum) Arity() int {
	return 2
}

// String representation of Sum.
func (op Sum) String() string {
	return "sum"
}
