package metric

import (
	"reflect"
	"testing"
)

func TestRecall(t *testing.T) {
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
			metric: BinaryRecall{Class: 0},
			score:  1,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: BinaryRecall{Class: 1},
			score:  0,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: BinaryRecall{Class: 2},
			score:  0,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeBinaryRecall{Class: 0},
			score:  -1,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: MicroRecall{},
			score:  2.0 / 6,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeMicroRecall{},
			score:  -2.0 / 6,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: MacroRecall{},
			score:  (1.0 + 0 + 0) / 3,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeMacroRecall{},
			score:  -(1.0 + 0 + 0) / 3,
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: WeightedRecall{},
			score:  (2*1.0 + 2*0 + 2*0) / 6, // 0.333...
			err:    nil,
		},
		{
			yTrue:  []float64{0, 1, 2, 0, 1, 2},
			yPred:  []float64{0, 2, 1, 0, 0, 1},
			metric: NegativeWeightedRecall{},
			score:  -(2*1.0 + 2*0 + 2*0) / 6, // -0.333...
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
