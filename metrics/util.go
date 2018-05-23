package metrics

import "strings"

// ParseMetric returns a Metric from it's String representation.
func ParseMetric(name string, class int) (Metric, error) {
	// If the name begins with "neg_" then a NegativeMetric is returns
	var negMetric = strings.HasPrefix(name, "neg_")
	if negMetric {
		name = strings.TrimLeft(name, "neg_")
	}
	var (
		klass      = float64(class)
		metric, ok = map[string]Metric{
			BinaryLogLoss{}.String():              BinaryLogLoss{},
			Accuracy{}.String():                   Accuracy{},
			BinaryPrecision{}.String():            BinaryPrecision{Class: klass},
			MacroPrecision{}.String():             MacroPrecision{},
			MicroPrecision{}.String():             MicroPrecision{},
			WeightedPrecision{}.String():          WeightedPrecision{},
			BinaryRecall{}.String():               BinaryRecall{Class: klass},
			MacroRecall{}.String():                MacroRecall{},
			MicroRecall{}.String():                MicroRecall{},
			WeightedRecall{}.String():             WeightedRecall{},
			BinaryF1{}.String():                   BinaryF1{Class: klass},
			MacroF1{}.String():                    MacroF1{},
			MicroF1{}.String():                    MicroF1{},
			WeightedF1{}.String():                 WeightedF1{},
			MeanAbsoluteError{}.String():          MeanAbsoluteError{},
			MeanSquaredError{}.String():           MeanSquaredError{},
			RootMeanSquaredError{}.String():       RootMeanSquaredError{},
			R2{}.String():                         R2{},
			ROCAUC{}.String():                     ROCAUC{},
			AbsolutePearsonCorrelation{}.String(): AbsolutePearsonCorrelation{},
		}[name]
	)
	if !ok {
		return nil, &errUnknownMetric{name}
	}
	if negMetric {
		metric = NegativeMetric{Metric: metric}
	}
	return metric, nil
}

func clip(f, min, max float64) float64 {
	if f < min {
		return min
	}
	if f > max {
		return max
	}
	return f
}
