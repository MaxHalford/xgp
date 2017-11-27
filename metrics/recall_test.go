package metrics

import (
	"fmt"
	"testing"
)

func TestRecall(t *testing.T) {
	var testCases = []metricTestCase{
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
			metric:  MicroRecall{},
			score:   2.0 / 6,
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
			metric:  WeightedRecall{},
			score:   (2*1.0 + 2*0 + 2*0) / 6, // 0.333...
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
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
