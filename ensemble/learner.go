package ensemble

// A Learner can be used by an Ensemble to train on a dataset.
type Learner interface {
	Fit(
		XTrain [][]float64,
		YTrain []float64,
		WTrain []float64,
		XVal [][]float64,
		YVal []float64,
		WVal []float64,
		verbose bool,
	) (Predictor, error)
}
