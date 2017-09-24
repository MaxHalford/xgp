package xgp

import (
	"fmt"
	"math"
)

// GetTransform returns a Transform from it's String representation.
func GetTransform(transformName string) (Transform, error) {
	var transform, ok = map[string]Transform{
		Identity{}.String(): Identity{},
		Binary{}.String():   Binary{},
		Sigmoid{}.String():  Sigmoid{},
	}[transformName]
	if !ok {
		return nil, fmt.Errorf("Unknown transform name '%s'", transformName)
	}
	return transform, nil
}

// A Transform is a 1D function which can be applied to a float64.
type Transform interface {
	Apply(x float64) float64
	String() string
}

// Identity returns a float64 without any modification.
type Identity struct{}

// Apply an Identity transform.
func (t Identity) Apply(x float64) float64 {
	return x
}

// String representation of an Identity transform.
func (t Identity) String() string {
	return "identity"
}

// Binary returns 0 if a float64 is negative and 1 if not.
type Binary struct{}

// Apply an Binary transform.
func (t Binary) Apply(x float64) float64 {
	if x < 0 {
		return 0
	}
	return 1
}

// String representation of a Binary transform.
func (t Binary) String() string {
	return "binary"
}

// Sigmoid applies the sigmoid function to a float64.
type Sigmoid struct{}

// Apply an Sigmoid transform.
func (t Sigmoid) Apply(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// String representation of a Sigmoid transform.
func (t Sigmoid) String() string {
	return "sigmoid"
}
