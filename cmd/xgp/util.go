package xgp

import (
	"fmt"
	"os"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/urfave/cli"
)

func fileExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("No file named '%s'", file), 1)
	}
	return nil
}

func getMetric(metricName string, class int) (metrics.Metric, error) {
	var (
		klass      = float64(class)
		metric, ok = map[string]metrics.Metric{
			"accuracy":               metrics.Accuracy{},
			"neg_accuracy":           metrics.NegativeAccuracy{},
			"binary_precision":       metrics.BinaryPrecision{Class: klass},
			"neg_binary_precision":   metrics.NegativeBinaryPrecision{Class: klass},
			"macro_precision":        metrics.MacroPrecision{},
			"neg_macro_precision":    metrics.NegativeMacroPrecision{},
			"micro_precision":        metrics.MicroPrecision{},
			"neg_micro_precision":    metrics.NegativeMicroPrecision{},
			"weighted_precision":     metrics.WeightedPrecision{},
			"neg_weighted_precision": metrics.NegativeWeightedPrecision{},
			"binary_recall":          metrics.BinaryRecall{Class: klass},
			"neg_binary_recall":      metrics.NegativeBinaryRecall{Class: klass},
			"macro_recall":           metrics.MacroRecall{},
			"neg_macro_recall":       metrics.NegativeMacroRecall{},
			"micro_recall":           metrics.MicroRecall{},
			"neg_micro_recall":       metrics.NegativeMicroRecall{},
			"weighted_recall":        metrics.WeightedRecall{},
			"neg_weighted_recall":    metrics.NegativeWeightedRecall{},
			"binary_f1_score":        metrics.BinaryF1Score{Class: klass},
			"neg_binary_f1_score":    metrics.NegativeBinaryF1Score{Class: klass},
			"macro_f1_score":         metrics.MacroF1Score{},
			"neg_macro_f1_score":     metrics.NegativeMacroF1Score{},
			"micro_f1_score":         metrics.MicroF1Score{},
			"neg_micro_f1_score":     metrics.NegativeMicroF1Score{},
			"weighted_f1_score":      metrics.WeightedF1Score{},
			"neg_weighted_f1_score":  metrics.NegativeWeightedF1Score{},
			"mae": metrics.MeanAbsoluteError{},
			"mse": metrics.MeanSquaredError{},
			"r2":  metrics.R2{},
		}[metricName]
	)
	if !ok {
		return nil, &errUnknownMetric{metricName}
	}
	return metric, nil
}

func isClassificationMetric(metricName string) (bool, error) {
	var is, ok = map[string]bool{
		"accuracy":               true,
		"neg_accuracy":           true,
		"binary_precision":       true,
		"neg_binary_precision":   true,
		"macro_precision":        true,
		"neg_macro_precision":    true,
		"micro_precision":        true,
		"neg_micro_precision":    true,
		"weighted_precision":     true,
		"neg_weighted_precision": true,
		"binary_recall":          true,
		"neg_binary_recall":      true,
		"macro_recall":           true,
		"neg_macro_recall":       true,
		"micro_recall":           true,
		"neg_micro_recall":       true,
		"weighted_recall":        true,
		"neg_weighted_recall":    true,
		"binary_f1_score":        true,
		"neg_binary_f1_score":    true,
		"macro_f1_score":         true,
		"neg_macro_f1_score":     true,
		"micro_f1_score":         true,
		"neg_micro_f1_score":     true,
		"weighted_f1_score":      true,
		"neg_weighted_f1_score":  true,
		"mae": false,
		"mse": false,
		"r2":  false,
	}[metricName]
	if !ok {
		return is, &errUnknownMetric{metricName}
	}
	return is, nil
}
