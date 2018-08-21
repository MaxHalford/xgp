package metrics

import "math"

// R2 measures the coefficient of determination.
type R2 struct{}

// Apply R2.
func (r2 R2) Apply(yTrue, yPred, weights []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(yPred)}
	}
	if weights != nil && len(yTrue) != len(weights) {
		return math.Inf(1), &errMismatchedLengths{len(yTrue), len(weights)}
	}

	// Compute the mean of the observed data
	var yMean float64
	if weights != nil {
		var ws float64
		for i, y := range yTrue {
			yMean += y * weights[i]
			ws += weights[i]
		}
		yMean /= ws
	} else {
		for _, y := range yTrue {
			yMean += y
		}
		yMean /= float64(len(yTrue))
	}

	var (
		SSR float64
		SST float64
	)
	if weights != nil {
		for i := range yTrue {
			SSR += math.Pow(yPred[i]-yTrue[i], 2) * weights[i]
			SST += math.Pow(yTrue[i]-yMean, 2) * weights[i]
		}
		return 1 - SSR/SST, nil
	}
	for i := range yTrue {
		SSR += math.Pow(yPred[i]-yTrue[i], 2)
		SST += math.Pow(yTrue[i]-yMean, 2)
	}
	return 1 - SSR/SST, nil
}

// Classification method of R2.
func (r2 R2) Classification() bool {
	return false
}

// BiggerIsBetter method of R2.
func (r2 R2) BiggerIsBetter() bool {
	return false
}

// NeedsProbabilities method of R2.
func (r2 R2) NeedsProbabilities() bool {
	return false
}

// String method of R2.
func (r2 R2) String() string {
	return "r2"
}
