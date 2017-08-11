package xgp

// Variable
import (
	"fmt"
)

// A Variable holds an index that can be used to access a certain value of a
// float64 vector. In other words a Variable represents a feature in a dataset.
type Variable struct {
	Index int
}

// Apply a Variable.
func (v Variable) Apply(X []float64) float64 {
	return X[v.Index]
}

// Arity of a Variable.
func (v Variable) Arity() int {
	return 0
}

// String representation of a Variable.
func (v Variable) String() string {
	return fmt.Sprintf("X[%d]", v.Index)
}
