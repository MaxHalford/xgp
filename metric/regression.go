package metric

import (
	"math"
)

// MeanAbsoluteError measures the mean absolute error (MAE).
type MeanAbsoluteError struct{}

// Apply MeanAbsoluteError.
func (metric MeanAbsoluteError) Apply(yTrue []float64, yPred []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}

	var sum float64
	for i := range yTrue {
		sum += math.Abs(yTrue[i] - yPred[i])
	}
	return sum / float64(len(yTrue)), nil
}

// MeanSquaredError measures the mean absolute error (MAE).
type MeanSquaredError struct{}

// Apply MeanSquaredError.
func (metric MeanSquaredError) Apply(yTrue []float64, yPred []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}

	var sum float64
	for i := range yTrue {
		sum += math.Pow(yTrue[i]-yPred[i], 2)
	}
	return sum / float64(len(yTrue)), nil
}

// R2 measures the coefficient of determination.
type R2 struct{}

// Apply R2.
func (metric R2) Apply(yTrue []float64, yPred []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}

	// Compute the mean of the observed data
	var yMean float64
	for _, y := range yTrue {
		yMean += y
	}
	yMean /= float64(len(yTrue))

	var (
		SST float64
		SSR float64
	)
	for i := range yTrue {
		SST += math.Pow(yTrue[i]-yMean, 2)
		SSR += math.Pow(yPred[i]-yTrue[i], 2)
	}

	return 1 - SSR/SST, nil
}
