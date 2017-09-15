package metrics

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRegression(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		metric  Metric
		score   float64
		err     error
	}{
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: nil,
			metric:  MeanAbsoluteError{},
			score:   0.5,
			err:     nil,
		},
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: []float64{2, 1, 1, 2},
			metric:  MeanAbsoluteError{},
			score:   (1 + 0.5 + 2) / 6, // 0.58333...
			err:     nil,
		},
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: nil,
			metric:  MeanSquaredError{},
			score:   0.375,
			err:     nil,
		},
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: []float64{2, 1, 1, 2},
			metric:  MeanSquaredError{},
			score:   (0.5 + 0.25 + 2) / 6, // 0.458333...
			err:     nil,
		},
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: nil,
			metric:  R2{},
			score:   1 - 1.5/29.1875,
			err:     nil,
		},
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: []float64{2, 1, 1, 2},
			metric:  R2{},
			score:   0.936354869816779,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		var score, err = tc.metric.Apply(tc.yTrue, tc.yPred, tc.weights)
		if fmt.Sprintf("%.15f", score) != fmt.Sprintf("%.15f", tc.score) || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Expected %.15f got %.15f in test case number %d", tc.score, score, i)
		}
	}
}
