package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseCode(t *testing.T) {
	var testCases = []struct {
		code string
		tree *Tree
	}{
		{
			code: "sum(X[0], 42)",
			tree: &Tree{
				Operator: Sum{},
				Branches: []*Tree{
					&Tree{Operator: Variable{0}},
					&Tree{Operator: Constant{42}},
				},
			},
		},
		{
			code: "cos(sum(X[0], 42))",
			tree: &Tree{
				Operator: Cos{},
				Branches: []*Tree{
					&Tree{
						Operator: Sum{},
						Branches: []*Tree{
							&Tree{Operator: Variable{0}},
							&Tree{Operator: Constant{42}},
						},
					},
				},
			},
		},
		{
			code: "sum(sum(X[0], 42), sum(X[1], 43))",
			tree: &Tree{
				Operator: Sum{},
				Branches: []*Tree{
					&Tree{
						Operator: Sum{},
						Branches: []*Tree{
							&Tree{Operator: Variable{0}},
							&Tree{Operator: Constant{42}},
						},
					},
					&Tree{
						Operator: Sum{},
						Branches: []*Tree{
							&Tree{Operator: Variable{1}},
							&Tree{Operator: Constant{43}},
						},
					},
				},
			},
		},
		{
			code: "cos(cos(cos(42)))",
			tree: &Tree{
				Operator: Cos{},
				Branches: []*Tree{
					&Tree{
						Operator: Cos{},
						Branches: []*Tree{
							&Tree{
								Operator: Cos{},
								Branches: []*Tree{
									&Tree{
										Operator: Constant{42},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var tree, err = ParseCode(tc.code)
			if err != nil {
				t.Errorf("%s", err)
			}
			if !reflect.DeepEqual(tc.tree, tree) {
				t.Errorf("Expected %v, got %v", tc.tree, tree)
			}
		})
	}
}
