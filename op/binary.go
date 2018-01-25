package op

import (
	"math"

	"github.com/gonum/floats"
)

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
	if x[1] == 0 {
		return 1
	}
	return x[0] / x[1]
}

// ApplyCols Division.
func (op Division) ApplyCols(X [][]float64) []float64 {
	for i, x := range X[1] {
		if x == 1 {
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

// OR boolean operator.
type OR struct{}

// ApplyRow OR.
func (op OR) ApplyRow(x []float64) float64 {
	if (x[0] == 1) || (x[1] == 1) {
		return 1
	}
	return 0
}

// ApplyCols OR.
func (op OR) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of OR.
func (op OR) Arity() int {
	return 2
}

// String representation of OR.
func (op OR) String() string {
	return "OR"
}

// AND boolean operator.
type AND struct{}

// ApplyRow AND.
func (op AND) ApplyRow(x []float64) float64 {
	if (x[0] == 1) && (x[1] == 1) {
		return 1
	}
	return 0
}

// ApplyCols AND.
func (op AND) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of AND.
func (op AND) Arity() int {
	return 2
}

// String representation of AND.
func (op AND) String() string {
	return "AND"
}

// XOR boolean operator.
type XOR struct{}

// ApplyRow XOR.
func (op XOR) ApplyRow(x []float64) float64 {
	if ((x[0] == 1) && (x[1] == 0)) || ((x[0] == 0) && (x[1] == 1)) {
		return 1
	}
	return 0
}

// ApplyCols XOR.
func (op XOR) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of XOR.
func (op XOR) Arity() int {
	return 2
}

// String representation of XOR.
func (op XOR) String() string {
	return "XOR"
}

// NAND boolean operator.
type NAND struct{}

// ApplyRow NAND.
func (op NAND) ApplyRow(x []float64) float64 {
	if (x[0] == 1) && (x[1] == 1) {
		return 0
	}
	return 1
}

// ApplyCols NAND.
func (op NAND) ApplyCols(X [][]float64) []float64 {
	var Y = make([]float64, len(X[0]))
	for i := range X[0] {
		Y[i] = op.ApplyRow([]float64{X[0][i], X[1][i]})
	}
	return Y
}

// Arity of NAND.
func (op NAND) Arity() int {
	return 2
}

// String representation of NAND.
func (op NAND) String() string {
	return "NAND"
}
