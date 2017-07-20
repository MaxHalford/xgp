package dataframe

import "math/rand"

// Sample k unique integers in range [min, max) using reservoir sampling,
// specifically Vitter's Algorithm R.
func randomInts(k, min, max int, rng *rand.Rand) (ints []int) {
	ints = make([]int, k)
	for i := 0; i < k; i++ {
		ints[i] = i + min
	}
	for i := k; i < max-min; i++ {
		var j = rng.Intn(i + 1)
		if j < k {
			ints[j] = i + min
		}
	}
	return
}

func containsInts(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
