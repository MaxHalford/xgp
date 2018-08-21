package meta

import (
	"fmt"
	"math"
	"time"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
)

// GradientBoosting implements gradient boosting on top of genetic programming.
type GradientBoosting struct {
	xgp.GPConfig         `json:"-"`
	NRounds              uint               `json:"-"`
	NEarlyStoppingRounds uint               `json:"-"`
	LearningRate         float64            `json:"-"`
	Loss                 metrics.DiffMetric `json:"-"`
	Programs             []xgp.Program      `json:"programs"`
	ValScores            []float64          `json:"val_scores"`
	TrainScores          []float64          `json:"train_scores"`
	YMean                float64            `json:"y_mean"` // Set during Fit
}

// NewGradientBoosting returns a GradientBoosting.
func NewGradientBoosting(conf xgp.GPConfig, n, k uint, lr float64, loss metrics.DiffMetric) (*GradientBoosting, error) {
	return &GradientBoosting{
		GPConfig:             conf,
		NRounds:              n,
		NEarlyStoppingRounds: k,
		LearningRate:         lr,
		Loss:                 loss,
		Programs:             make([]xgp.Program, 0),
		ValScores:            make([]float64, 0),
	}, nil
}

// Logistic function.
func expit(y []float64) []float64 {
	var p = make([]float64, len(y))
	for i, yi := range y {
		p[i] = 1 / (1 + math.Exp(-yi))
	}
	return p
}

func classify(y []float64) []float64 {
	var classes = make([]float64, len(y))
	for i, yi := range y {
		if yi > 0.5 {
			classes[i] = 1
		}
	}
	return classes
}

// Fit iteratively trains a GP on the gradient of the loss.
func (gb *GradientBoosting) Fit(
	// Required arguments
	X [][]float64,
	Y []float64,
	// Optional arguments (can safely be nil)
	W []float64,
	XVal [][]float64,
	YVal []float64,
	WVal []float64,
	verbose bool,
) error {
	var start = time.Now()

	// We use symbolic regression even if the task is classication
	if gb.Loss.Classification() {
		gb.LossMetric = metrics.MSE{}
	}

	// Start from the target mean
	gb.YMean = mean(Y)
	var YPred = make([]float64, len(Y))
	for i := range YPred {
		YPred[i] = gb.YMean
	}

	// Store the best validation score in order to check for early stopping
	var (
		bestVal          = math.Inf(1)
		earlyStopCounter = gb.NEarlyStoppingRounds
	)

	for i := uint(0); i < gb.NRounds; i++ {
		// Compute the gradient
		var (
			grad []float64
			err  error
		)
		if gb.Loss.Classification() {
			grad, err = gb.Loss.Gradient(Y, expit(YPred))
		} else {
			grad, err = gb.Loss.Gradient(Y, YPred)
		}
		if err != nil {
			return err
		}

		// Train a GP on the gradient
		gp, err := gb.NewGP()
		if err != nil {
			return err
		}
		err = gp.Fit(X, grad, W, nil, nil, nil, false)
		if err != nil {
			return err
		}

		// Extract the best obtained Program
		prog, err := gp.BestProgram()
		if err != nil {
			return err
		}
		gb.Programs = append(gb.Programs, prog)

		// Make predictions
		update, err := prog.Predict(X, false)
		if err != nil {
			return err
		}
		for i, u := range update {
			YPred[i] -= gb.LearningRate * u
		}

		// Compute training score
		var trainScore float64
		if gb.EvalMetric.Classification() {
			if gb.EvalMetric.NeedsProbabilities() {
				trainScore, err = gb.EvalMetric.Apply(Y, expit(YPred), nil)
			} else {
				trainScore, err = gb.EvalMetric.Apply(Y, classify(expit(YPred)), nil)
			}
		} else {
			trainScore, err = gb.EvalMetric.Apply(Y, YPred, nil)
		}
		if err != nil {
			return err
		}
		gb.TrainScores = append(gb.TrainScores, trainScore)

		// If there is no validation set then stop
		if XVal == nil || YVal == nil || gb.EvalMetric == nil {
			if verbose {
				fmt.Printf(
					"%s -- train %s: %.5f -- round %d\n",
					fmtDuration(time.Since(start)),
					gb.EvalMetric.String(),
					trainScore,
					i+1,
				)
			}
			continue
		}

		// Compute validation score
		YValPred, err := gb.Predict(XVal, gb.EvalMetric.NeedsProbabilities())
		if err != nil {
			return err
		}
		valScore, err := gb.EvalMetric.Apply(YVal, YValPred, nil)
		if err != nil {
			return err
		}
		gb.ValScores = append(gb.ValScores, valScore)

		// Display progress
		if verbose {
			fmt.Printf(
				"%s -- train %s: %.5f -- val %s: %.5f -- round %d\n",
				fmtDuration(time.Since(start)),
				gb.EvalMetric.String(),
				trainScore,
				gb.EvalMetric.String(),
				valScore,
				i+1,
			)
		}

		// Check for early stopping
		if valScore < bestVal {
			earlyStopCounter = gb.NEarlyStoppingRounds
			bestVal = valScore
		} else {
			earlyStopCounter--
		}
		if earlyStopCounter == 0 {
			if verbose {
				fmt.Println("Early stopping")
			}
			break
		}
	}
	return nil
}

// Predict accumulates the predictions of each stored Program.
func (gb GradientBoosting) Predict(X [][]float64, proba bool) ([]float64, error) {
	// Start from the target mean
	var YPred = make([]float64, len(X[0]))
	for i := range YPred {
		YPred[i] = gb.YMean
	}
	// Accumulate predictions
	for _, prog := range gb.Programs {
		update, err := prog.Predict(X, false)
		if err != nil {
			return nil, err
		}
		for i, u := range update {
			YPred[i] -= gb.LearningRate * u
		}
	}
	// Transform in case of classification
	if gb.Loss.Classification() {
		YPred = expit(YPred)
		if !proba {
			YPred = classify(YPred)
		}
	}
	return YPred, nil
}
