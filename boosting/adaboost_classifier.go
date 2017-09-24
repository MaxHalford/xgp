package boosting

import (
	"encoding/json"
	"math"
)

// AdaBoostClassifier implements adaptive boosting for classification. The
// implementation is based on the SAMME algorithm proposed in Zhu et al (2009).
type AdaBoostClassifier struct {
	rowWeights       []float64   `json:"row_weights"`
	predictors       []Predictor `json:"predictors"`
	predictorWeights []float64   `json:"predictor_weights"`
	nClasses         int         `json:"n_classes"`
}

// fit represents one round of boosting.
func (ada AdaBoostClassifier) fit(learner Learner, X [][]float64, Y []float64) (Predictor, float64, error) {
	// Fit the weak learner
	var (
		predictor, err = learner.Learn(X, Y)
		n              = len(X)
		yPred          = make([]float64, n)
	)
	if err != nil {
		return nil, 0, err
	}
	// Compute weighted error
	var E, W float64
	for i, x := range X {
		y, err := predictor.PredictRow(x)
		if err != nil {
			return nil, 0, err
		}
		yPred[i] = y
		if yPred[i] != Y[i] {
			E += ada.rowWeights[i]
		}
		W += ada.rowWeights[i]
	}
	E /= W
	// Compute the predictor weight
	var predWeight float64
	predWeight = math.Log((1-E)/E) + math.Log(float64(ada.nClasses)-1)
	// Update the row weights
	for i, w := range ada.rowWeights {
		if yPred[i] != Y[i] {
			ada.rowWeights[i] = w * math.Exp(predWeight)
		}
	}
	return predictor, predWeight, nil
}

// Fit AdaBoost
func (ada AdaBoostClassifier) Fit(learner Learner, X [][]float64, Y []float64, rounds int) error {
	// Determine the number of classes
	ada.nClasses = 2
	// Initialize weights
	var n = len(X)
	ada.rowWeights = make([]float64, n)
	ada.predictorWeights = make([]float64, n)
	for i := range ada.rowWeights {
		ada.rowWeights[i] = 1
	}
	// Initialize predictors
	ada.predictors = make([]Predictor, rounds)
	// Go through the rounds
	for i := range ada.predictors {
		var predictor, predWeight, err = ada.fit(learner, X, Y)
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
func (ada AdaBoostClassifier) PredictRow(x []float64) (float64, error) {
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

// Predict runs PredictRow on each set of features of a slice of features.
func (ada AdaBoostClassifier) Predict(X [][]float64) ([]float64, error) {
	var Y = make([]float64, len(X))
	for i, x := range X {
		var y, err = ada.PredictRow(x)
		if err != nil {
			return nil, err
		}
		Y[i] = y
	}
	return Y, nil
}

// MarshalJSON serializes an AdaBoostClassifier into JSON.
func (ada *AdaBoostClassifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ada)
}

// UnmarshalJSON parses JSON into an AdaBoostClassifier.
func (ada *AdaBoostClassifier) UnmarshalJSON(bytes []byte) error {
	if err := json.Unmarshal(bytes, &ada); err != nil {
		return err
	}
	return nil
}
