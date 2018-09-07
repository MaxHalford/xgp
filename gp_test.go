package xgp

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGP(t *testing.T) {
	// Define training variables
	var conf = NewDefaultGPConfig()
	conf.RNG = rand.New(rand.NewSource(42))
	conf.NIndividuals = 30
	conf.MinHeight = 1
	conf.MaxHeight = 3
	conf.PolishBest = false

	// Initialize a GP
	var gp, err = conf.NewGP()
	if err != nil {
		fmt.Println(err)
	}

	// Define the training set
	var (
		X = [][]float64{
			[]float64{1, 2, 3},
			[]float64{4, 5, 6},
		}
		Y = []float64{5, 7, 9}
	)

	// Train the GP
	gp.Fit(X, Y, nil, nil, nil, nil, false)
	var best, _ = gp.BestProgram()
	fmt.Println(best)

	// Make predictions
	var pred, _ = gp.PredictPartial([]float64{4, 7}, false)
	fmt.Println(pred)

	// Output:
	// x1+x0
	// 11
}

func TestGPProgress(t *testing.T) {
	var conf = NewDefaultGPConfig()
	conf.RNG = rand.New(rand.NewSource(42))
	conf.NIndividuals = 30
	conf.MinHeight = 1
	conf.MaxHeight = 3
	conf.PolishBest = false
	var gp, err = conf.NewGP()
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
	gp.Fit(X, Y, nil, nil, nil, nil, false)
	var (
		progress = gp.progress(time.Now())
		expected = "00:00:00, train mse: 0.00000"
	)
	if progress != expected {
		t.Errorf("Expected %s, got %s", expected, progress)
	}
}
