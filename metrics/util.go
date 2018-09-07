package metrics

import "strings"

// ParseMetric returns a Metric from it's String representation.
func ParseMetric(name string, class int) (Metric, error) {
	var neg = strings.HasPrefix(name, "neg_")
	if neg {
		name = strings.TrimLeft(name, "neg_")
	}
	var (
		metrics = []Metric{
			LogLoss{},
			Accuracy{},
			Precision{Class: float64(class)},
			MacroPrecision{},
			MicroPrecision{},
			WeightedPrecision{},
			Recall{Class: float64(class)},
			MacroRecall{},
			MicroRecall{},
			WeightedRecall{},
			F1{Class: float64(class)},
			MacroF1{},
			MicroF1{},
			WeightedF1{},
			MAE{},
			MSE{},
			RMSE{},
			R2{},
			ROCAUC{},
			AbsolutePearson{},
		}
		metric Metric
	)
	for _, m := range metrics {
		if m.String() == name {
			metric = m
			break
		}
	}
	if metric == nil {
		return nil, errUnknownMetric{name}
	}
	if neg {
		metric = Negative{Metric: metric}
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
