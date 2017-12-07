package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWeightedPicker(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			w         Weighting
			minHeight int
			maxHeight int
			in        *Tree
			out       *Tree
		}{
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 0,
					PFunction: 0,
				},
				minHeight: 0,
				maxHeight: 1,
				in:        mustParseCode("sum(42, X[0])"),
				out:       mustParseCode("42"),
			},
			{
				w: Weighting{
					PConstant: 0,
					PVariable: 1,
					PFunction: 0,
				},
				minHeight: 0,
				maxHeight: 1,
				in:        mustParseCode("sum(42, X[0])"),
				out:       mustParseCode("X[0]"),
			},
			{
				w: Weighting{
					PConstant: 0,
					PVariable: 0,
					PFunction: 1,
				},
				minHeight: 0,
				maxHeight: 1,
				in:        mustParseCode("sum(42, X[0])"),
				out:       mustParseCode("sum(42, X[0])"),
			},
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 1,
					PFunction: 1,
				},
				minHeight: 0,
				maxHeight: 0,
				in:        mustParseCode("cos(sin(42))"),
				out:       mustParseCode("42"),
			},
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 1,
					PFunction: 1,
				},
				minHeight: 1,
				maxHeight: 1,
				in:        mustParseCode("cos(sin(42))"),
				out:       mustParseCode("sin(42)"),
			},
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 1,
					PFunction: 1,
				},
				minHeight: 2,
				maxHeight: 2,
				in:        mustParseCode("cos(sin(42))"),
				out:       mustParseCode("cos(sin(42))"),
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				wp  = WeightedPicker{Weighting: tc.w}
				out = wp.Apply(tc.in, tc.minHeight, tc.maxHeight, rng)
			)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}

}
