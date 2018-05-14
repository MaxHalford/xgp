package meta

import (
	"math/rand"
	"sort"
)

// Sample with repetition k integers in range [0, n).
func randInts(k, n int, rng *rand.Rand) []int {
	var ints = make([]int, k)
	for i := 0; i < k; i++ {
		ints[i] = rng.Intn(n)
	}
	return ints
}

// Sample k unique integers in range [0, n) using reservoir sampling
func randIntsNoRep(k, n int, rng *rand.Rand) []int {
	var ints = make([]int, k)
	for i := range ints {
		ints[i] = i
	}
	// If k == n then simply all the integers in range [0, n]
	if k == n {
		return ints
	}
	// If not use reservoir sampling
	for i := k; i < n; i++ {
		var j = rng.Intn(i + 1)
		if j < k {
			ints[j] = i
		}
	}
	return ints
}

// Compute the cumulative sum of a float64 slice.
func cumsum(floats []float64) []float64 {
	var cs = make([]float64, len(floats))
	copy(cs, floats)
	for i := 1; i < len(cs); i++ {
		cs[i] += cs[i-1]
	}
	return cs
}

// Sample with repetition k integers in range [0, n) where n is the number of
// provided weights.
func randIntsWeighted(k int, weights []float64, rng *rand.Rand) []int {
	var (
		ints = make([]int, k)
		wcs  = cumsum(weights)
	)
	for i := range ints {
		r := rng.Float64() * wcs[len(wcs)-1]
		ints[i] = sort.SearchFloat64s(wcs, r)
	}
	return ints
}

func subsetFloat64Slice(X []float64, idxs []int) []float64 {
	var S = make([]float64, len(idxs))
	for i, idx := range idxs {
		S[i] = X[idx]
	}
	return S
}

func subsetColsFloat64Matrix(X [][]float64, idxs []int) [][]float64 {
	var S = make([][]float64, len(idxs))
	for i, c := range idxs {
		S[i] = make([]float64, len(X[c]))
		for j, x := range X[c] {
			S[i][j] = x
		}
	}
	return S
}

func subsetRowsFloat64Matrix(X [][]float64, idxs []int) [][]float64 {
	var S = make([][]float64, len(X))
	for i := range X {
		S[i] = make([]float64, len(idxs))
		for j, r := range idxs {
			S[i][j] = X[i][r]
		}
	}
	return S
}

func normalize(w []float64) []float64 {
	var s float64
	for _, wi := range w {
		s += wi
	}
	for i := range w {
		w[i] /= s
	}
	return w
}
