package metrics

import (
	"fmt"
	"testing"
)

func TestAccuracy(t *testing.T) {
	var testCases = []metricTestCase{
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
			weights: []float64{2, 1, 1, 2},
			metric:  Accuracy{},
			score:   4.0 / 6,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
