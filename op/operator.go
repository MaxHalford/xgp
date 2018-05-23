package op

// An Operator takes float64s as input and outputs a float64.
type Operator interface {
	Eval(X [][]float64) []float64
	Arity() int
	String() string
}
