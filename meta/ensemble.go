package meta

import (
	"math"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/metrics"
)

type Ensemble struct {
	Programs       []xgp.Program `json:"programs"`
	Weights        []float64     `json:"weights"`
	UsedCols       [][]int       `json:"used_columns"`
	Classification bool          `json:"is_classification"`
}

func (ensemble Ensemble) Predict(X [][]float64, proba bool) ([]float64, error) {
	var Y = make([]float64, len(X[0]))
	for i := range Y {
		var y = make([]float64, len(ensemble.Programs))
		// Collect the predictions of each Program
		for j, prog := range ensemble.Programs {
			// Build x with the right columns
			var x = make([]float64, len(ensemble.UsedCols[j]))
			for k, c := range ensemble.UsedCols[j] {
				x[k] = X[c][i]
			}
			yj, err := prog.PredictPartial(x, proba)
			if err != nil {
				return nil, err
			}
			y[j] = yj
		}
		// Aggregate predictions
		if ensemble.Classification && !proba {
			Y[i] = weightedVote(y, ensemble.Weights)
		} else {
			Y[i] = weightedMean(y, ensemble.Weights)
		}
	}
	return Y, nil
}

func (ensemble Ensemble) score(X [][]float64, YTrue []float64, metric metrics.Metric) (float64, error) {
	YPred, err := ensemble.Predict(X, metric.NeedsProbabilities())
	if err != nil {
		return math.NaN(), err
	}
	score, err := metric.Apply(YTrue, YPred, nil)
	if err != nil {
		return math.NaN(), err
	}
	return score, nil
}

func weightedMean(x, w []float64) float64 {
	var (
		xs float64
		ws float64
	)
	for i, wi := range w {
		xs += wi * x[i]
		ws += wi
	}
	return xs / ws
}

func weightedVote(x, w []float64) float64 {
	// Start by counting the votes
	var votes = make(map[float64]float64)
	for i, xi := range x {
		votes[xi] += w[i]
	}
	// Then find the class with the highest vote
	var (
		c float64
		v float64
	)
	for class, vote := range votes {
		if vote > v {
			c = class
			v = vote
		}
	}
	return c
}
