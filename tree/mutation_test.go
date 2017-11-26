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
						Weighting: Weighting{
							PConstant: 0.05,
							PVariable: 0.05,
							PFunction: 0.9,
						},
					},
				},
			},
			{
				mutator: SubTreeMutation{
					Crossover: SubTreeCrossover{
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
				tree   = randTree(rng)
				mutant = tree.Clone()
			)
			tc.mutator.Apply(mutant, rng)
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
			in      *Tree
			mutator HoistMutation
			out     *Tree
		}{
			{
				in: mustParseCode("cos(sin(42))"),
				mutator: HoistMutation{
					Picker: WeightedPicker{
						Weighting: Weighting{
							PConstant: 0.05,
							PVariable: 0.05,
							PFunction: 0.9,
						},
					},
				},
				out: mustParseCode("cos(42)"),
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
