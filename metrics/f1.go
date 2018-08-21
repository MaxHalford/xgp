package metrics

// F1 measures the F1-score.
type F1 struct {
	Class float64
}

// Apply F1.
func (f1 F1) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP        = cm.TruePositives(f1.Class)
		FP        = cm.FalsePositives(f1.Class)
		FN        = cm.FalseNegatives(f1.Class)
		precision = TP / (TP + FP)
		recall    = TP / (TP + FN)
	)
	if precision == 0 || recall == 0 {
		return 0, nil
	}
	return 2 * (precision * recall) / (precision + recall), nil
}

// Classification method of F1.
func (f1 F1) Classification() bool {
	return true
}

// BiggerIsBetter method of F1.
func (f1 F1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of F1.
func (f1 F1) NeedsProbabilities() bool {
	return false
}

// String method of F1.
func (f1 F1) String() string {
	return "f1"
}

// MicroF1 measures the global F1 score.
type MicroF1 struct{}

// Apply MicroF1.
func (f1 MicroF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
func (f1 MicroF1) Classification() bool {
	return true
}

// BiggerIsBetter method of MicroF1.
func (f1 MicroF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MicroF1.
func (f1 MicroF1) NeedsProbabilities() bool {
	return false
}

// String method of MicroF1.
func (f1 MicroF1) String() string {
	return "micro_f1"
}

// MacroF1 measures the global F1 score.
type MacroF1 struct{}

// Apply MacroF1.
func (f1 MacroF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
func (f1 MacroF1) Classification() bool {
	return true
}

// BiggerIsBetter method of MacroF1.
func (f1 MacroF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MacroF1.
func (f1 MacroF1) NeedsProbabilities() bool {
	return false
}

// String method of MacroF1.
func (f1 MacroF1) String() string {
	return "macro_f1"
}

// WeightedF1 measures the weighted average F1 score across all classes.
// This does take class imbalance into account.
type WeightedF1 struct{}

// Apply WeightedF1.
func (f1 WeightedF1) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
			f1, _ = F1{Class: class}.Apply(yTrue, yPred, weights)
			TP    = cm.TruePositives(class)
			FN    = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * f1
		n += (TP + FN)
	}
	return sum / n, nil
}

// Classification method of WeightedF1.
func (f1 WeightedF1) Classification() bool {
	return true
}

// BiggerIsBetter method of WeightedF1.
func (f1 WeightedF1) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of WeightedF1.
func (f1 WeightedF1) NeedsProbabilities() bool {
	return false
}

// String method of WeightedF1.
func (f1 WeightedF1) String() string {
	return "weighted_f1"
}
