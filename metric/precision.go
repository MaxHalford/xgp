package metric

// BinaryPrecision measures the fraction of times a class was correctly predicted.
type BinaryPrecision struct {
	Class float64
}

// Apply BinaryPrecision.
func (metric BinaryPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
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
	// If the class has never been predicted return 0
	if TP+FP == 0 {
		return 0, nil
	}
	return TP / (TP + FP), nil
}

// NegativeBinaryPrecision measures the negative precision.
type NegativeBinaryPrecision struct {
	Class float64
}

// Apply NegativeBinaryPrecision.
func (metric NegativeBinaryPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var precision, err = BinaryPrecision{Class: metric.Class}.Apply(yTrue, yPred)
	return -precision, err
}

// MicroPrecision measures the global precision by using the total true
// positives and false positives.
type MicroPrecision struct{}

// Apply MicroPrecision.
func (metric MicroPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
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

// NegativeMicroPrecision measures the negative micro precision.
type NegativeMicroPrecision struct{}

// Apply NegativeMicroPrecision.
func (metric NegativeMicroPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var precision, err = MicroPrecision{}.Apply(yTrue, yPred)
	return -precision, err
}

// MacroPrecision measures the unweighted average precision across all classes.
// This does not take class imbalance into account.
type MacroPrecision struct{}

// Apply MacroPrecision.
func (metric MacroPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var cm, err = MakeConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, class := range cm.Classes() {
		var precision, _ = BinaryPrecision{Class: class}.Apply(yTrue, yPred)
		sum += precision
	}
	return sum / float64(cm.NClasses()), nil
}

// NegativeMacroPrecision measures the negative micro precision.
type NegativeMacroPrecision struct{}

// Apply NegativeMacroPrecision.
func (metric NegativeMacroPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var precision, err = MacroPrecision{}.Apply(yTrue, yPred)
	return -precision, err
}

// WeightedPrecision measures the weighted average precision across all classes.
// This does take class imbalance into account.
type WeightedPrecision struct{}

// Apply WeightedPrecision.
func (metric WeightedPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
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
			precision, _ = BinaryPrecision{Class: class}.Apply(yTrue, yPred)
			TP, _        = cm.TruePositives(class)
			FN, _        = cm.FalseNegatives(class)
		)
		sum += (TP + FN) * precision
		n += (TP + FN)
	}
	return sum / n, nil
}

// NegativeWeightedPrecision measures the negative micro precision.
type NegativeWeightedPrecision struct{}

// Apply NegativeWeightedPrecision.
func (metric NegativeWeightedPrecision) Apply(yTrue []float64, yPred []float64) (float64, error) {
	var precision, err = WeightedPrecision{}.Apply(yTrue, yPred)
	return -precision, err
}
