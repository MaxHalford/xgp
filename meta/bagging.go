package meta

import (
	"fmt"
	"math/rand"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
)

// Bagging implements boostrap aggregation. NPrograms will be trained at most.
// At each round a uniform sample is constructed and fed to a new Program.
// Early stopping occurs if the evaluation score hasn't increased in at least
// EarlyStoppingRounds rounds.
type Bagging struct {
	NPrograms           uint
	RowSampling         float64
	ColSampling         float64
	EvalMetric          metrics.Metric
	EarlyStoppingRounds uint
	RNG                 *rand.Rand
}

// Train Bagging.
func (bag Bagging) Train(
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
		evalScores = make([]float64, 0)
		n          = len(XTrain[0])                    // Number of rows in the training set
		k          = int(bag.RowSampling * float64(n)) // Number of rows to sample
		p          = len(XTrain)                       // Number of columns in the training set
		q          = int(bag.ColSampling * float64(p)) // Number of columns to sample
	)

	for i := uint(0); i < bag.NPrograms; i++ {
		// Sample the training set
		var (
			rows = randInts(k, n, bag.RNG)
			cols = randIntsNoRep(q, p, bag.RNG)
		)
		// Train on the sample
		var RoundXVal [][]float64
		if XVal != nil {
			RoundXVal = subsetColsFloat64Matrix(XVal, cols)
		}
		var program, err = estimator.Fit(
			subsetRowsFloat64Matrix(subsetColsFloat64Matrix(XTrain, cols), rows),
			subsetFloat64Slice(YTrain, rows),
			nil,
			RoundXVal,
			YVal,
			nil,
			verbose,
		)
		if err != nil {
			return ensemble, err
		}
		// Store the program, it's weight, and the numbers of the columns it used
		ensemble.Programs = append(ensemble.Programs, program)
		ensemble.Weights = append(ensemble.Weights, 1)
		ensemble.UsedCols = append(ensemble.UsedCols, cols)
		// Compute the evaluation score
		if bag.EvalMetric == nil || XVal == nil {
			continue
		}
		score, err := ensemble.score(XVal, YVal, bag.EvalMetric)
		if err != nil {
			return ensemble, err
		}
		if verbose {
			fmt.Printf("%s: %.5f\n", bag.EvalMetric.String(), score)
		}
		// Store the evaluation score
		evalScores = append(evalScores, score)
		// Check if early stopping should occur
		if bag.EarlyStoppingRounds == 0 || uint(len(evalScores)) <= bag.EarlyStoppingRounds {
			continue
		}
		if checkEarlyStop(evalScores, bag.EarlyStoppingRounds, i, bag.EvalMetric) {
			if verbose {
				fmt.Println("Early stopping")
			}
			break
		}
	}

	return ensemble, nil
}
