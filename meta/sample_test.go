package meta

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRandIntsWeighted(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			k       int
			weights []float64
			ints    []int
		}{
			{
				k:       3,
				weights: []float64{1, 0, 0},
				ints:    []int{0, 0, 0},
			},
			{
				k:       3,
				weights: []float64{0, 1, 0},
				ints:    []int{1, 1, 1},
			},
			{
				k:       3,
				weights: []float64{0, 0, 1},
				ints:    []int{2, 2, 2},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var ints = randIntsWeighted(tc.k, tc.weights, rng)
			if !reflect.DeepEqual(ints, tc.ints) {
				t.Errorf("Expected %v, got %v", tc.ints, ints)
			}
		})
	}
}
