package tree

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/floats"
)

// parseFuncName returns a functional Operator from it's String representation.
func parseFuncName(funcName string) (Operator, error) {
	var f, ok = map[string]Operator{
		Cos{}.String():        Cos{},
		Sin{}.String():        Sin{},
		Log{}.String():        Log{},
		Exp{}.String():        Exp{},
		Max{}.String():        Max{},
		Min{}.String():        Min{},
		Sum{}.String():        Sum{},
		Difference{}.String(): Difference{},
		Division{}.String():   Division{},
		Product{}.String():    Product{},
		Power{}.String():      Power{},
	}[funcName]
	if !ok {
		return nil, fmt.Errorf("Unknown function name '%s'", funcName)
	}
	return f, nil
}

// 1D functions

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

// 2D operators

// Max returns the maximum of two operands.
type Max struct{}

// ApplyRow Max.
func (op Max) ApplyRow(x []float64) float64 {
	if x[0] > x[1] {
		return x[0]
	}
	return x[1]
}

// ApplyCols Max.
func (op Max) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if X[0][i] > X[1][i] {
			Y[i] = X[0][i]
		} else {
			Y[i] = X[1][i]
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

// ApplyRow Min.
func (op Min) ApplyRow(x []float64) float64 {
	if x[0] < x[1] {
		return x[0]
	}
	return x[1]
}

// ApplyCols Min.
func (op Min) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		if X[0][i] < X[1][i] {
			Y[i] = X[0][i]
		} else {
			Y[i] = X[1][i]
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

// ApplyRow Sum.
func (op Sum) ApplyRow(x []float64) float64 {
	return x[0] + x[1]
}

// ApplyCols Sum.
func (op Sum) ApplyCols(X [][]float64) []float64 {
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

// Difference returns the difference between two operands.
type Difference struct{}

// ApplyRow Difference.
func (op Difference) ApplyRow(x []float64) float64 {
	return x[0] - x[1]
}

// ApplyCols Difference.
func (op Difference) ApplyCols(X [][]float64) []float64 {
	floats.Sub(X[0], X[1])
	return X[0]
}

// Arity of Difference.
func (op Difference) Arity() int {
	return 2
}

// String representation of Difference.
func (op Difference) String() string {
	return "sub"
}

// Division returns the protected division of two operands. The left operand is
// the numerator and the right operand is the denominator.
type Division struct{}

// ApplyRow Division.
func (op Division) ApplyRow(x []float64) float64 {
	if math.Abs(x[1]) < math.SmallestNonzeroFloat64 {
		return 1
	}
	return x[0] / x[1]
}

// ApplyCols Division.
func (op Division) ApplyCols(X [][]float64) []float64 {
	for i, x := range X[1] {
		if math.Abs(x) < math.SmallestNonzeroFloat64 {
			X[0][i] = 1
		} else {
			X[0][i] /= x
		}
	}
	return X[0]
}

// Arity of Division.
func (op Division) Arity() int {
	return 2
}

// String representation of Division.
func (op Division) String() string {
	return "div"
}

// Product returns the product two operands.
type Product struct{}

// ApplyRow Product.
func (op Product) ApplyRow(X []float64) float64 {
	return X[0] * X[1]
}

// ApplyCols Product.
func (op Product) ApplyCols(X [][]float64) []float64 {
	floats.Mul(X[0], X[1])
	return X[0]
}

// Arity of Product.
func (op Product) Arity() int {
	return 2
}

// String representation of Product.
func (op Product) String() string {
	return "mul"
}

// Power computes the exponent of a first value by a second one.
type Power struct{}

// ApplyRow Power.
func (op Power) ApplyRow(X []float64) float64 {
	return math.Pow(X[0], X[1])
}

// ApplyCols Power.
func (op Power) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = math.Pow(X[0][i], X[1][i])
	}
	return Y
}

// Arity of Power.
func (op Power) Arity() int {
	return 2
}

// String representation of Power.
func (op Power) String() string {
	return "pow"
}
