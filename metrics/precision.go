package metrics

// BinaryPrecision measures the fraction of times a class was correctly predicted.
type BinaryPrecision struct {
	Class float64
}

// Apply BinaryPrecision.
func (metric BinaryPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var (
		TP = cm.TruePositives(metric.Class)
		FP = cm.FalsePositives(metric.Class)
	)
	// If the class has never been predicted return 0
	if TP+FP == 0 {
		return 0, nil
	}
	return TP / (TP + FP), nil
}

// Classification method of BinaryPrecision.
func (metric BinaryPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of BinaryPrecision.
func (metric BinaryPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of BinaryPrecision.
func (metric BinaryPrecision) NeedsProbabilities() bool {
	return false
}

// String method of BinaryPrecision.
func (metric BinaryPrecision) String() string {
	return "precision"
}

// MicroPrecision measures the global precision by using the total true
// positives and false positives.
type MicroPrecision struct{}

// Apply MicroPrecision.
func (metric MicroPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
func (metric MicroPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of MicroPrecision.
func (metric MicroPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MicroPrecision.
func (metric MicroPrecision) NeedsProbabilities() bool {
	return false
}

// String method of MicroPrecision.
func (metric MicroPrecision) String() string {
	return "micro_precision"
}

// MacroPrecision measures the unweighted average precision across all classes.
// This does not take class imbalance into account.
type MacroPrecision struct{}

// Apply MacroPrecision.
func (metric MacroPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, class := range cm.Classes() {
		var precision, _ = BinaryPrecision{Class: class}.Apply(yTrue, yPred, weights)
		sum += precision
	}
	return sum / float64(cm.NClasses()), nil
}

// Classification method of MacroPrecision.
func (metric MacroPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of MacroPrecision.
func (metric MacroPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of MacroPrecision.
func (metric MacroPrecision) NeedsProbabilities() bool {
	return false
}

// String method of MacroPrecision.
func (metric MacroPrecision) String() string {
	return "macro_precision"
}

// WeightedPrecision measures the weighted average precision across all classes.
// This does take class imbalance into account.
type WeightedPrecision struct{}

// Apply WeightedPrecision.
func (metric WeightedPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
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
			precision, _ = BinaryPrecision{Class: class}.Apply(yTrue, yPred, weights)
			TP           = cm.TruePositives(class)
			FN           = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * precision
		n += (TP + FN)
	}
	return sum / n, nil
}

// Classification method of WeightedPrecision.
func (metric WeightedPrecision) Classification() bool {
	return true
}

// BiggerIsBetter method of WeightedPrecision.
func (metric WeightedPrecision) BiggerIsBetter() bool {
	return true
}

// NeedsProbabilities method of WeightedPrecision.
func (metric WeightedPrecision) NeedsProbabilities() bool {
	return false
}

// String method of WeightedPrecision.
func (metric WeightedPrecision) String() string {
	return "weighted_precision"
}
