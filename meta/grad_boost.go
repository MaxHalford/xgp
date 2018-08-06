package meta

import (
	"fmt"

	"github.com/MaxHalford/xgp"
)

type GradBoost struct {
	xgp.GPConfig
	NRounds              uint
	NEarlyStoppingRounds uint
	LearningRate         float64
	Loss                 Loss
	Programs             []xgp.Program
	ValScores            []float64
}

func NewGradBoost(conf xgp.GPConfig, n, k uint, lr float64, loss Loss) (*GradBoost, error) {
	return &GradBoost{
		GPConfig:             conf,
		NRounds:              n,
		NEarlyStoppingRounds: k,
		LearningRate:         lr,
		Loss:                 loss,
		Programs:             make([]xgp.Program, 0),
		ValScores:            make([]float64, 0),
	}, nil
}

func (gb *GradBoost) Fit(X [][]float64, y []float64, XVal [][]float64, yVal []float64) error {
	var (
		yMean = mean(y)
		yPred = make([]float64, len(y))
	)
	for i := range yPred {
		yPred[i] = yMean
	}
	for i := uint(0); i < gb.NRounds; i++ {
		// Compute the gradient
		grad := gb.Loss.GradEval(y, yPred)
		// Initialize a GP
		gp, err := gb.NewGP()
		if err != nil {
			return err
		}
		// Train the GP on the gradient to find a good Program
		prog, err := gp.Fit(X, grad, nil, nil, nil, nil, false)
		if err != nil {
			return err
		}
		update, err := prog.Predict(X, gb.Loss.Classification())
		if err != nil {
			return err
		}
		for i, u := range update {
			yPred[i] -= u * gb.LearningRate
		}
		// Store the Program
		gb.Programs = append(gb.Programs, prog)
		// Compute validation score
		if XVal != nil && yVal != nil && gb.EvalMetric != nil {
			yValPred, err := gb.Predict(XVal)
			if err != nil {
				return err
			}
			score, err := gb.EvalMetric.Apply(yVal, yValPred, nil)
			if err != nil {
				return err
			}
			fmt.Printf("%s: %.5f\n", gb.EvalMetric.String(), score)
			// Store the validation score
			gb.ValScores = append(gb.ValScores, score)
			// Check for early stopping
			if uint(len(gb.Programs)) > gb.NEarlyStoppingRounds &&
				checkEarlyStop(gb.ValScores, gb.EvalMetric, i, gb.NEarlyStoppingRounds) {
				fmt.Println("Early stopping")
				break
			}
		}
	}
	return nil
}

func (gb GradBoost) Predict(X [][]float64) ([]float64, error) {
	var yPred = make([]float64, len(X[0]))
	for _, prog := range gb.Programs {
		update, err := prog.Predict(X, gb.Loss.Classification())
		if err != nil {
			return nil, err
		}
		for i, u := range update {
			yPred[i] -= u * gb.LearningRate
		}
	}
	return yPred, nil
}
