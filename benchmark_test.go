package xgp

import (
	"testing"

	"github.com/MaxHalford/xgp/dataset"
)

func BenchmarkFit(b *testing.B) {
	var estimator, err = NewEstimator(
		10,
		-10,
		"mae",
		"sum,sub,mul,div",
		"mae",
		3,
		6,
		10,
		0,
		0.5,
		0.3,
		42,
		0,
	)
	if err != nil {
		panic(err)
	}
	// Load the training set in memory
	boston, err := dataset.ReadCSV("examples/boston/train.csv", "y", false)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		estimator.Fit(boston.X, boston.Y, boston.XNames, false)
	}
}
