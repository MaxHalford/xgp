package tree

import (
	"testing"
)

func TestHoistMutation(t *testing.T) {
	var (
		tree = &Tree{
			Operator: Constant{1},
			Branches: []*Tree{
				&Tree{
					Operator: Constant{2},
					Branches: []*Tree{
						&Tree{
							Operator: Constant{3},
						},
					},
				},
			},
		}
	)
	HoistMutation{}.Apply(tree, newRand())
}
