package meta

import (
	"fmt"
	"math/rand"
	"time"
)

// mean computes the mean of a float64 slice.
func mean(x []float64) (m float64) {
	for _, xi := range x {
		m += xi
	}
	return m / float64(len(x))
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

// Sample k unique integers in range [min, max) using reservoir sampling,
// specifically Vitter's Algorithm R.
func randomInts(k uint, min, max int, rng *rand.Rand) []int {
	var ints = make([]int, k)
	for i := 0; i < int(k); i++ {
		ints[i] = i + min
	}
	for i := int(k); i < max-min; i++ {
		var j = rng.Intn(i + 1)
		if j < int(k) {
			ints[j] = i + min
		}
	}
	return ints
}

func selectCols(X [][]float64, cols []int) [][]float64 {
	var XX = make([][]float64, len(cols))
	for i, c := range cols {
		XX[i] = X[c]
	}
	return XX
}

func selectRows(X [][]float64, rows []int) [][]float64 {
	var XX = make([][]float64, len(X))
	for i, col := range X {
		XX[i] = make([]float64, len(col))
		for j, r := range rows {
			XX[i][j] = X[i][r]
		}
	}
	return XX
}
