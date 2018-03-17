package ensemble

import (
	"math/rand"
)

// BaggingRegressor implements boostrap aggregation for regression tasks.
type BaggingRegressor struct {
	NEstimators int        `json:"-"`
	RowSampling float64    `json:"-"` // Percentage of rows sampled for each estimator
	ColSampling float64    `json:"-"` // Percentage of columns sampled for each estimator
	RNG         *rand.Rand `json:"-"`

	Predictors    []Predictor `json:"predictors"`
	PredictorCols [][]int     `json:"predictor_cols"`
}

func (bag *BaggingRegressor) Fit(learner Learner, X [][]float64, Y []float64, W []float64, verbose bool) error {

	bag.Predictors = make([]Predictor, bag.NEstimators)
	bag.PredictorCols = make([][]int, bag.NEstimators)

	// If no weights are provided then uniform weights are used
	if W == nil {
		W = make([]float64, len(X[0]))
		for i := range W {
			W[i]++
		}
	}

	// Determine how many rows and columns to sample
	var (
		n = int(bag.RowSampling * float64(len(X[0])))
		p = int(bag.ColSampling * float64(len(X)))
	)

	for i := 0; i < bag.NEstimators; i++ {
		// Sample row and columns indices
		var (
			rowIdxs = sampleIndices(n, W, bag.RNG)
			colIdxs = sampleIndices(p, W, bag.RNG)
		)
		// Train on the sample
		var predictor, err = learner.Fit(
			subsetFloat64Matrix(X, rowIdxs, colIdxs),
			subsetFloat64Slice(Y, rowIdxs),
			nil,
			verbose,
		)
		if err != nil {
			return err
		}
		// Store the resulting predictor and the associated columns
		bag.Predictors = append(bag.Predictors, predictor)
		bag.PredictorCols = append(bag.PredictorCols, colIdxs)
	}

	return nil
}

func (bag BaggingRegressor) Predict(X [][]float64, predictProba bool) ([]float64, error) {

	var (
		YAll = make([]float64, len(bag.Predictors))
		Y    = make([]float64, len(X[0]))
	)

	// Iterate over each feature vector
	for i := range X[0] {
		// Collect each predictor's output
		for j, predictor := range bag.Predictors {
			var x = make([]float64, len(X))
			for k, c := range bag.PredictorCols[j] {
				x[k] = X[c][i]
			}
			var pred, err = predictor.PredictPartial(x, false)
			if err != nil {
				return nil, err
			}
			YAll[j] = pred
		}
		// Aggregate the individual predictions
		Y[i] = meanFloat64s(YAll)
	}

	return Y, nil
}
