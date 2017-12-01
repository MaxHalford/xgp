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
	var (
		TP = cm.TruePositives(metric.Class)
		FN = cm.FalseNegatives(metric.Class)
	)
	// If the class has never been predicted return 0
	if TP+FN == 0 {
		return 0, nil
	}
	return TP / (TP + FN), nil
}

// Classification method of BinaryRecall.
func (metric BinaryRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of BinaryRecall.
func (metric BinaryRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of BinaryRecall.
func (metric BinaryRecall) NeedsProbabilities() bool {
	return false
}

// String method of BinaryRecall.
func (metric BinaryRecall) String() string {
	return "recall"
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
		TP += cm.TruePositives(class)
		FN += cm.FalsePositives(class)
	}
	return TP / (TP + FN), nil
}

// Classification method of MicroRecall.
func (metric MicroRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of MicroRecall.
func (metric MicroRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MicroRecall.
func (metric MicroRecall) NeedsProbabilities() bool {
	return false
}

// String method of MicroRecall.
func (metric MicroRecall) String() string {
	return "micro_recall"
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

// Classification method of MacroRecall.
func (metric MacroRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of MacroRecall.
func (metric MacroRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MacroRecall.
func (metric MacroRecall) NeedsProbabilities() bool {
	return false
}

// String method of MacroRecall.
func (metric MacroRecall) String() string {
	return "macro_recall"
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
			TP        = cm.TruePositives(class)
			FN        = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * recall
		n += (TP + FN)
	}
	return sum / n, nil
}

// Classification method of WeightedRecall.
func (metric WeightedRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of WeightedRecall.
func (metric WeightedRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of WeightedRecall.
func (metric WeightedRecall) NeedsProbabilities() bool {
	return false
}

// String method of WeightedRecall.
func (metric WeightedRecall) String() string {
	return "weighted_recall"
}
