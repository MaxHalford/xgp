package metrics

// BinaryF1 measures the F1-score.
type BinaryF1 struct {
	Class float64
}

// Apply BinaryF1.
func (metric BinaryF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP        = cm.TruePositives(metric.Class)
		FP        = cm.FalsePositives(metric.Class)
		FN        = cm.FalseNegatives(metric.Class)
		precision = TP / (TP + FP)
		recall    = TP / (TP + FN)
	)
	if precision == 0 || recall == 0 {
		return 0, nil
	}
	var f1 = 2 * (precision * recall) / (precision + recall)
	return f1, nil
}

// Classification method of BinaryF1.
func (metric BinaryF1) Classification() bool {
	return true
}

// BiggerIsBetter method of BinaryF1.
func (metric BinaryF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of BinaryF1.
func (metric BinaryF1) NeedsProbabilities() bool {
	return false
}

// String method of BinaryF1.
func (metric BinaryF1) String() string {
	return "f1"
}

// MicroF1 measures the global F1 score.
type MicroF1 struct{}

// Apply MicroF1.
func (metric MicroF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var microPrecision, err = MicroPrecision{}.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	microRecall, err := MicroRecall{}.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var microF1 = 2 * (microPrecision * microRecall) / (microPrecision + microRecall)
	return microF1, nil
}

// Classification method of MicroF1.
func (metric MicroF1) Classification() bool {
	return true
}

// BiggerIsBetter method of MicroF1.
func (metric MicroF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MicroF1.
func (metric MicroF1) NeedsProbabilities() bool {
	return false
}

// String method of MicroF1.
func (metric MicroF1) String() string {
	return "micro_f1"
}

// MacroF1 measures the global F1 score.
type MacroF1 struct{}

// Apply MacroF1.
func (metric MacroF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var macroPrecision, err = MacroPrecision{}.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	macroRecall, err := MacroRecall{}.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var macroF1 = 2 * (macroPrecision * macroRecall) / (macroPrecision + macroRecall)
	return macroF1, nil
}

// Classification method of MacroF1.
func (metric MacroF1) Classification() bool {
	return true
}

// BiggerIsBetter method of MacroF1.
func (metric MacroF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MacroF1.
func (metric MacroF1) NeedsProbabilities() bool {
	return false
}

// String method of MacroF1.
func (metric MacroF1) String() string {
	return "macro_f1"
}

// WeightedF1 measures the weighted average F1 score across all classes.
// This does take class imbalance into account.
type WeightedF1 struct{}

// Apply WeightedF1.
func (metric WeightedF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
			f1, _ = BinaryF1{Class: class}.Apply(yTrue, yPred, weights)
			TP    = cm.TruePositives(class)
			FN    = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * f1
		n += (TP + FN)
	}
	return sum / n, nil
}

// Classification method of WeightedF1.
func (metric WeightedF1) Classification() bool {
	return true
}

// BiggerIsBetter method of WeightedF1.
func (metric WeightedF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of WeightedF1.
func (metric WeightedF1) NeedsProbabilities() bool {
	return false
}

// String method of WeightedF1.
func (metric WeightedF1) String() string {
	return "weighted_f1"
}
