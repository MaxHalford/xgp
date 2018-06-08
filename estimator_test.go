package xgp

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestEstimator(t *testing.T) {
	var conf = NewConfigWithDefaults()
	conf.RNG = rand.New(rand.NewSource(42))
	conf.NIndividuals = 20
	conf.MinHeight = 1
	conf.MaxHeight = 3
	var est, err = conf.NewEstimator()
	if err != nil {
		fmt.Println(err)
	}
	var (
		X = [][]float64{
			[]float64{1, 2, 3},
			[]float64{4, 5, 6},
		}
		Y = []float64{5, 7, 9}
	)
	est.Fit(X, Y, nil, nil, nil, nil, false)
	var prog = est.BestProgram()
	fmt.Println(prog)
	var pred, _ = prog.PredictPartial([]float64{4, 7}, false)
	fmt.Println(pred)
	// Output:
	// x1+x0
	// 11
}
