package xgp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/op"
)

func TestSubtreeCrossover(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in1  op.Operator
			in2  op.Operator
			sc   SubtreeCrossover
			out1 op.Operator
			out2 op.Operator
		}{
			{
				in1: op.Cos{op.Var{0}},
				in2: op.Sin{op.Var{1}},
				sc: SubtreeCrossover{
					Weight: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						switch operator.(type) {
						case op.Var:
							return 1
						default:
							return 0
						}
					},
				},
				out1: op.Cos{op.Var{1}},
				out2: op.Sin{op.Var{0}},
			},
			{
				in1: op.Add{op.Cos{op.Var{0}}, op.Sin{op.Const{0}}},
				in2: op.Add{op.Cos{op.Var{0}}, op.Sin{op.Const{1}}},
				sc: SubtreeCrossover{
					Weight: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						switch operator.(type) {
						case op.Const:
							return 1
						default:
							return 0
						}
					},
				},
				out1: op.Add{op.Cos{op.Var{0}}, op.Sin{op.Const{1}}},
				out2: op.Add{op.Cos{op.Var{0}}, op.Sin{op.Const{0}}},
			},
			{
				in1: op.Cos{op.Var{0}},
				in2: op.Sin{op.Const{1}},
				sc: SubtreeCrossover{
					Weight: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
						switch operator.(type) {
						case op.Const:
							return 0
						case op.Var:
							return 0
						default:
							return 1
						}
					},
				},
				out1: op.Sin{op.Const{1}},
				out2: op.Cos{op.Var{0}},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out1, out2 := tc.sc.Apply(tc.in1, tc.in2, rng)
			if out1 != tc.out1 {
				t.Errorf("Expected %s, got %s", tc.out1, out1)
			}
			if out2 != tc.out2 {
				t.Errorf("Expected %s, got %s", tc.out2, out2)
			}
		})
	}
}
