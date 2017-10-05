package xgp

import (
	"math"
	"math/rand"
	"time"
)

// randFloat64 returns a random float64 in [min, max].
func randFloat64(min, max float64, rng *rand.Rand) float64 {
	return min + rng.Float64()*(max-min)
}

// randInt returns a random int in [min, max].
func randInt(min, max int, rng *rand.Rand) int {
	return min + rng.Intn(max-min+1)
}

// meanFloat64s returns the mean of a float64 slice.
func meanFloat64s(fs []float64) float64 {
	var sum float64
	for _, f := range fs {
		sum += f
	}
	return sum / float64(len(fs))
}

// varianceFloat64s returns the variance of a float64 slice.
func varianceFloat64s(fs []float64) float64 {
	var (
		m  = meanFloat64s(fs)
		ss float64
	)
	for _, x := range fs {
		ss += math.Pow(x, 2)
	}
	return ss/float64(len(fs)) - math.Pow(m, 2)
}

// makeRNG returns a new random number generator with a random seed.
func makeRNG() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
