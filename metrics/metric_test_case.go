package metrics

import (
	"reflect"
	"testing"
)

type metricTestCase struct {
	yTrue   []float64
	yPred   []float64
	weights []float64
	metric  Metric
	score   float64
	err     error
}

func (tc metricTestCase) Run(t *testing.T) {
	var score, err = tc.metric.Apply(tc.yTrue, tc.yPred, tc.weights)
	if score != tc.score || !reflect.DeepEqual(err, tc.err) {
		t.Errorf("Expected %.3f, got %.3f", tc.score, score)
	}
}
