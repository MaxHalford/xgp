package ensemble

// A Predictor can be used to predict a single row of features. A set of
// Predictors can then be used by a Booster to vote for a final prediction.
type Predictor interface {
	PredictRow(x []float64, predictProba bool) (float64, error)
}
