package metrics

import "math"

// RootMeanSquaredError measures the root mean squared error (RMSE).
type RootMeanSquaredError struct{}

// Apply RootMeanSquaredError.
func (metric RootMeanSquaredError) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var mse, err = MeanSquaredError{}.Apply(yTrue, yPred, weights)
	if err != nil {
		return math.Inf(1), err
	}
	return math.Pow(mse, 0.5), nil
}

// Classification method of RootMeanSquaredError.
func (metric RootMeanSquaredError) Classification() bool {
	return false
}

// BiggerIsBetter method of RootMeanSquaredError.
func (metric RootMeanSquaredError) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of RootMeanSquaredError.
func (metric RootMeanSquaredError) NeedsProbabilities() bool {
	return false
}

// String method of RootMeanSquaredError.
func (metric RootMeanSquaredError) String() string {
	return "rmse"
}
