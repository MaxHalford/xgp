package cmd

type model interface {
	Predict(X [][]float64, proba bool) ([]float64, error)
}

type serialModel struct {
	Task   string `json:"task"`
	Flavor string `json:"flavor"`
	Model  model  `json:"model"`
}
