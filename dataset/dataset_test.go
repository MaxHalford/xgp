package dataset

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDatasetXT(t *testing.T) {
	var testCases = []struct {
		X  [][]float64
		XT [][]float64
	}{
		{
			X: [][]float64{
				[]float64{1, 2, 3},
				[]float64{4, 5, 6},
			},
			XT: [][]float64{
				[]float64{1, 4},
				[]float64{2, 5},
				[]float64{3, 6},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				ds = Dataset{X: tc.X}
				XT = ds.XT()
			)
			if !reflect.DeepEqual(XT, tc.XT) {
				t.Errorf("Got %v, expected %v", XT, tc.XT)
			}
		})
	}

}
