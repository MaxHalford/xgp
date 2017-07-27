package xgp

// Variable
import (
	"fmt"
)

type Variable struct {
	Index int
}

func (v Variable) Apply(X []float64) float64 {
	return X[v.Index]
}

func (v Variable) Arity() int {
	return 0
}

func (v Variable) String() string {
	return fmt.Sprintf("X[%d]", v.Index)
}
