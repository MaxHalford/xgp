package xgp

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

// 1D functions

// Cos computes the cosine of an operand.
type Cos struct{}

// Apply Cos.
func (op Cos) Apply(x []float64) float64 {
	return math.Cos(x[0])
}

// ApplyXT Cos.
func (op Cos) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i, x := range XT[0] {
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

// Sin computes the sine of an operand.
type Sin struct{}

// Apply Sin.
func (op Sin) Apply(X []float64) float64 {
	return math.Sin(X[0])
}

// ApplyXT Sin.
func (op Sin) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i, x := range XT[0] {
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

// Log computes the natural logarithm of an operand.
type Log struct{}

// Apply Log.
func (op Log) Apply(X []float64) float64 {
	return math.Log(X[0])
}

// ApplyXT Log.
func (op Log) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i, x := range XT[0] {
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

// Exp computes the exponential of an operand.
type Exp struct{}

// Apply Exp.
func (op Exp) Apply(X []float64) float64 {
	return math.Exp(X[0])
}

// ApplyXT Exp.
func (op Exp) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i, x := range XT[0] {
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

// 2D operators

// Max returns the maximum of two operands.
type Max struct{}

// Apply Max.
func (op Max) Apply(X []float64) float64 {
	if X[0] > X[1] {
		return X[0]
	}
	return X[1]
}

// ApplyXT Max.
func (op Max) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i := range XT[0] {
		if XT[0][i] > XT[1][i] {
			Y[i] = XT[0][i]
		} else {
			Y[i] = XT[1][i]
		}
	}
	return Y
}

// Arity of Max.
func (op Max) Arity() int {
	return 2
}

// String representation of Max.
func (op Max) String() string {
	return "max"
}

// Min returns the minimum of two operands.
type Min struct{}

// Apply Min.
func (op Min) Apply(X []float64) float64 {
	if X[0] < X[1] {
		return X[0]
	}
	return X[1]
}

// ApplyXT Min.
func (op Min) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i := range XT[0] {
		if XT[0][i] < XT[1][i] {
			Y[i] = XT[0][i]
		} else {
			Y[i] = XT[1][i]
		}
	}
	return Y
}

// Arity of Min.
func (op Min) Arity() int {
	return 2
}

// String representation of Min.
func (op Min) String() string {
	return "min"
}

// Sum returns the sum of two operands.
type Sum struct{}

// Apply Sum.
func (op Sum) Apply(X []float64) float64 {
	return X[0] + X[1]
}

// ApplyXT Sum.
func (op Sum) ApplyXT(XT [][]float64) []float64 {
	floats.Add(XT[0], XT[1])
	return XT[0]
}

// Arity of Sum.
func (op Sum) Arity() int {
	return 2
}

// String representation of String.
func (op Sum) String() string {
	return "+"
}

// Difference returns the difference between two operands.
type Difference struct{}

// Apply Difference.
func (op Difference) Apply(X []float64) float64 {
	return X[0] - X[1]
}

// ApplyXT Difference.
func (op Difference) ApplyXT(XT [][]float64) []float64 {
	floats.Sub(XT[0], XT[1])
	return XT[0]
}

// Arity of Difference.
func (op Difference) Arity() int {
	return 2
}

// String representation of Difference.
func (op Difference) String() string {
	return "-"
}

// Division returns the division of two operands. The left operand is the
// numerator and the right operand is the denominator. The division is protected
// so that if the denominator's value is in range [-0.001, 0.001] the operator
// returns 1.
type Division struct{}

// Apply Division.
func (op Division) Apply(X []float64) float64 {
	if math.Abs(X[1]) < 0.001 {
		return 1
	}
	return X[0] / X[1]
}

// ApplyXT Division.
func (op Division) ApplyXT(XT [][]float64) []float64 {
	floats.Div(XT[0], XT[1])
	return XT[0]
}

// Arity of Division.
func (op Division) Arity() int {
	return 2
}

// String representation of Division.
func (op Division) String() string {
	return "/"
}

// Product returns the product two operands.
type Product struct{}

// Apply Product.
func (op Product) Apply(X []float64) float64 {
	return X[0] * X[1]
}

// ApplyXT Product.
func (op Product) ApplyXT(XT [][]float64) []float64 {
	floats.Mul(XT[0], XT[1])
	return XT[0]
}

// Arity of Product.
func (op Product) Arity() int {
	return 2
}

// String representation of Product.
func (op Product) String() string {
	return "*"
}

// Power computes the exponent of a first value by a second one.
type Power struct{}

// Apply Power.
func (op Power) Apply(X []float64) float64 {
	return math.Pow(X[0], X[1])
}

// ApplyXT Power.
func (op Power) ApplyXT(XT [][]float64) []float64 {
	var Y = make([]float64, len(XT[0]))
	for i := range XT[0] {
		Y[i] = math.Pow(XT[0][i], XT[1][i])
	}
	return Y
}

// Arity of Power.
func (op Power) Arity() int {
	return 2
}

// String representation of Power.
func (op Power) String() string {
	return "^"
}
