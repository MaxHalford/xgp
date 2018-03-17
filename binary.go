package xgp

import "math"

// sigmoid applies the sigmoid transform.
func sigmoid(y float64) float64 {
	return 1 / (1 + math.Exp(-y))
}

// binary converts a
func binary(y float64) float64 {
	if y > 0.5 {
		return 1
	}
	return 0
}
