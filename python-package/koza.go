package main

// #cgo pkg-config: python
// #define Py_LIMITED_API
// #include <Python.h>
// PyObject* Py_String(char *pystring);
import "C"
import (
	"github.com/MaxHalford/koza"
)

var estimator *koza.Estimator

// Fit an Estimator to a dataset.
//export Fit
func Fit(
	X [][]float64,
	Y []float64,
	XNames []string,
	constMax float64,
	constMin float64,
	evalMetricName string,
	funcsString string,
	lossMetricName string,
	maxHeight int,
	minHeight int,
	nGenerations int,
	nPops int,
	nRounds int,
	pConstant float64,
	pCrossover float64,
	pFull float64,
	pHoistMutation float64,
	pPointMutation float64,
	pSubTreeMutation float64,
	pTerminal float64,
	parsimonyCoeff float64,
	pointMutationRate float64,
	popSize int,
	seed int,
	tuningNGenerations int,
	verbose bool,
) *C.char {
	var err error
	estimator, err = koza.NewEstimator(
		constMax,
		constMin,
		evalMetricName,
		funcsString,
		nGenerations,
		lossMetricName,
		maxHeight,
		minHeight,
		nPops,
		pConstant,
		pCrossover,
		pFull,
		pHoistMutation,
		pPointMutation,
		pSubTreeMutation,
		pTerminal,
		parsimonyCoeff,
		pointMutationRate,
		popSize,
		int64(seed),
		tuningNGenerations,
	)

	// Fit the Estimator
	err = estimator.Fit(X, Y, XNames, verbose)
	if err != nil {
		panic(err)
	}

	// Get the best obtained Program
	best, err := estimator.BestProgram()
	if err != nil {
		panic(err)
	}

	return C.CString(best.String())
}

// Predict a dataset. The predicitions will be copied onto YPred, this way the
// results can be accessed on the Python side.
//export Predict
func Predict(X [][]float64, predictProba bool, YPred []float64) {
	var Y, err = estimator.Predict(X, predictProba)
	if err != nil {
		panic(err)
	}
	copy(YPred, Y)
}

func main() {}
