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
	TP, err := cm.TruePositives(metric.Class)
	// Check class exists
	if err != nil {
		return 0, err
	}
	FP, err := cm.FalsePositives(metric.Class)
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

// String method of BinaryPrecision.
func (metric BinaryPrecision) String() string {
	return "precision"
}

// NegativeBinaryPrecision measures the negative precision.
type NegativeBinaryPrecision struct {
	Class float64
}

// Apply NegativeBinaryPrecision.
func (metric NegativeBinaryPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var precision, err = BinaryPrecision{Class: metric.Class}.Apply(yTrue, yPred, weights)
	return -precision, err
}

// Classification method of NegativeBinaryPrecision.
func (metric NegativeBinaryPrecision) Classification() bool {
	return true
}

// String method of NegativeBinaryPrecision.
func (metric NegativeBinaryPrecision) String() string {
	return "neg_precision"
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
		tp, err := cm.TruePositives(class)
		if err != nil {
			return 0, err
		}
		fp, err := cm.FalsePositives(class)
		if err != nil {
			return 0, err
		}
		TP += tp
		FP += fp
	}
	return TP / (TP + FP), nil
}

// Classification method of MicroPrecision.
func (metric MicroPrecision) Classification() bool {
	return true
}

// String method of MicroPrecision.
func (metric MicroPrecision) String() string {
	return "micro_precision"
}

// NegativeMicroPrecision measures the negative micro precision.
type NegativeMicroPrecision struct{}

// Apply NegativeMicroPrecision.
func (metric NegativeMicroPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var precision, err = MicroPrecision{}.Apply(yTrue, yPred, weights)
	return -precision, err
}

// Classification method of NegativeMicroPrecision.
func (metric NegativeMicroPrecision) Classification() bool {
	return true
}

// String method of NegativeMicroPrecision.
func (metric NegativeMicroPrecision) String() string {
	return "neg_micro_precision"
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

// String method of MacroPrecision.
func (metric MacroPrecision) String() string {
	return "macro_precision"
}

// NegativeMacroPrecision measures the negative micro precision.
type NegativeMacroPrecision struct{}

// Apply NegativeMacroPrecision.
func (metric NegativeMacroPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var precision, err = MacroPrecision{}.Apply(yTrue, yPred, weights)
	return -precision, err
}

// Classification method of NegativeMacroPrecision.
func (metric NegativeMacroPrecision) Classification() bool {
	return true
}

// String method of NegativeMacroPrecision.
func (metric NegativeMacroPrecision) String() string {
	return "neg_macro_precision"
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
			TP, _        = cm.TruePositives(class)
			FN, _        = cm.FalseNegatives(class)
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

// String method of WeightedPrecision.
func (metric WeightedPrecision) String() string {
	return "weighted_precision"
}

// NegativeWeightedPrecision measures the negative micro precision.
type NegativeWeightedPrecision struct{}

// Apply NegativeWeightedPrecision.
func (metric NegativeWeightedPrecision) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var precision, err = WeightedPrecision{}.Apply(yTrue, yPred, weights)
	return -precision, err
}

// Classification method of NegativeWeightedPrecision.
func (metric NegativeWeightedPrecision) Classification() bool {
	return true
}

// String method of NegativeWeightedPrecision.
func (metric NegativeWeightedPrecision) String() string {
	return "neg_weighted_precision"
}
