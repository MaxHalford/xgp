package ensemble

import (
	"fmt"
	"math/rand"

	"github.com/MaxHalford/koza/metrics"
)

func sample(
	X [][]float64,
	Y []float64,
	W []float64,
	rowSampling float64,
	colSampling float64,
	boostrapRows bool,
	bootstrapCols bool,
	rng *rand.Rand,
) ([][]float64, []float64, []float64, []int, []int, error) {
	// Sample row indexes
	var n = int(rowSampling * float64(len(X[0])))
	rowIdxs, err := randomInts(n, 0, len(X[0]), boostrapRows, rng)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	// Sample column indexes
	var p = int(colSampling * float64(len(X)))
	colIdxs, err := randomInts(p, 0, len(X), bootstrapCols, rng)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	// Create the sample
	var (
		XSam = make([][]float64, len(colIdxs))
		YSam = make([]float64, len(rowIdxs))
		WSam []float64
	)
	if W != nil {
		WSam = make([]float64, len(rowIdxs))
	}
	for i := range colIdxs {
		XSam[i] = make([]float64, len(rowIdxs))
	}
	for i, r := range rowIdxs {
		for j, c := range colIdxs {
			XSam[j][i] = X[c][r]
		}
		YSam[i] = Y[r]
		if W != nil {
			WSam[i] = W[r]
		}
	}
	if W != nil {
		return XSam, YSam, WSam, rowIdxs, colIdxs, nil
	}
	return XSam, YSam, nil, rowIdxs, colIdxs, nil
}

// BaggingRegressor implements boostrap aggregation for regression tasks.
type BaggingRegressor struct {
	NEstimators   int     `json:"-"`
	RowSampling   float64 `json:"-"` // Percentage of rows sampled for each estimator
	ColSampling   float64 `json:"-"` // Percentage of columns sampled for each estimator
	BootstrapRows bool    `json:"-"` // Whether to sample columns with replacement or not
	BootstrapCols bool    `json:"-"` // Whether to sample columns with replacement or not

	Predictors    []Predictor `json:"predictors"`
	PredictorCols [][]int     `json:"predictor_cols"`
}

func (bag *BaggingRegressor) Fit(
	learner Learner,
	XTrain [][]float64,
	YTrain []float64,
	WTrain []float64,
	XVal [][]float64,
	YVal []float64,
	WVal []float64,
	verbose bool,
	evalMetric metrics.Metric,
	rng *rand.Rand,
) error {

	bag.Predictors = make([]Predictor, 0)
	bag.PredictorCols = make([][]int, 0)

	for i := 0; i < bag.NEstimators; i++ {
		// Sample
		var XSam, YSam, WSam, _, colIdxs, err = sample(
			XTrain,
			YTrain,
			WTrain,
			bag.RowSampling,
			bag.ColSampling,
			bag.BootstrapRows,
			bag.BootstrapCols,
			rng,
		)
		if err != nil {
			return err
		}
		// Train on the sample
		predictor, err := learner.Fit(
			XSam,
			YSam,
			WSam,
			XVal,
			YVal,
			WVal,
			verbose,
		)
		if err != nil {
			return err
		}
		// Store the resulting predictor and the associated columns
		bag.Predictors = append(bag.Predictors, predictor)
		bag.PredictorCols = append(bag.PredictorCols, colIdxs)
		// Calculate the current scores
		if !verbose || evalMetric == nil {
			continue
		}
		YTrainPred, err := bag.Predict(XTrain, false)
		if err != nil {
			return err
		}
		trainScore, err := evalMetric.Apply(YTrain, YTrainPred, WTrain)
		if err != nil {
			return err
		}
		fmt.Printf("Bagging train %s: %.5f", evalMetric, trainScore)
		if XVal != nil && YVal != nil {
			YValPred, err := bag.Predict(XVal, false)
			if err != nil {
				return err
			}
			valScore, err := evalMetric.Apply(YVal, YValPred, WVal)
			if err != nil {
				return err
			}
			fmt.Printf(", val %s: %.5f\n", evalMetric, valScore)
		}
	}

	return nil
}

func (bag BaggingRegressor) Predict(X [][]float64, predictProba bool) ([]float64, error) {

	var (
		Y       = make([]float64, len(X[0]))
		rowPred = make([]float64, len(bag.Predictors))
	)

	// Iterate over each feature vector
	for i := range X[0] {
		// Get the individual predictions
		for j, predictor := range bag.Predictors {
			var x = make([]float64, len(X))
			for k, c := range bag.PredictorCols[j] {
				x[k] = X[c][i]
			}
			var pred, err = predictor.PredictRow(x, false)
			if err != nil {
				return nil, err
			}
			rowPred[j] = pred
		}
		// Aggregate the individual predictions
		Y[i] = meanFloat64s(rowPred)
	}

	return Y, nil
}
