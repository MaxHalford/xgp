package dataset

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestReadCSV(t *testing.T) {
	var testCases = []struct {
		str            string
		target         string
		classification bool
		dataset        *Dataset
	}{
		{
			str:            "x0,x1,y\n1,2,3\n-1,-2,-3",
			target:         "y",
			classification: false,
			dataset: &Dataset{
				X:      [][]float64{[]float64{1, -1}, []float64{2, -2}},
				XNames: []string{"x0", "x1"},
				Y:      []float64{3, -3},
				YName:  "y",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				path = fmt.Sprintf("test_read_csv_%d.csv", i)
				err  = ioutil.WriteFile(path, []byte(tc.str), 0644)
			)
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
			dataset, err := ReadCSV(path, tc.target, tc.classification)
			if err != nil {
				t.Errorf("Expected nil, got %s", err)
			}
			if !reflect.DeepEqual(dataset, tc.dataset) {
				t.Errorf("Expected\n%s, got\n%s", tc.dataset, dataset)
			}
		})
	}
}
