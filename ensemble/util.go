package ensemble

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

// Sample k ints in [min, max) using Vitter's reservoir sampling algorithm if
// bootstrap is false.
func randomInts(k int, min, max int, bootstrap bool, rng *rand.Rand) ([]int, error) {
	if max < min {
		return nil, fmt.Errorf("max has to be greater or equal to min: %d < %d", max, min)
	}
	var ints = make([]int, k)
	// With replacement
	if bootstrap {
		for i := 0; i < k; i++ {
			ints[i] = randInt(min, max, rng)
		}
		return ints, nil
	}
	// Without replacement
	for i := 0; i < k; i++ {
		ints[i] = i + min
	}
	for i := k; i < max-min; i++ {
		var j = rng.Intn(i + 1)
		if j < k {
			ints[j] = i + min
		}
	}
	return ints, nil
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

func sample(
	X [][]float64,
	Y []float64,
	W []float64,
	rowSampling float64,
	colSampling float64,
	boostrapRows bool,
	bootstrapCols bool,
	rng *rand.Rand,
) ([][]float64, []float64, []float64, []int, []int, error) {

	// Sample row indexes
	var n = int(rowSampling * float64(len(X[0])))
	rowIdxs, err := randomInts(n, 0, len(X[0]), boostrapRows, rng)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Sample column indexes
	var p = int(colSampling * float64(len(X)))
	colIdxs, err := randomInts(p, 0, len(X), bootstrapCols, rng)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Create the sample
	var (
		XSam = make([][]float64, len(colIdxs))
		YSam = make([]float64, len(rowIdxs))
		WSam []float64
	)
	if W != nil {
		WSam = make([]float64, len(rowIdxs))
	}
	for i := range colIdxs {
		XSam[i] = make([]float64, len(rowIdxs))
	}
	for i, r := range rowIdxs {
		for j, c := range colIdxs {
			XSam[j][i] = X[c][r]
		}
		YSam[i] = Y[r]
		if W != nil {
			WSam[i] = W[r]
		}
	}
	if W != nil {
		return XSam, YSam, WSam, rowIdxs, colIdxs, nil
	}
	return XSam, YSam, nil, rowIdxs, colIdxs, nil
}
