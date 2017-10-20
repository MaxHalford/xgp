package tree

import (
	"math/rand"
	"time"
)

// newRand returns a new random number generator with a random seed.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// randFloat64 returns a random float64 in [min, max].
func randFloat64(min, max float64, rng *rand.Rand) float64 {
	return min + rng.Float64()*(max-min)
}

// randInt returns a random int in [min, max].
func randInt(min, max int, rng *rand.Rand) int {
	return min + rng.Intn(max-min+1)
}
