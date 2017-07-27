package metric

import (
	"reflect"
	"testing"
)

func TestF1Score(t *testing.T) {
	var testCases = []struct {
		yTrue  []float64
		yPred  []float64
		metric Metric
		score  float64
		err    error
	}{
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: BinaryF1Score{Class: 0},
			score:  2 * (2.0 / 3 * 1) / (2.0/3 + 1), // 0.8
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: BinaryF1Score{Class: 1},
			score:  0,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: BinaryF1Score{Class: 2},
			score:  0,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeBinaryF1Score{Class: 0},
			score:  -2 * (2.0 / 3 * 1) / (2.0/3 + 1), // -0.8
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: MicroF1Score{},
			score:  2.0 * (1.0 / 3 * 1.0 / 3) / (1.0/3 + 1.0/3), // 0.333...
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeMicroF1Score{},
			score:  -2.0 * (1.0 / 3 * 1.0 / 3) / (1.0/3 + 1.0/3), // -0.333...
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: MacroF1Score{},
			score:  2.0 * (2.0 / 9 * 1.0 / 3) / (2.0/9 + 1.0/3), // 0.266...
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeMacroF1Score{},
			score:  -2.0 * (2.0 / 9 * 1.0 / 3) / (2.0/9 + 1.0/3), // -0.266...
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: WeightedF1Score{},
			score:  (2.0*2*(2.0/3*1)/(2.0/3+1) + 0 + 0) / 6, // 0.266...
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeWeightedF1Score{},
			score:  -(2.0*2*(2.0/3*1)/(2.0/3+1) + 0 + 0) / 6, // -0.266...
			err:    nil,
		},
	}
	for i, tc := range testCases {
		var score, err = tc.metric.Apply(tc.yTrue, tc.yPred)
		if score != tc.score || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Expected %.3f got %.3f in test case number %d", tc.score, score, i)
		}
	}
}
