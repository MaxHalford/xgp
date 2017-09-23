package main

import "C"
import (
	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
)

var est xgp.Estimator

//export Fit
func Fit(
	X [][]float64,
	y []float64,
	metricName string,
	generations int,
	tuningGenerations int,
) {
	// Set parameters
	//var metric, err = metrics.GetMetric(metricName, 1)
	//fmt.Println(metric, err)
	//fmt.Println(generations, tuningGenerations)
	est.Metric = metrics.MeanSquaredError{}
	est.Transform = xgp.Identity{}
	est.PVariable = 0.5
	est.NodeInitializer = xgp.RampedHaldAndHalfInitializer{
		MinHeight: 2,
		MaxHeight: 5,
		PLeaf:     0.5,
	}
	est.FunctionSet = map[int][]xgp.Operator{
		2: []xgp.Operator{
			xgp.Sum{},
			xgp.Difference{},
			xgp.Product{},
			xgp.Division{},
		},
	}
	est.Generations = generations
	est.TuningGenerations = tuningGenerations
	est.ProgressChan = make(chan float64, generations+tuningGenerations)

	// Set the initial GA
	est.GA = &gago.GA{
		NewGenome: est.NewProgram,
		NPops:     1,
		PopSize:   100,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}

	// Set the tuning GA
	est.TuningGA = &gago.GA{
		NewGenome: est.NewProgramTuner,
		NPops:     1,
		PopSize:   20,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}

	// Fit the Estimator
	//_ = est.Fit(X, y)
}

//export Predict
func Predict(X [][]float64) []float64 {
	var yPred, _ = est.Predict(X)
	return yPred
}

func main() {}
