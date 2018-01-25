package op

// Variable
import (
	"fmt"
)

// A Variable holds an index that can be used to access a certain value of a
// float64 vector. In other words a Variable represents a feature in a dataset.
type Variable struct {
	Index int
}

// ApplyRow Variable.
func (v Variable) ApplyRow(x []float64) float64 {
	return x[v.Index]
}

// ApplyCols Variable.
func (v Variable) ApplyCols(X [][]float64) []float64 {
	var V = make([]float64, len(X[v.Index]))
	copy(V, X[v.Index])
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
