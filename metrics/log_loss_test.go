package metrics

import (
	"fmt"
	"testing"
)

func TestLogLoss(t *testing.T) {
	var testCases = []metricTestCase{
		{
			yTrue:   []float64{1},
			yPred:   []float64{0.5},
			weights: nil,
			metric:  LogLoss{},
			score:   0.69315,
			err:     nil,
		},
		{
			yTrue:   []float64{1},
			yPred:   []float64{0.9},
			weights: nil,
			metric:  LogLoss{},
			score:   0.10536,
			err:     nil,
		},
		{
			yTrue:   []float64{1, 1},
			yPred:   []float64{0.5, 0.9},
			weights: nil,
			metric:  LogLoss{},
			score:   0.39925, // (0.69315 + 0.10536) / 2
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
