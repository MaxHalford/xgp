package tree

// A Weighting is a convinience structure for assigning weights to Operators
// for selection purposes.
type Weighting struct {
	PConstant float64
	PVariable float64
	PFunction float64
}

func (w Weighting) apply(op Operator) float64 {
	switch op.(type) {
	case Constant:
		return w.PConstant
	case Variable:
		return w.PVariable
	default:
		return w.PFunction
	}
}
