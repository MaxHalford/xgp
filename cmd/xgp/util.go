package main

import (
	"fmt"
	"os"

	"github.com/MaxHalford/xgp/metric"
	"github.com/urfave/cli"
)

func fileExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("No file named '%s'", file), 1)
	}
	return nil
}

func getMetric(metricName string, class float64) (metric.Metric, error) {
	var metric, ok = map[string]metric.Metric{
		"accuracy":               metric.Accuracy{},
		"neg_accuracy":           metric.NegativeAccuracy{},
		"binary_precision":       metric.BinaryPrecision{Class: class},
		"neg_binary_precision":   metric.NegativeBinaryPrecision{Class: class},
		"macro_precision":        metric.MacroPrecision{},
		"neg_macro_precision":    metric.NegativeMacroPrecision{},
		"micro_precision":        metric.MicroPrecision{},
		"neg_micro_precision":    metric.NegativeMicroPrecision{},
		"weighted_precision":     metric.WeightedPrecision{},
		"neg_weighted_precision": metric.NegativeWeightedPrecision{},
		"binary_recall":          metric.BinaryRecall{Class: class},
		"neg_binary_recall":      metric.NegativeBinaryRecall{Class: class},
		"macro_recall":           metric.MacroRecall{},
		"neg_macro_recall":       metric.NegativeMacroRecall{},
		"micro_recall":           metric.MicroRecall{},
		"neg_micro_recall":       metric.NegativeMicroRecall{},
		"weighted_recall":        metric.WeightedRecall{},
		"neg_weighted_recall":    metric.NegativeWeightedRecall{},
		"binary_f1_score":        metric.BinaryF1Score{Class: class},
		"neg_binary_f1_score":    metric.NegativeBinaryF1Score{Class: class},
		"macro_f1_score":         metric.MacroF1Score{},
		"neg_macro_f1_score":     metric.NegativeMacroF1Score{},
		"micro_f1_score":         metric.MicroF1Score{},
		"neg_micro_f1_score":     metric.NegativeMicroF1Score{},
		"weighted_f1_score":      metric.WeightedF1Score{},
		"neg_weighted_f1_score":  metric.NegativeWeightedF1Score{},
		"mean_absolute_error":    metric.MeanAbsoluteError{},
		"mean_squared_error":     metric.MeanSquaredError{},
		"r2":                     metric.R2{},
	}[metricName]
	if !ok {
		return metric, fmt.Errorf("Unknown metric name '%s'", metricName)
	}
	return metric, nil
}
