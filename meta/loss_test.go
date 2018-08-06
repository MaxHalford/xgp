package meta

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSquareLoss(t *testing.T) {
	var testCases = []struct {
		yTrue, yPred, loss, grad []float64
	}{
		{
			yTrue: []float64{1, 2, 3},
			yPred: []float64{1, 2, 3},
			loss:  []float64{0, 0, 0},
			grad:  []float64{0, 0, 0},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var loss = SquareLoss{}.Eval(tc.yTrue, tc.yPred)
			if !reflect.DeepEqual(loss, tc.loss) {
				t.Errorf("Expected %v, got %v", tc.loss, loss)
			}
		})
	}
}
