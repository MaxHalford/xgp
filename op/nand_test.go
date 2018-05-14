package op

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNANDEval(t *testing.T) {
	var testCases = []struct {
		X [][]float64
		Y []float64
	}{
		{
			X: [][]float64{
				[]float64{1, 0, 1, 0, 1, 2},
				[]float64{1, 1, 0, 0, 2, 1},
			},
			Y: []float64{0, 1, 1, 1, 1, 1},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var Y = NAND{}.Eval(tc.X)
			if !reflect.DeepEqual(Y, tc.Y) {
				t.Errorf("Expected %v, got %v", tc.Y, Y)
			}
		})
	}
}
