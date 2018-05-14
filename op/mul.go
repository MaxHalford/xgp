package op

import "github.com/gonum/floats"

// Mul returns the product two operands.
type Mul struct{}

// Eval Mul.
func (op Mul) Eval(X [][]float64) []float64 {
	floats.Mul(X[0], X[1])
	return X[0]
}

// Arity of Mul.
func (op Mul) Arity() int {
	return 2
}

// String representation of Mul.
func (op Mul) String() string {
	return "mul"
}
