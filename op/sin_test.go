package op

import (
	"fmt"
	"math"
	"testing"
)

func TestSinSimplify(t *testing.T) {
	var testCases = []struct {
		in  Sin
		out Operator
	}{
		{
			in:  Sin{Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Sin{Mul{Const{3}, Var{0}}},
		},
		{
			in:  Sin{Const{0.5 * math.Pi}},
			out: Const{1},
		},
		{
			in:  Sin{Var{0}},
			out: Sin{Var{0}},
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
