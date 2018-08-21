package metrics

import "math"

// RMSE measures the root mean squared error (RMSE).
type RMSE struct{}

// Apply RMSE.
func (rmse RMSE) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var mse, err = MSE{}.Apply(yTrue, yPred, weights)
	if err != nil {
		return math.Inf(1), err
	}
	return math.Pow(mse, 0.5), nil
}

// Classification method of RMSE.
func (rmse RMSE) Classification() bool {
	return false
}

// BiggerIsBetter method of RMSE.
func (rmse RMSE) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of RMSE.
func (rmse RMSE) NeedsProbabilities() bool {
	return false
}

// String method of RMSE.
func (rmse RMSE) String() string {
	return "rmse"
}
