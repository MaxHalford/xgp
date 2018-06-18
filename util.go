package xgp

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// newRand returns a new random number generator with a random seed.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// randInt returns a random int in [min, max].
func randInt(min, max int, rng *rand.Rand) int {
	return min + rng.Intn(max-min+1)
}

// sigmoid applies the sigmoid transform.
func sigmoid(y float64) float64 {
	return 1 / (1 + math.Exp(-y))
}

// binary converts a float64 to 0 or 1.
func binary(y float64) float64 {
	if y > 0.5 {
		return 1
	}
	return 0
}

// countDistinct returns the number of unique elements in a slice of float64s.
func countDistinct(x []float64) int {
	var seen = make(map[float64]bool)
	for _, xi := range x {
		if _, ok := seen[xi]; !ok {
			seen[xi] = true
		}
	}
	return len(seen)
}

// fmtDuration the "hours:minutes:seconds" representation of a time.Duration.
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
