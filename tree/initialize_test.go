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
			height     int
			nOperators int
		}{
			{
				height:     0,
				nOperators: 1,
			},
			{
				height:     1,
				nOperators: 3,
			},
			{
				height:     2,
				nOperators: 7,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var tree = FullInitializer{Height: tc.height}.Apply(of, rng)
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
			pLeaf      float64
			maxHeight  int
			nOperators int
		}{
			{
				pLeaf:      0,
				maxHeight:  0,
				nOperators: 1,
			},
			{
				pLeaf:      1,
				maxHeight:  0,
				nOperators: 1,
			},
			{
				pLeaf:      0,
				maxHeight:  1,
				nOperators: 3,
			},
			{
				pLeaf:      1,
				maxHeight:  1,
				nOperators: 1,
			},
			{
				pLeaf:      0,
				maxHeight:  2,
				nOperators: 7,
			},
			{
				pLeaf:      1,
				maxHeight:  2,
				nOperators: 1,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				initializer = GrowInitializer{MaxHeight: tc.maxHeight, PLeaf: tc.pLeaf}
				tree        = initializer.Apply(of, rng)
			)
			if tree.NOperators() != tc.nOperators {
				t.Errorf("Expected %d trees, got %d", tc.nOperators, tree.NOperators())
			}
		})
	}
}
