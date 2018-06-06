package xgp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/op"
)

func TestFullInit(t *testing.T) {
	var (
		newOp = func(leaf bool, rng *rand.Rand) op.Operator {
			if leaf {
				return op.Const{1}
			}
			return op.Add{}
		}
		rng       = newRand()
		testCases = []struct {
			maxHeight uint
			minHeight uint
			height    uint
			size      uint
		}{
			{
				minHeight: 0,
				maxHeight: 0,
				height:    0,
				size:      1,
			},
			{
				minHeight: 1,
				maxHeight: 1,
				height:    1,
				size:      3,
			},
			{
				minHeight: 2,
				maxHeight: 2,
				height:    2,
				size:      7,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			operator := FullInit{}.Apply(tc.minHeight, tc.maxHeight, newOp, rng)
			h := op.CalcHeight(operator)
			if h != tc.height {
				t.Errorf("Expected height %d, got %d", tc.height, h)
			}
			s := op.CountOps(operator)
			if s != tc.size {
				t.Errorf("Expected size %d, got %d", tc.size, s)
			}
		})
	}
}

func TestGrowInit(t *testing.T) {
	var (
		newOp = func(leaf bool, rng *rand.Rand) op.Operator {
			if leaf {
				return op.Const{1}
			}
			return op.Add{}
		}
		rng       = newRand()
		testCases = []struct {
			pLeaf     float64
			minHeight uint
			maxHeight uint
			height    uint
			size      uint
		}{
			{
				pLeaf:     0,
				minHeight: 0,
				maxHeight: 0,
				height:    0,
				size:      1,
			},
			{
				pLeaf:     1,
				minHeight: 0,
				maxHeight: 0,
				height:    0,
				size:      1,
			},
			{
				pLeaf:     0,
				minHeight: 1,
				maxHeight: 1,
				height:    1,
				size:      3,
			},
			{
				pLeaf:     1,
				minHeight: 0,
				maxHeight: 1,
				height:    0,
				size:      1,
			},
			{
				pLeaf:     0,
				minHeight: 2,
				maxHeight: 2,
				height:    2,
				size:      7,
			},
			{
				pLeaf:     1,
				minHeight: 0,
				maxHeight: 2,
				height:    0,
				size:      1,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			operator := GrowInit{tc.pLeaf}.Apply(tc.minHeight, tc.maxHeight, newOp, rng)
			h := op.CalcHeight(operator)
			if h != tc.height {
				t.Errorf("Expected height %d, got %d", tc.height, h)
			}
			s := op.CountOps(operator)
			if s != tc.size {
				t.Errorf("Expected size %d, got %d", tc.size, s)
			}
		})
	}
}
