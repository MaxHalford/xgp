package metrics

import (
	"fmt"
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

func fmtScore(score float64) string {
	return fmt.Sprintf("%.5f", score)
}

func (tc metricTestCase) Run(t *testing.T) {
	var score, err = tc.metric.Apply(tc.yTrue, tc.yPred, tc.weights)
	if fmtScore(score) != fmtScore(tc.score) || !reflect.DeepEqual(err, tc.err) {
		t.Errorf("Expected %s, got %s", fmtScore(tc.score), fmtScore(score))
	}
}
