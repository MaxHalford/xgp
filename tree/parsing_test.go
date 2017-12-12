package tree

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MaxHalford/koza/tree/op"
)

func TestParseCode(t *testing.T) {
	var testCases = []struct {
		code string
		tree *Tree
	}{
		{
			code: "sum(X[0], 42)",
			tree: &Tree{
				Operator: op.Sum{},
				Branches: []*Tree{
					&Tree{Operator: op.Variable{0}},
					&Tree{Operator: op.Constant{42}},
				},
			},
		},
		{
			code: "cos(sum(X[0], 42))",
			tree: &Tree{
				Operator: op.Cos{},
				Branches: []*Tree{
					&Tree{
						Operator: op.Sum{},
						Branches: []*Tree{
							&Tree{Operator: op.Variable{0}},
							&Tree{Operator: op.Constant{42}},
						},
					},
				},
			},
		},
		{
			code: "sum(sum(X[0], 42), sum(X[1], 43))",
			tree: &Tree{
				Operator: op.Sum{},
				Branches: []*Tree{
					&Tree{
						Operator: op.Sum{},
						Branches: []*Tree{
							&Tree{Operator: op.Variable{0}},
							&Tree{Operator: op.Constant{42}},
						},
					},
					&Tree{
						Operator: op.Sum{},
						Branches: []*Tree{
							&Tree{Operator: op.Variable{1}},
							&Tree{Operator: op.Constant{43}},
						},
					},
				},
			},
		},
		{
			code: "cos(cos(cos(42)))",
			tree: &Tree{
				Operator: op.Cos{},
				Branches: []*Tree{
					&Tree{
						Operator: op.Cos{},
						Branches: []*Tree{
							&Tree{
								Operator: op.Cos{},
								Branches: []*Tree{
									&Tree{
										Operator: op.Constant{42},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			code: "mul(cos(X[0]), log(sin(X[1])))",
			tree: &Tree{
				Operator: op.Product{},
				Branches: []*Tree{
					&Tree{
						Operator: op.Cos{},
						Branches: []*Tree{
							&Tree{
								Operator: op.Variable{0},
							},
						},
					},
					&Tree{
						Operator: op.Log{},
						Branches: []*Tree{
							&Tree{
								Operator: op.Sin{},
								Branches: []*Tree{
									&Tree{
										Operator: op.Variable{1},
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
