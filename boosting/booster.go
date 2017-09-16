package boosting

import "github.com/MaxHalford/xgp/dataframe"

// A Booster can fit a dataframe.DataFrame given a "weak" learner.
type Booster interface {
	Fit(learner Learner, df *dataframe.DataFrame, rounds int) error
	Predict(df *dataframe.DataFrame) ([]float64, error)
}
