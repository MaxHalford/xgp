package metric

import (
	"reflect"
	"testing"
)

func TestPrecision(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		metrics Metric
		score   float64
		err     error
	}{
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: BinaryPrecision{Class: 0},
			score:   2.0 / 3,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: BinaryPrecision{Class: 1},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: BinaryPrecision{Class: 2},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: NegativeBinaryPrecision{Class: 0},
			score:   -2.0 / 3,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: MicroPrecision{},
			score:   2.0 / 6,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: NegativeMicroPrecision{},
			score:   -2.0 / 6,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: MacroPrecision{},
			score:   (2.0/3 + 0 + 0) / 3, // 0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: NegativeMacroPrecision{},
			score:   -(2.0/3 + 0 + 0) / 3, // -0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: WeightedPrecision{},
			score:   (2*2.0/3 + 2*0 + 1*0) / 6, // 0.222...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			metrics: NegativeWeightedPrecision{},
			score:   -(2*2.0/3 + 2*0 + 1*0) / 6, // -0.222...
			err:     nil,
		},
	}
	for i, tc := range testCases {
		var score, err = tc.metrics.Apply(tc.yTrue, tc.yPred)
		if score != tc.score || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Expected %.3f got %.3f in test case number %d", tc.score, score, i)
		}
	}
}
