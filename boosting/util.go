package boosting

import "math"

// ArgMax returns the index of the largest value in a float64 slice.
func ArgMax(x []float64) int {
	var (
		arg int
		max = math.Inf(-1)
	)
	for i, f := range x {
		if f > max {
			arg = i
			max = f
		}
	}
	return arg
}
