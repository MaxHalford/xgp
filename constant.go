package xgp

import (
	"fmt"
)

// A Constant holds a float64 value.
type Constant struct {
	Value float64
}

// Apply Constant.
func (c Constant) Apply(X []float64) float64 {
	return c.Value
}

// Arity of a Constant.
func (c Constant) Arity() int {
	return 0
}

// String representation of a Constant.
func (c Constant) String() string {
	return fmt.Sprintf("%.2f", c.Value)
}
