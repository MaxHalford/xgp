package tree

import "github.com/MaxHalford/koza/tree/op"

// A Weighting is a convinience structure for assigning weights to Operators
// for selection purposes.
type Weighting struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

func (w Weighting) apply(operator op.Operator) float64 {
	switch operator.(type) {
	case op.Constant:
		return w.PConstant
	case op.Variable:
		return w.PVariable
	default:
		return w.PFunction
	}
}
