package boosting

import "github.com/MaxHalford/xgp/dataframe"

// A Learner can be used by a Booster to fit a dataframe.DataFrame.
type Learner interface {
	Learn(df *dataframe.DataFrame) (Predictor, error)
}
