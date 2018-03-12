package ensemble

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// newRand returns a new random number generator with a random seed.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// argMax returns the index of the largest value in a float64 slice.
func argMax(x []float64) int {
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

// randInt returns a random int in [min, max).
func randInt(min, max int, rng *rand.Rand) int {
	return min + rng.Intn(max-min)
}

// randFloat64 returns a random float64 in [min, max).
func randFloat64(min, max float64, rng *rand.Rand) float64 {
	return min + rng.Float64()*(max-min)
}

// Compute the sum of a float64 slice.
func sumFloat64s(floats []float64) (sum float64) {
	for _, v := range floats {
		sum += v
	}
	return
}

// Compute the mean of a float64 slice.
func meanFloat64s(floats []float64) float64 {
	return sumFloat64s(floats) / float64(len(floats))
}

// Divide each element in a float64 slice by it's sum.
func normalizeFloat64s(floats []float64) []float64 {
	var (
		sum        = sumFloat64s(floats)
		normalized = make([]float64, len(floats))
	)
	for i, f := range floats {
		normalized[i] = f / sum
	}
	return normalized
}

// Compute the cumulative sum of a float64 slice.
func cumsum(floats []float64) []float64 {
	var summed = make([]float64, len(floats))
	copy(summed, floats)
	for i := 1; i < len(summed); i++ {
		summed[i] += summed[i-1]
	}
	return summed
}

// Sample k integers from a list of weights of size n. Sampling is done with
// replacement.
func sampleIndices(k int, weights []float64, rng *rand.Rand) []int {
	var (
		sample = make([]int, k)
		wheel  = cumsum(weights)
		a      = wheel[0]
		b      = wheel[len(wheel)-1]
	)
	for i := range sample {
		sample[i] = sort.SearchFloat64s(wheel, randFloat64(a, b, rng))
	}
	return sample
}

func subsetFloat64Matrix(X [][]float64, rowIdxs []int, colIdxs []int) [][]float64 {
	var S = make([][]float64, len(colIdxs))
	for i, c := range colIdxs {
		S[i] = make([]float64, len(rowIdxs))
		for j, r := range rowIdxs {
			S[i][j] = X[c][r]
		}
	}
	return S
}

func subsetFloat64Slice(X []float64, idxs []int) []float64 {
	var S = make([]float64, len(idxs))
	for i, idx := range idxs {
		S[i] = X[idx]
	}
	return S
}
