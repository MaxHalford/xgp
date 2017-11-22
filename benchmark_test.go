package koza

import (
	"testing"

	"github.com/MaxHalford/koza/dataset"
)

func BenchmarkFit(b *testing.B) {
	var estimator, err = NewEstimator(
		10,                // constMax
		-10,               // constMin
		"mae",             // evalMetric
		"sum,sub,mul,div", // funcs
		10,                // generations
		"mae",             // lossMetric
		6,                 // maxHeight
		3,                 // minHeight
		1,                 // nPops
		0,                 // parsimonyCoeff
		0.5,               // pConstant
		0.5,               // pFull
		0.01,              // pHoistMutation
		0.01,              // pPointMutation
		0.9,               // pSubTreeCrossover
		0.01,              // PSubTreeMutation
		0.3,               // pTerminal
		30,                // popSize
		42,                // seed
		0,                 // tuningGenerations
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
