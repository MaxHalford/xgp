package metrics

import (
	"fmt"
	"testing"
)

func TestF1(t *testing.T) {
	var testCases = []metricTestCase{
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  F1{Class: 0},
			score:   2 * (2.0 / 3 * 1) / (2.0/3 + 1), // 0.8
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  F1{Class: 1},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  F1{Class: 2},
			score:   0,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  MicroF1{},
			score:   2.0 * (1.0 / 3 * 1.0 / 3) / (1.0/3 + 1.0/3), // 0.333...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  MacroF1{},
			score:   2.0 * (2.0 / 9 * 1.0 / 3) / (2.0/9 + 1.0/3), // 0.266...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 0, 1, 2},
			yPred:   []float64{0, 2, 1, 0, 0, 1},
			weights: nil,
			metric:  WeightedF1{},
			score:   (2.0*2*(2.0/3*1)/(2.0/3+1) + 0 + 0) / 6, // 0.266...
			err:     nil,
		},
		{
			yTrue:   []float64{0, 0, 0, 1, 1, 1},
			yPred:   []float64{1, 0, 0, 1, 1, 0},
			weights: []float64{2, 1, 1, 1, 1, 2},
			metric:  F1{Class: 1},
			score:   4.0 / 8,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
