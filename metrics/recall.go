package metrics

// Recall measures the fraction of times a true class was predicted.
type Recall struct {
	Class float64
}

// Apply Recall.
func (recall Recall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP = cm.TruePositives(recall.Class)
		FN = cm.FalseNegatives(recall.Class)
	)
	// If the class has never been predicted return 0
	if TP+FN == 0 {
		return 0, nil
	}
	return TP / (TP + FN), nil
}

// Classification method of Recall.
func (recall Recall) Classification() bool {
	return true
}

// BiggerIsBetter method of Recall.
func (recall Recall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of Recall.
func (recall Recall) NeedsProbabilities() bool {
	return false
}

// String method of Recall.
func (recall Recall) String() string {
	return "recall"
}

// MicroRecall measures the global recall by using the total true positives and
// false negatives.
type MicroRecall struct{}

// Apply MicroRecall.
func (recall MicroRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
func (recall MicroRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of MicroRecall.
func (recall MicroRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MicroRecall.
func (recall MicroRecall) NeedsProbabilities() bool {
	return false
}

// String method of MicroRecall.
func (recall MicroRecall) String() string {
	return "micro_recall"
}

// MacroRecall measures the unweighted average recall across all classes.
// This does not take class imbalance into account.
type MacroRecall struct{}

// Apply MacroRecall.
func (recall MacroRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, class := range cm.Classes() {
		var recall, _ = Recall{Class: class}.Apply(yTrue, yPred, weights)
		sum += recall
	}
	return sum / float64(cm.NClasses()), nil
}

// Classification method of MacroRecall.
func (recall MacroRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of MacroRecall.
func (recall MacroRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MacroRecall.
func (recall MacroRecall) NeedsProbabilities() bool {
	return false
}

// String method of MacroRecall.
func (recall MacroRecall) String() string {
	return "macro_recall"
}

// WeightedRecall measures the weighted average recall across all classes.
// This does take class imbalance into account.
type WeightedRecall struct{}

// Apply WeightedRecall.
func (recall WeightedRecall) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
			recall, _ = Recall{Class: class}.Apply(yTrue, yPred, weights)
			TP        = cm.TruePositives(class)
			FN        = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * recall
		n += (TP + FN)
	}
	return sum / n, nil
}

// Classification method of WeightedRecall.
func (recall WeightedRecall) Classification() bool {
	return true
}

// BiggerIsBetter method of WeightedRecall.
func (recall WeightedRecall) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of WeightedRecall.
func (recall WeightedRecall) NeedsProbabilities() bool {
	return false
}

// String method of WeightedRecall.
func (recall WeightedRecall) String() string {
	return "weighted_recall"
}
