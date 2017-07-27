package metric

import (
	"reflect"
	"testing"
)

func TestAccuracy(t *testing.T) {
	var testCases = []struct {
		yTrue       []float64
		yPred       []float64
		metricuator Metric
		score       float64
		err         error
	}{
		{
			yTrue:       []float64{0, 2, 1, 3},
			yPred:       []float64{0, 1, 2, 3},
			metricuator: Accuracy{},
			score:       2.0 / 4,
			err:         nil,
		},
		{
			yTrue:       []float64{0, 2, 1, 3},
			yPred:       []float64{0, 1, 2, 3},
			metricuator: NegativeAccuracy{},
			score:       -2.0 / 4,
			err:         nil,
		},
	}
	for i, tc := range testCases {
		var score, err = tc.metricuator.Apply(tc.yTrue, tc.yPred)
		if score != tc.score || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Expected %.3f got %.3f in test case number %d", tc.score, score, i)
		}
	}
}
