package metrics

import (
	"reflect"
	"testing"
)

func TestRecall(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		metric  Metric
		score   float64
		err     error
	}{
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  BinaryRecall{Class: 0},
			score:   1,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  BinaryRecall{Class: 1},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  BinaryRecall{Class: 2},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeBinaryRecall{Class: 0},
			score:   -1,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  MicroRecall{},
			score:   2.0 / 6,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeMicroRecall{},
			score:   -2.0 / 6,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  MacroRecall{},
			score:   (1.0 + 0 + 0) / 3,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeMacroRecall{},
			score:   -(1.0 + 0 + 0) / 3,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  WeightedRecall{},
			score:   (2*1.0 + 2*0 + 2*0) / 6, // 0.333...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeWeightedRecall{},
			score:   -(2*1.0 + 2*0 + 2*0) / 6, // -0.333...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 0, 0, 1, 1, 1},
			yPred:   []float64{1, 0, 0, 1, 1, 0},
			weights: []float64{2, 1, 1, 1, 1, 2},
			metric:  BinaryRecall{Class: 1},
			score:   4.0 / 8,
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
