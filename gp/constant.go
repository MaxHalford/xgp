package gp

import (
	"fmt"
)

// Constant

type Constant struct {
	Value float64
}

func (c Constant) Apply(X []float64) float64 {
	return c.Value
}

func (c Constant) Arity() int {
	return 0
}

func (c Constant) String() string {
	return fmt.Sprintf("%.2f", c.Value)
}
