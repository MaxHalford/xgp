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
