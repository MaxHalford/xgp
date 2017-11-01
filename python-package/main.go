package main

import "C"
import (
	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/tree"
)

var est xgp.Estimator

//export Fit
func Fit(
	X [][]float64,
	Y []float64,
	XNames []string,
	constMax float64,
	constMin float64,
	evalMetricName string,
	funcsString string,
	generations int,
	lossMetricName string,
	maxHeight int,
	minHeight int,
	parsimonyCoeff float64,
	pTerminal float64,
	pConstant float64,
	rounds int,
	tuningGenerations int,
	verbose bool,
) {
	// Set the constant minimum and maximum values
	est.ConstMin = constMin
	est.ConstMax = constMax

	// Set evaluation metric
	var evalMetric, err = metrics.GetMetric(evalMetricName, 1)
	if err != nil {
		panic(err)
	}
	est.EvalMetric = evalMetric

	// Set the functions to use
	est.Functions = []tree.Operator{tree.Sum{}, tree.Difference{}, tree.Product{}, tree.Division{}}

	// Set the number of generations
	est.Generations = generations

	// Set loss metric
	lossMetric, err := metrics.GetMetric(lossMetricName, 1)
	if err != nil {
		panic(err)
	}
	est.LossMetric = lossMetric

	// Set the tree initializer
	est.TreeInitializer = tree.RampedHaldAndHalfInitializer{
		MinHeight: minHeight,
		MaxHeight: maxHeight,
		PTerminal: pTerminal,
	}

	// Set the parsimony coefficient
	est.ParsimonyCoeff = parsimonyCoeff

	// Set the probability of producing a Constant when producing a terminal branch
	est.PConstant = pConstant

	// Set the number of tuning generations
	est.TuningGenerations = tuningGenerations

	// Set the initial GA
	est.GA = &gago.GA{
		NewGenome: est.NewProgram,
		NPops:     1,
		PopSize:   1000,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}

	// Fit the Estimator
	err = est.Fit(X, Y, XNames, verbose)
	if err != nil {
		panic(err)
	}
}

//export Predict
func Predict(X [][]float64, predictProba bool) []float64 {
	var yPred, err = est.Predict(X, predictProba)
	if err != nil {
		panic(err)
	}
	return yPred
}

func main() {}
