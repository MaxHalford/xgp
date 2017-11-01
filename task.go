package xgp

import "github.com/MaxHalford/xgp/metrics"

// A Task contains information a Program needs to know in order to at least
// make predictions.
type Task struct {
	Metric   metrics.Metric
	NClasses int // Should be equal to 0 if Classification is false
}

func (t Task) binaryClassification() bool {
	return t.Metric.Classification() && t.NClasses == 2
}

func (t Task) multiClassClassification() bool {
	return t.Metric.Classification() && t.NClasses > 2
}
