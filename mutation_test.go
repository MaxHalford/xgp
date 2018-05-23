package xgp

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
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
						Weighting: Weighting{
							PConstant: 0.05,
							PVariable: 0.05,
							PFunction: 0.9,
						},
					},
				},
			},
			{
				mutator: SubtreeMutation{
					Crossover: SubtreeCrossover{
						Picker: WeightedPicker{
							Weighting: Weighting{
								PConstant: 0.1,
								PVariable: 0.1,
								PFunction: 0.1,
							},
						},
					},
					NewTree: randTree,
				},
			},
			{
				mutator: PointMutation{
					Weighting: Weighting{
						PConstant: 1,
						PVariable: 1,
						PFunction: 0,
					},
					MutateOperator: func(op.Operator, *rand.Rand) op.Operator {
						return op.Constant{42}
					},
				},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				tree   = randTree(rng)
				mutant = tree.Clone()
			)
			tc.mutator.Apply(&mutant, rng)
			if reflect.DeepEqual(mutant, tree) {
				t.Error("Mutation did not make any difference")
			}
		})
	}
}

func TestHoistMutation(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in      tree.Tree
			mutator HoistMutation
			out     tree.Tree
		}{
			{
				in: tree.MustParseCode("42"),
				mutator: HoistMutation{
					Picker: WeightedPicker{
						Weighting: Weighting{
							PConstant: 0.05,
							PVariable: 0.05,
							PFunction: 0.9,
						},
					},
				},
				out: tree.MustParseCode("42"),
			},
			{
				in: tree.MustParseCode("cos(42)"),
				mutator: HoistMutation{
					Picker: WeightedPicker{
						Weighting: Weighting{
							PConstant: 0.05,
							PVariable: 0.05,
							PFunction: 0.9,
						},
					},
				},
				out: tree.MustParseCode("42"),
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.mutator.Apply(&tc.in, rng)
			if !reflect.DeepEqual(tc.in, tc.out) {
				t.Errorf("Expected %s, got %s", tc.out, tc.in)
			}
		})
	}
}
