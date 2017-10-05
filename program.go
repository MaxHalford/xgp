package xgp

// A Program holds a tree composed of Nodes and also holds a reference to an
// Estimator. A Program is simply an abstraction of top of a Node that allows
// not having to store the Estimator reference in each Node.
type Program struct {
	Root      *Node                  `json:"root"`
	Estimator *Estimator             `json:"-"`
	DRS       *DynamicRangeSelection `json:"drs"`
}

// String representation of a Program.
func (prog Program) String() string {
	return prog.Root.String()
}

// Clone a Program.
func (prog Program) clone() Program {
	var clone = Program{
		Root:      prog.Root.clone(),
		Estimator: prog.Estimator,
	}
	if prog.DRS != nil {
		clone.DRS = prog.DRS.clone()
	}
	return clone
}

// PredictRow predicts the output of some features.
func (prog Program) PredictRow(x []float64) (float64, error) {
	var y = prog.Root.evaluateRow(x)
	if prog.DRS != nil {
		return prog.DRS.PredictRow(y), nil
	}
	return y, nil
}

// Predict predicts the output of a slice of features.
func (prog Program) Predict(XT [][]float64) (yPred []float64, err error) {
	if prog.Estimator != nil {
		yPred, err = prog.Root.evaluateXT(XT, prog.Estimator.nodeCache)
	} else {
		yPred, err = prog.Root.evaluateXT(XT, nil)
	}
	if err != nil {
		return nil, err
	}
	if prog.DRS != nil {
		return prog.DRS.Predict(yPred), nil
	}
	return yPred, nil
}
