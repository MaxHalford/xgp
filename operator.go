package xgp

// An Operator takes float64s as input and outputs a float64.
type Operator interface {
	Apply(x []float64) float64
	ApplyXT(XT [][]float64) []float64
	Arity() int
	String() string
}
