package metrics

import (
	"reflect"
	"testing"
)

func TestPrecision(t *testing.T) {
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
			metric:  BinaryPrecision{Class: 0},
			score:   2.0 / 3,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  BinaryPrecision{Class: 1},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  BinaryPrecision{Class: 2},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeBinaryPrecision{Class: 0},
			score:   -2.0 / 3,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  MicroPrecision{},
			score:   2.0 / 6,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeMicroPrecision{},
			score:   -2.0 / 6,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  MacroPrecision{},
			score:   (2.0/3 + 0 + 0) / 3, // 0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeMacroPrecision{},
			score:   -(2.0/3 + 0 + 0) / 3, // -0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  WeightedPrecision{},
			score:   (2*2.0/3 + 2*0 + 1*0) / 6, // 0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  NegativeWeightedPrecision{},
			score:   -(2*2.0/3 + 2*0 + 1*0) / 6, // -0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 0, 0, 1, 1, 1},
			yPred:   []float64{1, 0, 0, 1, 1, 0},
			weights: []float64{2, 1, 1, 1, 1, 2},
			metric:  BinaryPrecision{Class: 1},
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
