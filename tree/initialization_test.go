package tree

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestFullInitializer(t *testing.T) {
	var (
		of = OperatorFactory{
			PConstant:   1,
			NewConstant: func(rng *rand.Rand) Constant { return Constant{1} },
			NewFunction: func(rng *rand.Rand) Operator { return Sum{} },
		}
		rng       = newRand()
		testCases = []struct {
			maxHeight  int
			minHeight  int
			nOperators int
		}{
			{
				minHeight:  0,
				maxHeight:  0,
				nOperators: 1,
			},
			{
				minHeight:  1,
				maxHeight:  1,
				nOperators: 3,
			},
			{
				minHeight:  2,
				maxHeight:  2,
				nOperators: 7,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var tree = FullInitializer{}.Apply(tc.minHeight, tc.maxHeight, of, rng)
			if tree.NOperators() != tc.nOperators {
				t.Errorf("Expected %d trees, got %d", tc.nOperators, tree.NOperators())
			}
		})
	}

}

func TestGrowInitializer(t *testing.T) {
	var (
		of = OperatorFactory{
			PConstant:   1,
			NewConstant: func(rng *rand.Rand) Constant { return Constant{1} },
			NewFunction: func(rng *rand.Rand) Operator { return Sum{} },
		}
		rng       = newRand()
		testCases = []struct {
			minHeight  int
			maxHeight  int
			pLeaf      float64
			nOperators int
		}{
			{
				pLeaf:      0,
				minHeight:  0,
				maxHeight:  0,
				nOperators: 1,
			},
			{
				pLeaf:      1,
				minHeight:  0,
				maxHeight:  0,
				nOperators: 1,
			},
			{
				pLeaf:      0,
				minHeight:  1,
				maxHeight:  1,
				nOperators: 3,
			},
			{
				pLeaf:      1,
				minHeight:  0,
				maxHeight:  1,
				nOperators: 1,
			},
			{
				pLeaf:      0,
				minHeight:  2,
				maxHeight:  2,
				nOperators: 7,
			},
			{
				pLeaf:      1,
				minHeight:  0,
				maxHeight:  2,
				nOperators: 1,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				initializer = GrowInitializer{
					PTerminal: tc.pLeaf,
				}
				tree = initializer.Apply(tc.minHeight, tc.maxHeight, of, rng)
			)
			if tree.NOperators() != tc.nOperators {
				t.Errorf("Expected %d operator(s), got %d", tc.nOperators, tree.NOperators())
			}
		})
	}
}
