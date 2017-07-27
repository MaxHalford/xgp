package metric

// BinaryF1Score measures the F1-score.
type BinaryF1Score struct {
	Class float64
}

// Apply BinaryF1Score.
func (metric BinaryF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	TP, err := cm.TruePositives(metric.Class)
	// Check class exists
	if err != nil {
		return 0, err
	}
	FP, err := cm.FalsePositives(metric.Class)
	FN, err := cm.FalseNegatives(metric.Class)
	// If the class has never been predicted return 0
	var (
		precision = TP / (TP + FP)
		recall    = TP / (TP + FN)
	)
	if precision == 0 || recall == 0 {
		return 0, nil
	}
	var f1Score = 2 * (precision * recall) / (precision + recall)
	return f1Score, nil
}

// NegativeBinaryF1Score measures the negative F1-score.
type NegativeBinaryF1Score struct {
	Class float64
}

// Apply NegativeBinaryF1Score.
func (metric NegativeBinaryF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var f1Score, err = BinaryF1Score{Class: metric.Class}.Apply(yTrue, yPred)
	return -f1Score, err
}

// MicroF1Score measures the global F1 score.
type MicroF1Score struct{}

// Apply MicroF1Score.
func (metric MicroF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var microPrecision, err = MicroPrecision{}.Apply(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	microRecall, err := MicroRecall{}.Apply(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	var microF1Score = 2 * (microPrecision * microRecall) / (microPrecision + microRecall)
	return microF1Score, nil
}

// NegativeMicroF1Score measures the negative micro F1 score.
type NegativeMicroF1Score struct{}

// Apply NegativeMicroF1Score.
func (metric NegativeMicroF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var f1Score, err = MicroF1Score{}.Apply(yTrue, yPred)
	return -f1Score, err
}

// MacroF1Score measures the global F1 score.
type MacroF1Score struct{}

// Apply MacroF1Score.
func (metric MacroF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var macroPrecision, err = MacroPrecision{}.Apply(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	macroRecall, err := MacroRecall{}.Apply(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	var macroF1Score = 2 * (macroPrecision * macroRecall) / (macroPrecision + macroRecall)
	return macroF1Score, nil
}

// NegativeMacroF1Score measures the negative micro precision.
type NegativeMacroF1Score struct{}

// Apply NegativeMacroF1Score.
func (metric NegativeMacroF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var f1Score, err = MacroF1Score{}.Apply(yTrue, yPred)
	return -f1Score, err
}

// WeightedF1Score measures the weighted average F1 score across all classes.
// This does take class imbalance into account.
type WeightedF1Score struct{}

// Apply WeightedF1Score.
func (metric WeightedF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	var (
		sum float64
		n   float64
	)
	for _, class := range cm.Classes() {
		var (
			f1Score, _ = BinaryF1Score{Class: class}.Apply(yTrue, yPred)
			TP, _      = cm.TruePositives(class)
			FN, _      = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * f1Score
		n += (TP + FN)
	}
	return sum / n, nil
}

// NegativeWeightedF1Score measures the negative micro precision.
type NegativeWeightedF1Score struct{}

// Apply NegativeWeightedF1Score.
func (metric NegativeWeightedF1Score) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var f1Score, err = WeightedF1Score{}.Apply(yTrue, yPred)
	return -f1Score, err
}
