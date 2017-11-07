package main

// #cgo pkg-config: python
// #define Py_LIMITED_API
// #include <Python.h>
// PyObject* Py_String(char *pystring);
import "C"
import (
	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/tree"
)

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
) *C.char {

	// Determine the fitness and evaluation metrics to use
	lossMetric, err := metrics.GetMetric(lossMetricName, 1)
	if err != nil {
		panic(err)
	}

	// Default the evaluation metric to the fitness metric if it's nil
	var evalMetric metrics.Metric
	if evalMetricName == "" {
		evalMetric = lossMetric
	} else {
		metric, err := metrics.GetMetric(evalMetricName, 1)
		if err != nil {
			panic(err)
		}
		evalMetric = metric
	}

	// The convention is to use a fitness metric which has to be minimized
	if lossMetric.BiggerIsBetter() {
		lossMetric = metrics.NegativeMetric{Metric: lossMetric}
	}

	// Determine the functions to use
	functions, err := tree.ParseStringFuncs(funcsString)
	if err != nil {
		panic(err)
	}

	var estimator = xgp.Estimator{
		EvalMetric: evalMetric,
		LossMetric: lossMetric,
		ConstMin:   constMin,
		ConstMax:   constMax,
		PConstant:  pConstant,
		TreeInitializer: tree.RampedHaldAndHalfInitializer{
			MinHeight: minHeight,
			MaxHeight: maxHeight,
			PTerminal: pTerminal,
		},
		Functions:         functions,
		Generations:       generations,
		TuningGenerations: tuningGenerations,
		ParsimonyCoeff:    parsimonyCoeff,
	}

	// Set the initial GA
	estimator.GA = &gago.GA{
		NewGenome: estimator.NewProgram,
		NPops:     1,
		PopSize:   100,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}

	// Fit the Estimator
	err = estimator.Fit(X, Y, XNames, verbose)
	if err != nil {
		panic(err)
	}

	// Get the best obtained program
	best, err := estimator.BestProgram()
	if err != nil {
		panic(err)
	}

	var numpy = tree.NumpyDisplay{}.Apply(best.Tree)
	return C.CString(numpy)

}

func main() {}
