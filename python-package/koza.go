package main

// #cgo pkg-config: python
// #define Py_LIMITED_API
// #include <Python.h>
// PyObject* Py_String(char *pystring);
import "C"
import (
	"reflect"
	"unsafe"

	"github.com/MaxHalford/koza"
	"github.com/MaxHalford/koza/tree"
)

var estimator *koza.Estimator

// ArrayToSlice converts a C double array to a Go Slice.
//export ArrayToSlice
func ArrayToSlice(a *C.double, length int) *[]float64 {
	var sh = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(a)),
		Len:  length,
		Cap:  length,
	}
	return (*[]float64)(unsafe.Pointer(&sh))
}

// Fit an Estimator to a dataset
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
	pFull float64,
	pHoistMutation float64,
	pPointMutation float64,
	pSubTreeCrossover float64,
	pSubTreeMutation float64,
	pTerminal float64,
	parsimonyCoeff float64,
	pointMutationRate float64,
	popSize int,
	seed int,
	tuningNGenerations int,
	verbose bool,
) *C.char {

	estimator, err := koza.NewEstimator(
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
		pFull,
		pHoistMutation,
		pPointMutation,
		pSubTreeCrossover,
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

	// Get the best obtained program
	best, err := estimator.BestProgram()
	if err != nil {
		panic(err)
	}

	var numpy = tree.NumpyDisplay{}.Apply(best.Tree)
	return C.CString(numpy)

}

// Predict a dataset
//export Predict
func Predict(X [][]float64, YPred *[]float64, predictProba bool) {
	var Y, err = estimator.Predict(X, predictProba)
	if err != nil {
		panic(err)
	}
	for i, y := range Y {
		(*YPred)[i] = y
	}
}

func main() {}
