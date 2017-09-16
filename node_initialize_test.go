package xgp

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/MaxHalford/xgp/tree"
)

func TestFullNodeInitializer(t *testing.T) {
	var (
		newOperator = func(terminal bool, rng *rand.Rand) Operator {
			if terminal {
				return Constant{1}
			}
			return Sum{}
		}
		rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
		testCases = []struct {
			height int
			nnodes int
		}{
			{
				height: 0,
				nnodes: 1,
			},
			{
				height: 1,
				nnodes: 3,
			},
			{
				height: 2,
				nnodes: 7,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var node = FullNodeInitializer{Height: tc.height}.Apply(newOperator, rng)
			if tree.GetNNodes(node) != tc.nnodes {
				t.Errorf("Expected %d nodes, got %d", tc.nnodes, tree.GetNNodes(node))
			}
		})
	}

}

func TestGrowNodeInitializer(t *testing.T) {
	var (
		newOperator = func(terminal bool, rng *rand.Rand) Operator {
			if terminal {
				return Constant{1}
			}
			return Sum{}
		}
		rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
		testCases = []struct {
			pLeaf     float64
			maxHeight int
			nnodes    int
		}{
			{
				pLeaf:     0,
				maxHeight: 0,
				nnodes:    1,
			},
			{
				pLeaf:     1,
				maxHeight: 0,
				nnodes:    1,
			},
			{
				pLeaf:     0,
				maxHeight: 1,
				nnodes:    3,
			},
			{
				pLeaf:     1,
				maxHeight: 1,
				nnodes:    1,
			},
			{
				pLeaf:     0,
				maxHeight: 2,
				nnodes:    7,
			},
			{
				pLeaf:     1,
				maxHeight: 2,
				nnodes:    1,
			},
		}
	)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				initializer = GrowNodeInitializer{MaxHeight: tc.maxHeight, PLeaf: tc.pLeaf}
				node        = initializer.Apply(newOperator, rng)
			)
			if tree.GetNNodes(node) != tc.nnodes {
				t.Errorf("Expected %d nodes, got %d", tc.nnodes, tree.GetNNodes(node))
			}
		})
	}
}
