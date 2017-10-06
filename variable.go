package xgp

// Variable
import (
	"fmt"
	"math/rand"
)

// newVariable returns a Variable with an index in range [0, p).
func newVariable(p int, rng *rand.Rand) Variable {
	return Variable{
		Index: rng.Intn(p),
	}
}

// A Variable holds an index that can be used to access a certain value of a
// float64 vector. In other words a Variable represents a feature in a dataset.
type Variable struct {
	Index int
}

// Apply Variable.
func (v Variable) Apply(x []float64) float64 {
	return x[v.Index]
}

// ApplyXT Variable.
func (v Variable) ApplyXT(XT [][]float64) []float64 {
	var V = make([]float64, len(XT[v.Index]))
	copy(V, XT[v.Index])
	return V
}

// Arity of a Variable is 0 because it is a terminal operator.
func (v Variable) Arity() int {
	return 0
}

// String representation of a Variable.
func (v Variable) String() string {
	return fmt.Sprintf("X[%d]", v.Index)
}
