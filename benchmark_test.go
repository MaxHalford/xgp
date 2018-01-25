package koza

import (
	"math/rand"
	"testing"
)

func BenchmarkEvaluateCols(b *testing.B) {
	var (
		n, p = 10000, 10
		X    = make([][]float64, p)
		Y    = make([]float64, n)
	)
	for i := 0; i < p; i++ {
		X[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			X[i][j] = rand.Float64()
		}
		Y[i] = rand.Float64()
	}
	var estimator, _ = NewConfigWithDefaults().NewEstimator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		estimator.Fit(X, Y, nil, nil, nil, nil, false)
	}
}
