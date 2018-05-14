package op

import "github.com/gonum/floats"

// Sub returns the difference between two operands.
type Sub struct{}

// Eval Sub.
func (op Sub) Eval(X [][]float64) []float64 {
	floats.Sub(X[0], X[1])
	return X[0]
}

// Arity of Sub.
func (op Sub) Arity() int {
	return 2
}

// String representation of Sub.
func (op Sub) String() string {
	return "sub"
}
