package metrics

import (
	"fmt"
	"testing"
)

func TestMeanSquaredError(t *testing.T) {
	var testCases = []metricTestCase{
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: nil,
			metric:  MeanSquaredError{},
			score:   0.375,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
