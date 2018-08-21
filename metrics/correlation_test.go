package metrics

import (
	"fmt"
	"testing"
)

func TestAbsolutePearson(t *testing.T) {
	var testCases = []metricTestCase{
		{
			yTrue:   []float64{0, 1, 2, 3},
			yPred:   []float64{1, 2, 3, 4},
			weights: nil,
			metric:  AbsolutePearson{},
			score:   1,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 2, 3},
			yPred:   []float64{-1, -2, -3, -4},
			weights: nil,
			metric:  AbsolutePearson{},
			score:   1,
			err:     nil,
		},
		{
			yTrue:   []float64{0, 1, 0, 1},
			yPred:   []float64{1, 0, 1, 0},
			weights: nil,
			metric:  AbsolutePearson{},
			score:   1,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
