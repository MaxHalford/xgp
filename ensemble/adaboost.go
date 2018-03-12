package ensemble

import (
	"fmt"
	"math/rand"
)

type AdaBoostRegressor struct {
	NEstimators      int         `json:"-"`
	RNG              *rand.Rand  `json:"-"`
	Predictors       []Predictor `json:"predictors"`
	PredictorErrors  []float64   `json:"predictor_errors"`
	PredictorWeights []float64   `json:"predictor_weights"`
}

func (ada *AdaBoostRegressor) Fit(learner Learner, X [][]float64, Y []float64, W []float64, verbose bool) error {

	ada.Predictors = make([]Predictor, ada.NEstimators)
	ada.PredictorErrors = make([]float64, ada.NEstimators)
	ada.PredictorWeights = make([]float64, ada.NEstimators)

	// If no weights are provided then uniform weights are used
	if W == nil {
		W = make([]float64, len(X[0]))
		for i := range W {
			W[i]++
		}
	}

	// Normalize the weights
	W = normalizeFloat64s(W)

	for i := 0; i < ada.NEstimators; i++ {
		// Sample row and columns indices
		var (
			rowIdxs = sampleIndices(len(X[0]), W, ada.RNG)
			colIdxs = sampleIndices(len(X), W, ada.RNG)
		)
		// Train on the sample
		var predictor, err = learner.Fit(
			subsetFloat64Matrix(X, rowIdxs, colIdxs),
			subsetFloat64Slice(Y, rowIdxs),
			nil,
			verbose,
		)
		fmt.Println(predictor, err)
	}

	return nil

}
