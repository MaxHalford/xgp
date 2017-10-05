package xgp

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNodeSimplify(t *testing.T) {
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
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.node.Simplify()
			if !reflect.DeepEqual(tc.node, tc.prunedNode) {
				t.Errorf("Expected %v, got %v", tc.prunedNode, tc.node)
			}
		})
	}
}
