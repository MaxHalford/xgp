package gp

import (
	"math"
)

// 1D functions

// Cos computes the cosine of an operand.
type Cos struct{}

func (op Cos) Apply(X []float64) float64 {
	return math.Cos(X[0])
}

func (op Cos) Arity() int {
	return 1
}

func (op Cos) String() string {
	return "cos"
}

// Sin computes the sine of an operand.
type Sin struct{}

func (op Sin) Apply(X []float64) float64 {
	return math.Sin(X[0])
}

func (op Sin) Arity() int {
	return 1
}

func (op Sin) String() string {
	return "sin"
}

// Log

type Log struct{}

func (op Log) Apply(X []float64) float64 {
	return math.Log(X[0])
}

func (op Log) Arity() int {
	return 1
}

func (op Log) String() string {
	return "log"
}

// Exp

type Exp struct{}

func (op Exp) Apply(X []float64) float64 {
	return math.Exp(X[0])
}

func (op Exp) Arity() int {
	return 1
}

func (op Exp) String() string {
	return "exp"
}

// 2D operators

// Max

type Max struct{}

func (op Max) Apply(X []float64) float64 {
	if X[0] > X[1] {
		return X[0]
	}
	return X[1]
}

func (op Max) Arity() int {
	return 2
}

func (op Max) String() string {
	return "max"
}

// Min

type Min struct{}

func (op Min) Apply(X []float64) float64 {
	if X[0] < X[1] {
		return X[0]
	}
	return X[1]
}

func (op Min) Arity() int {
	return 2
}

func (op Min) String() string {
	return "min"
}

// Sum

type Sum struct{}

func (op Sum) Apply(X []float64) float64 {
	return X[0] + X[1]
}

func (op Sum) Arity() int {
	return 2
}

func (op Sum) String() string {
	return "+"
}

// Difference

type Difference struct{}

func (op Difference) Apply(X []float64) float64 {
	return X[0] - X[1]
}

func (op Difference) Arity() int {
	return 2
}

func (op Difference) String() string {
	return "-"
}

// Division returns the division of two operands. The left operand is the
// numerator and the right operand is the denominator. The division is protected
// so that if the denominator's value is in range [-0.001, 0.001] the operator
// returns 1.
type Division struct{}

func (op Division) Apply(X []float64) float64 {
	if math.Abs(X[1]) < 0.001 {
		return 1
	}
	return X[0] / X[1]
}

func (op Division) Arity() int {
	return 2
}

func (op Division) String() string {
	return "/"
}

// Product returns the product of the operands.
type Product struct{}

func (op Product) Apply(X []float64) float64 {
	return X[0] * X[1]
}

func (op Product) Arity() int {
	return 2
}

func (op Product) String() string {
	return "*"
}

// Power

type Power struct{}

func (op Power) Apply(X []float64) float64 {
	return math.Pow(X[0], X[1])
}

func (op Power) Arity() int {
	return 2
}

func (op Power) String() string {
	return "^"
}

// Greater than

type GreaterThan struct{}

func (op GreaterThan) Apply(X []float64) float64 {
	if X[0] > X[1] {
		return 1
	}
	return 0
}

func (op GreaterThan) Arity() int {
	return 2
}

func (op GreaterThan) String() string {
	return ">"
}

// Lesser than

type LesserThan struct{}

func (op LesserThan) Apply(X []float64) float64 {
	if X[0] < X[1] {
		return 1
	}
	return 0
}

func (op LesserThan) Arity() int {
	return 2
}

func (op LesserThan) String() string {
	return "<"
}
