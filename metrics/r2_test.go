package metrics

import (
	"fmt"
	"testing"
)

func TestR2(t *testing.T) {
	var testCases = []metricTestCase{
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: nil,
			metric:  R2{},
			score:   1 - 1.5/29.1875,
			err:     nil,
		},
		{
			yTrue:   []float64{3, -0.5, 2, 7},
			yPred:   []float64{2.5, 0, 2, 8},
			weights: []float64{2, 1, 1, 2},
			metric:  R2{},
			score:   0.936354869816779,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
