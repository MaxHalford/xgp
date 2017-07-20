package gp

import (
	"math/rand"
	"testing"
	"time"

	"github.com/MaxHalford/tiago/tree"
)

func TestFullNodeInitializer(t *testing.T) {
	var (
		newNode = func(leaf bool, rng *rand.Rand) *Node {
			var node = &Node{}
			if !leaf {
				node.Children = make([]*Node, 2)
			}
			return node
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

	for _, tc := range testCases {
		var node = FullNodeInitializer{Height: tc.height}.Apply(newNode, rng)
		if tree.GetNNodes(node) != tc.nnodes {
			t.Errorf("Expected %d nodes, got %d", tc.nnodes, tree.GetNNodes(node))
		}
	}

}

func TestGrowNodeInitializer(t *testing.T) {
	var (
		newNode = func(leaf bool, rng *rand.Rand) *Node {
			var node = &Node{}
			if !leaf {
				node.Children = make([]*Node, 2)
			}
			return node
		}
		rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
		testCases = []struct {
			pLeaf    float64
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

	for _, tc := range testCases {
		var (
			initializer = GrowNodeInitializer{MaxHeight: tc.maxHeight, PLeaf: tc.pLeaf}
			node        = initializer.Apply(newNode, rng)
		)
		if tree.GetNNodes(node) != tc.nnodes {
			t.Errorf("Expected %d nodes, got %d", tc.nnodes, tree.GetNNodes(node))
		}
	}

}
