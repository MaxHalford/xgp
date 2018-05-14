package meta

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
)

// AdaBoost implements adaptive boosting.
type AdaBoost struct {
	NPrograms           uint
	RowSampling         float64
	EvalMetric          metrics.Metric
	EarlyStoppingRounds uint
	RNG                 *rand.Rand
}

// Train Bagging.
func (ada AdaBoost) Train(
	estimator *xgp.Estimator,
	XTrain [][]float64,
	YTrain []float64,
	XVal [][]float64,
	YVal []float64,
	verbose bool,
) (Ensemble, error) {

	var (
		ensemble = Ensemble{
			Programs:       make([]xgp.Program, 0),
			Weights:        make([]float64, 0),
			UsedCols:       make([][]int, 0),
			Classification: estimator.LossMetric.Classification(),
		}
		n          = len(XTrain[0])                    // Number of rows in the training set
		k          = int(ada.RowSampling * float64(n)) // Number of rows to sample
		p          = len(XTrain)                       // Number of columns in the training set
		W          = make([]float64, n)                // Row weights
		evalScores = make([]float64, 0)
	)

	// Uniformly initialize the row weights
	for i := range W {
		W[i] = 1 / float64(n)
	}

	for i := uint(0); i < ada.NPrograms; i++ {
		// Sample the training set
		var (
			rows = randIntsWeighted(k, W, ada.RNG)
			cols = randIntsNoRep(p, p, ada.RNG)
		)
		// Train on the sample
		var program, err = estimator.Fit(
			subsetRowsFloat64Matrix(XTrain, rows),
			subsetFloat64Slice(YTrain, rows),
			nil,
			XVal,
			YVal,
			nil,
			verbose,
		)
		if err != nil {
			return ensemble, err
		}
		// Compute prediction error
		YPred, err := program.Predict(XTrain, false)
		if err != nil {
			return ensemble, err
		}
		var weight float64
		// For classification use the weights update rule from Freund & Schapire 1997
		if estimator.LossMetric.Classification() {
			accuracy, err := metrics.Accuracy{}.Apply(YTrain, YPred, W)
			if err != nil {
				return ensemble, err
			}
			alpha := 0.5 * math.Log(accuracy/(1-accuracy))
			for j := range W {
				if YTrain[j] == YPred[j] {
					W[j] *= math.Exp(-alpha)
				} else {
					W[j] *= math.Exp(alpha)
				}
			}
			weight = alpha
		}
		// Normalize the weights
		W = normalize(W)
		// Store the program, it's weight, and the numbers of the columns it used
		ensemble.Programs = append(ensemble.Programs, program)
		ensemble.Weights = append(ensemble.Weights, weight)
		ensemble.UsedCols = append(ensemble.UsedCols, cols)
		// Compute the evaluation score
		if ada.EvalMetric == nil || XVal == nil {
			continue
		}
		score, err := ensemble.score(XVal, YVal, ada.EvalMetric)
		if err != nil {
			return ensemble, err
		}
		if verbose {
			fmt.Printf("%s: %.5f\n", ada.EvalMetric.String(), score)
		}
		// Store the evaluation score
		evalScores = append(evalScores, score)
		// Check if early stopping should occur
		if ada.EarlyStoppingRounds == 0 || uint(len(evalScores)) <= ada.EarlyStoppingRounds {
			continue
		}
		if checkEarlyStop(evalScores, ada.EarlyStoppingRounds, i, ada.EvalMetric) {
			if verbose {
				fmt.Println("Early stopping")
			}
			break
		}
	}

	return ensemble, nil
}
