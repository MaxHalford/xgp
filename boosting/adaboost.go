package boosting

import (
	"math"

	"github.com/MaxHalford/xgp/dataframe"
)

// AdaBoost implements adaptive boosting. The implementation is based on the
// SAMME algorithm proposed in Zhu et al (2009).
type AdaBoost struct {
	Regression       bool
	rowWeights       []float64
	predictors       []Predictor
	predictorWeights []float64
	nClasses         int
}

// fit represents one round of boosting.
func (ada AdaBoost) fit(learner Learner, df *dataframe.DataFrame) (Predictor, float64, error) {
	// Fit the weak learner
	var (
		predictor, err = learner.Learn(df)
		yPred          = make([]float64, df.NRows())
	)
	if err != nil {
		return nil, 0, err
	}
	// Compute weighted error
	var E, W float64
	for i, x := range df.X {
		y, err := predictor.PredictRow(x)
		if err != nil {
			return nil, 0, err
		}
		yPred[i] = y
		if yPred[i] != df.Y[i] {
			E += ada.rowWeights[i]
		}
		W += ada.rowWeights[i]
	}
	E /= W
	// Compute the predictor weight
	var predWeight float64
	if !ada.Regression {
		predWeight = math.Log((1-E)/E) + math.Log(float64(ada.nClasses)-1)
	} else {
		predWeight = math.Log((1 - E) / E)
	}
	// Update the row weights
	for i, w := range ada.rowWeights {
		if yPred[i] != df.Y[i] {
			ada.rowWeights[i] = w * math.Exp(predWeight)
		}
	}
	return predictor, predWeight, nil
}

// Fit AdaBoost
func (ada AdaBoost) Fit(learner Learner, df *dataframe.DataFrame, rounds int) error {
	// Determine number of classes if the task is classification
	if !ada.Regression {
		var n, err = df.NClasses()
		if err != nil {
			return err
		}
		ada.nClasses = n
	}
	// Initialize weights
	ada.rowWeights = make([]float64, df.NRows())
	ada.predictorWeights = make([]float64, df.NRows())
	for i := range ada.rowWeights {
		ada.rowWeights[i] = 1
	}
	// Initialize predictors
	ada.predictors = make([]Predictor, rounds)
	// Go through the rounds
	for i := range ada.predictors {
		var predictor, predWeight, err = ada.fit(learner, df)
		if err != nil {
			return err
		}
		ada.predictors[i] = predictor
		ada.predictorWeights[i] = predWeight
	}
	return nil
}

// PredictRow aggregates the votes of all the Predictors and outputs a final
// value or class depending on the task at hand.
func (ada AdaBoost) PredictRow(x []float64) (float64, error) {
	// Regression
	if ada.Regression {
		var (
			y float64
			w float64
		)
		for i, predictor := range ada.predictors {
			var p, err = predictor.PredictRow(x)
			if err != nil {
				return 0, err
			}
			y += p * ada.predictorWeights[i]
			w += ada.predictorWeights[i]
		}
		return y / w, nil
	}
	// Classification
	var votes = make([]float64, ada.nClasses)
	for i, predictor := range ada.predictors {
		var c, err = predictor.PredictRow(x)
		if err != nil {
			return 0, err
		}
		votes[int(c)] += ada.predictorWeights[i]
	}
	return float64(ArgMax(votes)), nil
}

// Predict runs PredictRow on each row in a dataframe.DataFrame.
func (ada AdaBoost) Predict(df *dataframe.DataFrame) ([]float64, error) {
	var Y = make([]float64, df.NRows())
	for i, x := range df.X {
		var y, err = ada.PredictRow(x)
		if err != nil {
			return nil, err
		}
		Y[i] = y
	}
	return Y, nil
}
