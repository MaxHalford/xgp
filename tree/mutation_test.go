package tree

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestModification(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			mutator Mutator
		}{
			{
				mutator: HoistMutation{
					Picker: WeightedPicker{
						PConstant: 0.05,
						PVariable: 0.05,
						PFunction: 0.9,
					},
				},
			},
			{
				mutator: SubTreeMutation{
					Picker: WeightedPicker{
						PConstant: 0.05,
						PVariable: 0.05,
						PFunction: 0.9,
					},
					NewTree: func(minHeight, maxHeight int, rng *rand.Rand) *Tree {
						return &Tree{Operator: Constant{42}}
					},
				},
			},
			{
				mutator: PointMutation{
					Picker: WeightedPicker{
						PConstant: 1,
						PVariable: 1,
						PFunction: 0,
					},
					MutateOperator: func(op Operator, rng *rand.Rand) Operator {
						return Constant{42}
					},
				},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				tree  = randTree(rng)
				clone = tree.Clone()
			)
			tc.mutator.Apply(clone, rng)
			if reflect.DeepEqual(clone, tree) {
				t.Error("Mutation should not have affected the original tree")
			}
		})
	}
}

func TestHoistMutation(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in      *Tree
			mutator HoistMutation
			out     *Tree
		}{
			{
				in: &Tree{
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
				},
				mutator: HoistMutation{
					Picker: WeightedPicker{
						PConstant: 0.05,
						PVariable: 0.05,
						PFunction: 0.9,
					},
				},
				out: &Tree{
					Operator: Constant{1},
					Branches: []*Tree{
						&Tree{Operator: Constant{3}},
					},
				},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.mutator.Apply(tc.in, rng)
			if !reflect.DeepEqual(tc.in, tc.out) {
				t.Errorf("Expected %s, got %s", tc.out, tc.in)
			}
		})
	}
}
