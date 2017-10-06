package xgp

import (
	"fmt"
	"math/rand"

	"gonum.org/v1/gonum/floats"
)

// newConstant returns a Constant whose value is sampled from [min, max].
func newConstant(min, max float64, rng *rand.Rand) Constant {
	return Constant{
		Value: randFloat64(min, max, rng),
	}
}

// A Constant holds a float64 value.
type Constant struct {
	Value float64
}

// Apply Constant.
func (c Constant) Apply(x []float64) float64 {
	return c.Value
}

// ApplyXT Constant.
func (c Constant) ApplyXT(XT [][]float64) []float64 {
	var C = make([]float64, len(XT[0]))
	floats.AddConst(c.Value, C)
	return C
}

// Arity of a Constant is 0 because it is a terminal operator.
func (c Constant) Arity() int {
	return 0
}

// String representation of a Constant.
func (c Constant) String() string {
	return fmt.Sprintf("%.2f", c.Value)
}
