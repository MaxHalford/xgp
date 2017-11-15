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
	generations int,
	lossMetricName string,
	maxHeight int,
	minHeight int,
	parsimonyCoeff float64,
	pTerminal float64,
	pConstant float64,
	rounds int,
	seed int64,
	tuningGenerations int,
	verbose bool,
) *C.char {

	var estimator, err = xgp.NewEstimator(
		constMax,
		constMin,
		evalMetricName,
		funcsString,
		lossMetricName,
		maxHeight,
		minHeight,
		generations,
		parsimonyCoeff,
		pConstant,
		pTerminal,
		seed,
		tuningGenerations,
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
