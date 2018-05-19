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
	LearningRate        float64
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
		n          = len(XTrain[0])
		p          = len(XTrain)        // Number of columns in the training set
		W          = make([]float64, n) // Row weights
		evalScores = make([]float64, 0)
	)

	// Uniformly initialize the row weights
	for i := range W {
		W[i] = 1 / float64(n)
	}

	for i := uint(0); i < ada.NPrograms; i++ {
		fit := func(X [][]float64, Y []float64) (xgp.Program, error) {
			return estimator.Fit(X, Y, nil, XVal, YVal, nil, verbose)
		}
		var (
			prog   xgp.Program
			weight float64
			err    error
		)
		if estimator.LossMetric.Classification() {
			prog, W, weight, err = ada.boostSAMMER(fit, XTrain, YTrain, W, ada.RNG)
		} else {
			prog, W, weight, err = ada.boostR2(fit, XTrain, YTrain, W, ada.RNG)
		}
		if err != nil {
			return ensemble, err
		}
		if W == nil || weight == 0 {
			continue
		}
		// Normalize the weights so that they sum to 1
		W = normalize(W)
		// Store the program, it's weight, and the numbers of the columns it used
		ensemble.Programs = append(ensemble.Programs, prog)
		ensemble.Weights = append(ensemble.Weights, weight)
		ensemble.UsedCols = append(ensemble.UsedCols, randIntsNoRep(p, p, ada.RNG))
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

func (ada AdaBoost) boostR2(
	fit func(X [][]float64, Y []float64) (xgp.Program, error),
	X [][]float64,
	Y []float64,
	W []float64,
	rng *rand.Rand,
) (xgp.Program, []float64, float64, error) {
	// Sample the data and obtain an estimator
	var (
		k            = int(ada.RowSampling * float64(len(W)))
		rows         = randIntsWeighted(k, W, ada.RNG)
		program, err = fit(
			subsetRowsFloat64Matrix(X, rows),
			subsetFloat64Slice(Y, rows),
		)
	)
	if err != nil {
		return program, nil, math.NaN(), err
	}
	// Make predictions
	YPred, err := program.Predict(X, false)
	if err != nil {
		return program, nil, math.NaN(), err
	}
	// Compute individual errors
	var (
		errors   = make([]float64, len(W))
		maxError float64
	)
	for i := range errors {
		e := math.Abs(Y[i] - YPred[i])
		errors[i] = e
		if e > maxError {
			maxError = e
		}
	}
	if maxError > 0 {
		for i := range errors {
			errors[i] /= maxError
		}
	}
	for i, e := range errors {
		errors[i] = 1 - math.Exp(-e)
	}
	// Compute average error
	var meanError float64
	for i, w := range W {
		meanError += w * errors[i]
	}
	// Check error
	if meanError == 0 {
		return program, W, 1, nil
	}
	if meanError >= 0.5 {
		return program, nil, 0, nil
	}
	// Determine program weight
	beta := meanError / (1 - meanError)
	weight := ada.LearningRate * math.Log(1/beta)
	// Update sample weights
	for i, e := range errors {
		W[i] *= math.Pow(beta, (1-e)*ada.LearningRate)
	}
	return program, W, weight, nil
}

func (ada AdaBoost) boostSAMMER(
	fit func(X [][]float64, Y []float64) (xgp.Program, error),
	X [][]float64,
	Y []float64,
	W []float64,
	rng *rand.Rand,
) (xgp.Program, []float64, float64, error) {
	// Sample the data and obtain an estimator
	var (
		k            = int(ada.RowSampling * float64(len(W)))
		rows         = randIntsWeighted(k, W, ada.RNG)
		program, err = fit(
			subsetRowsFloat64Matrix(X, rows),
			subsetFloat64Slice(Y, rows),
		)
	)
	if err != nil {
		return program, nil, math.NaN(), err
	}
	// Make predictions
	YPred, err := program.Predict(X, true)
	if err != nil {
		return program, nil, math.NaN(), err
	}
	// Compute error
	var e float64
	for i, y := range Y {
		d := y - YPred[i]
		if math.Abs(d) > 0.5 {
			e += math.Log(math.Max(YPred[i], 0.00001))
		} else {
			e -= math.Log(math.Max(YPred[i], 0.00001))
		}
	}
	// Determine program weight
	weight := ada.LearningRate * 0.5 * e
	// Update sample weights
	for i, w := range W {
		if weight > 0 || w < 0 {
			W[i] *= math.Exp(w)
		}
	}
	return program, W, weight, nil
}
