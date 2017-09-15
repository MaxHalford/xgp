package metrics

import (
	"reflect"
	"testing"
)

func TestAccuracy(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		metric  Metric
		score   float64
		err     error
	}{
		{
			yTrue:   []float64{0, 2, 1, 3},
			yPred:   []float64{0, 1, 2, 3},
			weights: nil,
			metric:  Accuracy{},
			score:   2.0 / 4,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 2, 1, 3},
			yPred:   []float64{0, 1, 2, 3},
			weights: nil,
			metric:  NegativeAccuracy{},
			score:   -2.0 / 4,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 2, 1, 3},
			yPred:   []float64{0, 1, 2, 3},
			weights: []float64{2, 1, 1, 2},
			metric:  Accuracy{},
			score:   4.0 / 6, // 0.666...
			err:     nil,
		},
	}
	for i, tc := range testCases {
		var score, err = tc.metric.Apply(tc.yTrue, tc.yPred, tc.weights)
		if score != tc.score || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Expected %.3f got %.3f in test case number %d", tc.score, score, i)
		}
	}
}
