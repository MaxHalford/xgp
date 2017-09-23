package metrics

// GetMetric returns a Metric from it's String representation.
func GetMetric(metricName string, class int) (Metric, error) {
	var (
		klass      = float64(class)
		metric, ok = map[string]Metric{
			Accuracy{}.String():                  Accuracy{},
			NegativeAccuracy{}.String():          NegativeAccuracy{},
			BinaryPrecision{}.String():           BinaryPrecision{Class: klass},
			NegativeBinaryPrecision{}.String():   NegativeBinaryPrecision{Class: klass},
			MacroPrecision{}.String():            MacroPrecision{},
			NegativeMacroPrecision{}.String():    NegativeMacroPrecision{},
			MicroPrecision{}.String():            MicroPrecision{},
			NegativeMicroPrecision{}.String():    NegativeMicroPrecision{},
			WeightedPrecision{}.String():         WeightedPrecision{},
			NegativeWeightedPrecision{}.String(): NegativeWeightedPrecision{},
			BinaryRecall{}.String():              BinaryRecall{Class: klass},
			NegativeBinaryRecall{}.String():      NegativeBinaryRecall{Class: klass},
			MacroRecall{}.String():               MacroRecall{},
			NegativeMacroRecall{}.String():       NegativeMacroRecall{},
			MicroRecall{}.String():               MicroRecall{},
			NegativeMicroRecall{}.String():       NegativeMicroRecall{},
			WeightedRecall{}.String():            WeightedRecall{},
			NegativeWeightedRecall{}.String():    NegativeWeightedRecall{},
			BinaryF1Score{}.String():             BinaryF1Score{Class: klass},
			NegativeBinaryF1Score{}.String():     NegativeBinaryF1Score{Class: klass},
			MacroF1Score{}.String():              MacroF1Score{},
			NegativeMacroF1Score{}.String():      NegativeMacroF1Score{},
			MicroF1Score{}.String():              MicroF1Score{},
			NegativeMicroF1Score{}.String():      NegativeMicroF1Score{},
			WeightedF1Score{}.String():           WeightedF1Score{},
			NegativeWeightedF1Score{}.String():   NegativeWeightedF1Score{},
			MeanAbsoluteError{}.String():         MeanAbsoluteError{},
			MeanSquaredError{}.String():          MeanSquaredError{},
			R2{}.String():                        R2{},
			NegativeR2{}.String():                NegativeR2{},
		}[metricName]
	)
	if !ok {
		return nil, &errUnknownMetric{metricName}
	}
	return metric, nil
}
