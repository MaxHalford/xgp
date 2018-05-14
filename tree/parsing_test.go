package tree

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MaxHalford/xgp/op"
)

func TestParseCode(t *testing.T) {
	var testCases = []struct {
		code string
		tree Tree
	}{
		{
			code: "sum(X[0], 42)",
			tree: Tree{
				op: op.Sum{},
				branches: []*Tree{
					&Tree{op: op.Variable{0}},
					&Tree{op: op.Constant{42}},
				},
			},
		},
		{
			code: "cos(sum(X[0], 42))",
			tree: Tree{
				op: op.Cos{},
				branches: []*Tree{
					&Tree{
						op: op.Sum{},
						branches: []*Tree{
							&Tree{op: op.Variable{0}},
							&Tree{op: op.Constant{42}},
						},
					},
				},
			},
		},
		{
			code: "sum(sum(X[0], 42), sum(X[1], 43))",
			tree: Tree{
				op: op.Sum{},
				branches: []*Tree{
					&Tree{
						op: op.Sum{},
						branches: []*Tree{
							&Tree{op: op.Variable{0}},
							&Tree{op: op.Constant{42}},
						},
					},
					&Tree{
						op: op.Sum{},
						branches: []*Tree{
							&Tree{op: op.Variable{1}},
							&Tree{op: op.Constant{43}},
						},
					},
				},
			},
		},
		{
			code: "cos(cos(cos(42)))",
			tree: Tree{
				op: op.Cos{},
				branches: []*Tree{
					&Tree{
						op: op.Cos{},
						branches: []*Tree{
							&Tree{
								op: op.Cos{},
								branches: []*Tree{
									&Tree{
										op: op.Constant{42},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			code: "mul(cos(X[0]), exp(sin(X[1])))",
			tree: Tree{
				op: op.Mul{},
				branches: []*Tree{
					&Tree{
						op: op.Cos{},
						branches: []*Tree{
							&Tree{
								op: op.Variable{0},
							},
						},
					},
					&Tree{
						op: op.Exp{},
						branches: []*Tree{
							&Tree{
								op: op.Sin{},
								branches: []*Tree{
									&Tree{
										op: op.Variable{1},
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
