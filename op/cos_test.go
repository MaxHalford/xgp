package op

import (
	"fmt"
	"math"
	"testing"
)

func TestCosSimplify(t *testing.T) {
	var testCases = []struct {
		in  Cos
		out Operator
	}{
		{
			in:  Cos{Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Cos{Mul{Const{3}, Var{0}}},
		},
		{
			in:  Cos{Const{math.Pi}},
			out: Const{-1},
		},
		{
			in:  Cos{Var{0}},
			out: Cos{Var{0}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := tc.in.Simplify()
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}
