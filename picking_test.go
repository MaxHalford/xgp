package xgp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MaxHalford/xgp/tree"
)

func TestWeightedPicker(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			w         Weighting
			minHeight int
			maxHeight int
			in        tree.Tree
			out       tree.Tree
		}{
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 0,
					PFunction: 0,
				},
				minHeight: 0,
				maxHeight: 1,
				in:        tree.MustParseCode("sum(42, X[0])"),
				out:       tree.MustParseCode("42"),
			},
			{
				w: Weighting{
					PConstant: 0,
					PVariable: 1,
					PFunction: 0,
				},
				minHeight: 0,
				maxHeight: 1,
				in:        tree.MustParseCode("sum(42, X[0])"),
				out:       tree.MustParseCode("X[0]"),
			},
			{
				w: Weighting{
					PConstant: 0,
					PVariable: 0,
					PFunction: 1,
				},
				minHeight: 0,
				maxHeight: 1,
				in:        tree.MustParseCode("sum(42, X[0])"),
				out:       tree.MustParseCode("sum(42, X[0])"),
			},
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 1,
					PFunction: 1,
				},
				minHeight: 0,
				maxHeight: 0,
				in:        tree.MustParseCode("cos(sin(42))"),
				out:       tree.MustParseCode("42"),
			},
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 1,
					PFunction: 1,
				},
				minHeight: 1,
				maxHeight: 1,
				in:        tree.MustParseCode("cos(sin(42))"),
				out:       tree.MustParseCode("sin(42)"),
			},
			{
				w: Weighting{
					PConstant: 1,
					PVariable: 1,
					PFunction: 1,
				},
				minHeight: 2,
				maxHeight: 2,
				in:        tree.MustParseCode("cos(sin(42))"),
				out:       tree.MustParseCode("cos(sin(42))"),
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				wp  = WeightedPicker{Weighting: tc.w}
				out = wp.Apply(&tc.in, tc.minHeight, tc.maxHeight, rng)
			)
			if !reflect.DeepEqual(out, &tc.out) {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}

}
