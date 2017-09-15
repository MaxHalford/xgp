package metrics

// BinaryRecall measures the fraction of times a true class was predicted.
type BinaryRecall struct {
	Class float64
}

// Apply BinaryRecall.
func (metric BinaryRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	TP, err := cm.TruePositives(metric.Class)
	// Check class exists
	if err != nil {
		return 0, err
	}
	FN, err := cm.FalseNegatives(metric.Class)
	// If the class has never been predicted return 0
	if TP+FN == 0 {
		return 0, nil
	}
	return TP / (TP + FN), nil
}

// NegativeBinaryRecall measures the negative recall.
type NegativeBinaryRecall struct {
	Class float64
}

// Apply NegativeBinaryRecall.
func (metric NegativeBinaryRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var recall, err = BinaryRecall{Class: metric.Class}.Apply(yTrue, yPred, weights)
	return -recall, err
}

// MicroRecall measures the global recall by using the total true positives and
// false negatives.
type MicroRecall struct{}

// Apply MicroRecall.
func (metric MicroRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP float64
		FN float64
	)
	for _, class := range cm.Classes() {
		tp, err := cm.TruePositives(class)
		if err != nil {
			return 0, err
		}
		fn, err := cm.FalsePositives(class)
		if err != nil {
			return 0, err
		}
		TP += tp
		FN += fn
	}
	return TP / (TP + FN), nil
}

// NegativeMicroRecall measures the negative micro recall.
type NegativeMicroRecall struct{}

// Apply NegativeMicroRecall.
func (metric NegativeMicroRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var recall, err = MicroRecall{}.Apply(yTrue, yPred, weights)
	return -recall, err
}

// MacroRecall measures the unweighted average recall across all classes.
// This does not take class imbalance into account.
type MacroRecall struct{}

// Apply MacroRecall.
func (metric MacroRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, class := range cm.Classes() {
		var recall, _ = BinaryRecall{Class: class}.Apply(yTrue, yPred, weights)
		sum += recall
	}
	return sum / float64(cm.NClasses()), nil
}

// NegativeMacroRecall measures the negative micro recall.
type NegativeMacroRecall struct{}

// Apply NegativeMacroRecall.
func (metric NegativeMacroRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var recall, err = MacroRecall{}.Apply(yTrue, yPred, weights)
	return -recall, err
}

// WeightedRecall measures the weighted average recall across all classes.
// This does take class imbalance into account.
type WeightedRecall struct{}

// Apply WeightedRecall.
func (metric WeightedRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		sum float64
		n   float64
	)
	for _, class := range cm.Classes() {
		var (
			recall, _ = BinaryRecall{Class: class}.Apply(yTrue, yPred, weights)
			TP, _     = cm.TruePositives(class)
			FN, _     = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * recall
		n += (TP + FN)
	}
	return sum / n, nil
}

// NegativeWeightedRecall measures the negative micro recall.
type NegativeWeightedRecall struct{}

// Apply NegativeWeightedRecall.
func (metric NegativeWeightedRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var recall, err = WeightedRecall{}.Apply(yTrue, yPred, weights)
	return -recall, err
}
