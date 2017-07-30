package xgp

import (
	"reflect"
	"testing"
)

func TestNodePrune(t *testing.T) {
	var testCases = []struct {
		node       *Node
		prunedNode *Node
	}{
		{
			node: &Node{
				Operator: Sum{},
				Children: []*Node{
					&Node{Operator: Constant{1}},
					&Node{Operator: Constant{2}},
				},
			},
			prunedNode: &Node{
				Operator: Constant{3},
			},
		},
		{
			node: &Node{
				Operator: Sum{},
				Children: []*Node{
					&Node{Operator: Variable{0}},
					&Node{Operator: Constant{42}},
				},
			},
			prunedNode: &Node{
				Operator: Sum{},
				Children: []*Node{
					&Node{Operator: Variable{0}},
					&Node{Operator: Constant{42}},
				},
			},
		},
	}
	for i, tc := range testCases {
		tc.node.Prune()
		if !reflect.DeepEqual(tc.node, tc.prunedNode) {
			t.Errorf("Error in test case %d", i)
		}
	}
}
