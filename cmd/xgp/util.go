package main

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

func getMetric(metricName string, class float64) (metrics.Metric, error) {
	var metric, ok = map[string]metrics.Metric{
		"accuracy":               metrics.Accuracy{},
		"neg_accuracy":           metrics.NegativeAccuracy{},
		"binary_precision":       metrics.BinaryPrecision{Class: class},
		"neg_binary_precision":   metrics.NegativeBinaryPrecision{Class: class},
		"macro_precision":        metrics.MacroPrecision{},
		"neg_macro_precision":    metrics.NegativeMacroPrecision{},
		"micro_precision":        metrics.MicroPrecision{},
		"neg_micro_precision":    metrics.NegativeMicroPrecision{},
		"weighted_precision":     metrics.WeightedPrecision{},
		"neg_weighted_precision": metrics.NegativeWeightedPrecision{},
		"binary_recall":          metrics.BinaryRecall{Class: class},
		"neg_binary_recall":      metrics.NegativeBinaryRecall{Class: class},
		"macro_recall":           metrics.MacroRecall{},
		"neg_macro_recall":       metrics.NegativeMacroRecall{},
		"micro_recall":           metrics.MicroRecall{},
		"neg_micro_recall":       metrics.NegativeMicroRecall{},
		"weighted_recall":        metrics.WeightedRecall{},
		"neg_weighted_recall":    metrics.NegativeWeightedRecall{},
		"binary_f1_score":        metrics.BinaryF1Score{Class: class},
		"neg_binary_f1_score":    metrics.NegativeBinaryF1Score{Class: class},
		"macro_f1_score":         metrics.MacroF1Score{},
		"neg_macro_f1_score":     metrics.NegativeMacroF1Score{},
		"micro_f1_score":         metrics.MicroF1Score{},
		"neg_micro_f1_score":     metrics.NegativeMicroF1Score{},
		"weighted_f1_score":      metrics.WeightedF1Score{},
		"neg_weighted_f1_score":  metrics.NegativeWeightedF1Score{},
		"mean_absolute_error":    metrics.MeanAbsoluteError{},
		"mean_squared_error":     metrics.MeanSquaredError{},
		"r2":                     metrics.R2{},
	}[metricName]
	if !ok {
		return metric, fmt.Errorf("Unknown metric name '%s'", metricName)
	}
	return metric, nil
}
