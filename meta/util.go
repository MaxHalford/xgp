package meta

import (
	"math/rand"
	"time"

	"github.com/MaxHalford/xgp/metrics"
)

// newRand returns a new random number generator with a random seed.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func checkEarlyStop(scores []float64, rounds, i uint, metric metrics.Metric) bool {
	return (metric.BiggerIsBetter() && (scores[i] <= scores[i-rounds])) ||
		(!metric.BiggerIsBetter() && (scores[i] >= scores[i-rounds]))
}
