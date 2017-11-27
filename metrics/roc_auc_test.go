package metrics

import (
	"fmt"
	"testing"
)

func TestROCAUC(t *testing.T) {
	var testCases = []metricTestCase{
		{
			yTrue:   []float64{0, 0, 1, 1},
			yPred:   []float64{0.1, 0.4, 0.35, 0.8},
			weights: nil,
			metric:  ROCAUC{},
			score:   0.75,
			err:     nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.Run(t)
		})
	}
}
