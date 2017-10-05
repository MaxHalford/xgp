package metrics

import (
	"fmt"
	"strings"
)

// GetMetric returns a Metric from it's String representation.
func GetMetric(metricName string, class int) (Metric, error) {
	// If the metricName begins with "neg_" then a NegativeMetric is returns
	var negMetric = strings.HasPrefix(metricName, "neg_")
	if negMetric {
		metricName = strings.TrimLeft(metricName, "neg_")
	}
	var (
		klass      = float64(class)
		metric, ok = map[string]Metric{
			Accuracy{}.String():          Accuracy{},
			BinaryPrecision{}.String():   BinaryPrecision{Class: klass},
			MacroPrecision{}.String():    MacroPrecision{},
			MicroPrecision{}.String():    MicroPrecision{},
			WeightedPrecision{}.String(): WeightedPrecision{},
			BinaryRecall{}.String():      BinaryRecall{Class: klass},
			MacroRecall{}.String():       MacroRecall{},
			MicroRecall{}.String():       MicroRecall{},
			WeightedRecall{}.String():    WeightedRecall{},
			BinaryF1Score{}.String():     BinaryF1Score{Class: klass},
			MacroF1Score{}.String():      MacroF1Score{},
			MicroF1Score{}.String():      MicroF1Score{},
			WeightedF1Score{}.String():   WeightedF1Score{},
			MeanAbsoluteError{}.String(): MeanAbsoluteError{},
			R2{}.String():                R2{},
		}[metricName]
	)
	if !ok {
		return nil, &errUnknownMetric{metricName}
	}
	if negMetric {
		metric = NegativeMetric{Metric: metric}
	}
	return metric, nil
}

// A NegativeMetric returns the negative output of a given Metric.
type NegativeMetric struct {
	Metric Metric
}

// Apply NegativeMetric.
func (metric NegativeMetric) Apply(yTrue, yPred, weights []float64) (float64, error) {
	var output, err = metric.Metric.Apply(yTrue, yPred, weights)
	if err != nil {
		return 0, err
	}
	return -output, nil
}

// Classification method of NegativeMetric.
func (metric NegativeMetric) Classification() bool {
	return metric.Metric.Classification()
}

// BiggerIsBetter method of NegativeMetric.
func (metric NegativeMetric) BiggerIsBetter() bool {
	return !metric.Metric.BiggerIsBetter()
}

// String method of NegativeMetric.
func (metric NegativeMetric) String() string {
	return fmt.Sprintf("neg_%s", metric.Metric.String())
}
