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

// checkEarlyStop checks if early stopping should occur at round i given k
// early stopping rounds.
func checkEarlyStop(scores []float64, metric metrics.Metric, i, k uint) bool {
	return (metric.BiggerIsBetter() && (scores[i] <= scores[i-k])) ||
		(!metric.BiggerIsBetter() && (scores[i] >= scores[i-k]))
}

func mean(x []float64) (m float64) {
	for _, xi := range x {
		m += xi
	}
	return m / float64(len(x))
}
