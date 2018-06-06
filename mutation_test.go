package xgp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/op"
)

func TestPointMutation(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in  op.Operator
			pm  PointMutation
			out op.Operator
		}{
			{
				in: op.Sub{op.Var{0}, op.Sin{op.Const{1}}},
				pm: PointMutation{
					Rate: 1,
					Mutate: func(operator op.Operator, rng *rand.Rand) op.Operator {
						switch operator.Arity() {
						case 0:
							return op.Const{42}
						case 1:
							return op.Cos{operator.Operand(0)}
						}
						return op.Add{operator.Operand(0), operator.Operand(1)}
					},
				},
				out: op.Add{op.Const{42}, op.Cos{op.Const{42}}},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.pm.Apply(tc.in, rng)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestHoistMutation(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in  op.Operator
			hm  HoistMutation
			out op.Operator
		}{
			{
				in: op.Const{42},
				hm: HoistMutation{
					Weight1: func(operator op.Operator, depth uint, rng *rand.Rand) float64 { return 1 },
					Weight2: func(operator op.Operator, depth uint, rng *rand.Rand) float64 { return 1 },
				},
				out: op.Const{42},
			},
			{
				in: op.Cos{op.Sin{op.Const{42}}},
				hm: HoistMutation{
					Weight1: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						if depth == 0 {
							return 1
						}
						return 0
					},
					Weight2: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						if depth != 2 {
							return 0
						}
						return 1
					},
				},
				out: op.Const{42},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.hm.Apply(tc.in, rng)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestSubtreeMutation(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in  op.Operator
			sm  SubtreeMutation
			out op.Operator
		}{
			{
				in: op.Const{42},
				sm: SubtreeMutation{
					func(operator op.Operator, depth uint, rng *rand.Rand) float64 { return 1 },
					func(rng *rand.Rand) op.Operator { return op.Cos{op.Var{7}} },
				},
				out: op.Cos{op.Var{7}},
			},
			{
				in: op.Sin{op.Cos{op.Const{42}}},
				sm: SubtreeMutation{
					func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						if operator.Arity() > 0 {
							return 0
						}
						return 1
					},
					func(rng *rand.Rand) op.Operator { return op.Var{1} },
				},
				out: op.Sin{op.Cos{op.Var{1}}},
			},
			{
				in: op.Sin{op.Cos{op.Const{42}}},
				sm: SubtreeMutation{
					func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						if depth != 1 {
							return 0
						}
						return 1
					},
					func(rng *rand.Rand) op.Operator { return op.Var{1} },
				},
				out: op.Sin{op.Var{1}},
			},
			{
				in: op.Sin{op.Cos{op.Const{42}}},
				sm: SubtreeMutation{
					func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						if depth != 0 {
							return 0
						}
						return 1
					},
					func(rng *rand.Rand) op.Operator { return op.Var{1} },
				},
				out: op.Var{1},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.sm.Apply(tc.in, rng)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}
