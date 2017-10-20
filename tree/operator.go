package tree

// An Operator takes float64s as input and outputs a float64.
type Operator interface {
	ApplyRow(x []float64) float64
	ApplyCols(X [][]float64) []float64
	Arity() int
	String() string
}
