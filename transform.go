package xgp

import "math"

// Identity returns a float64 without any modification.
func Identity(x float64) float64 {
	return x
}

// Binary returns 0 if a float64 is negative and 1 if not.
func Binary(x float64) float64 {
	if x < 0 {
		return 0
	}
	return 1
}

// Sigmoid applies the sigmoid function to a float64.
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}
