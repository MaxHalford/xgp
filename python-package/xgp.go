package main

// #cgo pkg-config: python
// #define Py_LIMITED_API
// #include <Python.h>
// PyObject* Py_String(char *pystring);
import "C"
import (
	"github.com/MaxHalford/xgp"
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
	lossMetricName string,
	maxHeight int,
	minHeight int,
	nGenerations int,
	nPops int,
	nRounds int,
	parsimonyCoeff float64,
	pConstant float64,
	pHoistMutation float64,
	pPointMutation float64,
	pSubTreeCrossover float64,
	pSubTreeMutation float64,
	pTerminal float64,
	popSize int,
	seed int,
	tuningNGenerations int,
	verbose bool,
) *C.char {

	var estimator, err = xgp.NewEstimator(
		constMax,
		constMin,
		evalMetricName,
		funcsString,
		nGenerations,
		lossMetricName,
		maxHeight,
		minHeight,
		nPops,
		parsimonyCoeff,
		pConstant,
		pHoistMutation,
		pPointMutation,
		pSubTreeCrossover,
		pSubTreeMutation,
		pTerminal,
		popSize,
		int64(seed),
		tuningNGenerations,
	)

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
