package xgp

type Operator interface {
	Apply(X []float64) float64
	Arity() int
	String() string
}
