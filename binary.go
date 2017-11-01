package xgp

import "math"

// sigmoid applies the sigmoid transform to each value in a slice.
func sigmoid(Y []float64) []float64 {
	var Z = make([]float64, len(Y))
	for i, y := range Y {
		Z[i] = 1 / (1 + math.Exp(-y))
	}
	return Z
}

// binary converts each value in a slice to 0 or 1.
func binary(Y []float64) []float64 {
	var Z = make([]float64, len(Y))
	for i, y := range Y {
		if y > 0.5 {
			Z[i] = 1
			continue
		}
		Z[i] = 0
	}
	return Z
}
