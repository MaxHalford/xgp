package xgp

// An Operator takes float64s as input and outputs a float64.
type Operator interface {
	Apply(X []float64) float64
	Arity() int
	String() string
}
