package xgp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/op"
)

func TestFullInitializer(t *testing.T) {
	var (
		of = OperatorFactory{
			PConstant:   1,
			NewConstant: func(rng *rand.Rand) op.Constant { return op.Constant{1} },
			NewFunction: func(rng *rand.Rand) op.Operator { return op.Sum{} },
		}
		rng       = newRand()
		testCases = []struct {
			maxHeight int
			minHeight int
			size      int
		}{
			{
				minHeight: 0,
				maxHeight: 0,
				size:      1,
			},
			{
				minHeight: 1,
				maxHeight: 1,
				size:      3,
			},
			{
				minHeight: 2,
				maxHeight: 2,
				size:      7,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var tree = FullInitializer{}.Apply(tc.minHeight, tc.maxHeight, of, rng)
			if tree.Size() != tc.size {
				t.Errorf("Expected %d, got %d", tc.size, tree.Size())
			}
		})
	}

}

func TestGrowInitializer(t *testing.T) {
	var (
		of = OperatorFactory{
			PConstant:   1,
			NewConstant: func(rng *rand.Rand) op.Constant { return op.Constant{1} },
			NewFunction: func(rng *rand.Rand) op.Operator { return op.Sum{} },
		}
		rng       = newRand()
		testCases = []struct {
			minHeight int
			maxHeight int
			pLeaf     float64
			size      int
		}{
			{
				pLeaf:     0,
				minHeight: 0,
				maxHeight: 0,
				size:      1,
			},
			{
				pLeaf:     1,
				minHeight: 0,
				maxHeight: 0,
				size:      1,
			},
			{
				pLeaf:     0,
				minHeight: 1,
				maxHeight: 1,
				size:      3,
			},
			{
				pLeaf:     1,
				minHeight: 0,
				maxHeight: 1,
				size:      1,
			},
			{
				pLeaf:     0,
				minHeight: 2,
				maxHeight: 2,
				size:      7,
			},
			{
				pLeaf:     1,
				minHeight: 0,
				maxHeight: 2,
				size:      1,
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
			if tree.Size() != tc.size {
				t.Errorf("Expected %d operator(s), got %d", tc.size, tree.Size())
			}
		})
	}
}
