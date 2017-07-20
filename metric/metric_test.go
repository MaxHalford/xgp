package metric

import (
	"reflect"
	"testing"
)

func TestMetrics(t *testing.T) {
	var testCases = []struct {
		yTrue     []float64
		yPred     []float64
		evaluator Metric
		score     float64
		err       error
	}{
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: Accuracy{},
			score:     4.0 / 6,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: InverseAccuracy{},
			score:     -4.0 / 6,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: Precision{Class: 0},
			score:     1.0 / 3,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: InversePrecision{Class: 0},
			score:     -1.0 / 3,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: Precision{Class: 1},
			score:     1,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: InversePrecision{Class: 1},
			score:     -1,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: Recall{Class: 0},
			score:     1,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: InverseRecall{Class: 0},
			score:     -1,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: Recall{Class: 2},
			score:     2.0 / 3,
			err:       nil,
		},
		{
			yTrue:     []float64{0, 1, 1, 2, 2, 2},
			yPred:     []float64{0, 1, 0, 2, 2, 0},
			evaluator: InverseRecall{Class: 2},
			score:     -2.0 / 3,
			err:       nil,
		},
	}
	for i, tc := range testCases {
		var score, err = tc.evaluator.Apply(tc.yTrue, tc.yPred)
		if score != tc.score || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Error in test case number %d", i)
		}
	}
}
