package koza

import (
	"math/rand"
	"time"
)

// newRand returns a new random number generator with a random seed.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// countDistinct returns the number of unique elements in a slice of float64s.
func countDistinct(X []float64) int {
	var seen = make(map[float64]bool)
	for _, x := range X {
		if _, ok := seen[x]; !ok {
			seen[x] = true
		}
	}
	return len(seen)
}
