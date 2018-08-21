package metrics

// Precision measures the fraction of times a class was correctly predicted.
type Precision struct {
	Class float64
}

// Apply Precision.
func (precision Precision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP = cm.TruePositives(precision.Class)
		FP = cm.FalsePositives(precision.Class)
	)
	// If the class has never been predicted return 0
	if TP+FP == 0 {
		return 0, nil
	}
	return TP / (TP + FP), nil
}

// Classification method of Precision.
func (precision Precision) Classification() bool {
	return true
}

// BiggerIsBetter method of Precision.
func (precision Precision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of Precision.
func (precision Precision) NeedsProbabilities() bool {
	return false
}

// String method of Precision.
func (precision Precision) String() string {
	return "precision"
}

// MicroPrecision measures the global precision by using the total true
// positives and false positives.
type MicroPrecision struct{}

// Apply MicroPrecision.
func (precision MicroPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP float64
		FP float64
	)
	for _, class := range cm.Classes() {
		TP += cm.TruePositives(class)
		FP += cm.FalsePositives(class)
	}
	return TP / (TP + FP), nil
}

// Classification method of MicroPrecision.
func (precision MicroPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of MicroPrecision.
func (precision MicroPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MicroPrecision.
func (precision MicroPrecision) NeedsProbabilities() bool {
	return false
}

// String method of MicroPrecision.
func (precision MicroPrecision) String() string {
	return "micro_precision"
}

// MacroPrecision measures the unweighted average precision across all classes.
// This does not take class imbalance into account.
type MacroPrecision struct{}

// Apply MacroPrecision.
func (precision MacroPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, class := range cm.Classes() {
		var precision, _ = Precision{Class: class}.Apply(yTrue, yPred, weights)
		sum += precision
	}
	return sum / float64(cm.NClasses()), nil
}

// Classification method of MacroPrecision.
func (precision MacroPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of MacroPrecision.
func (precision MacroPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MacroPrecision.
func (precision MacroPrecision) NeedsProbabilities() bool {
	return false
}

// String method of MacroPrecision.
func (precision MacroPrecision) String() string {
	return "macro_precision"
}

// WeightedPrecision measures the weighted average precision across all classes.
// This does take class imbalance into account.
type WeightedPrecision struct{}

// Apply WeightedPrecision.
func (precision WeightedPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
			precision, _ = Precision{Class: class}.Apply(yTrue, yPred, weights)
			TP           = cm.TruePositives(class)
			FN           = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * precision
		n += (TP + FN)
	}
	return sum / n, nil
}

// Classification method of WeightedPrecision.
func (precision WeightedPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of WeightedPrecision.
func (precision WeightedPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of WeightedPrecision.
func (precision WeightedPrecision) NeedsProbabilities() bool {
	return false
}

// String method of WeightedPrecision.
func (precision WeightedPrecision) String() string {
	return "weighted_precision"
}
