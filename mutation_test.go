package xgp

import (
	"testing"
)

func TestHoistMutation(t *testing.T) {
	var (
		prog = Program{
			Root: &Node{
				Operator: Constant{1},
				Children: []*Node{
					&Node{
						Operator: Constant{2},
						Children: []*Node{
							&Node{
								Operator: Constant{3},
							},
						},
					},
				},
			},
		}
	)
	HoistMutation{}.Apply(&prog, makeRNG())
}
