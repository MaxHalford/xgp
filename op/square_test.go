package op

import (
	"fmt"
	"testing"
)

func TestSquareSimplify(t *testing.T) {
	var testCases = []struct {
		in  Square
		out Operator
	}{
		{
			in:  Square{Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Square{Mul{Const{3}, Var{0}}},
		},
		{
			in:  Square{Const{3}},
			out: Const{9},
		},
		{
			in:  Square{Neg{Var{0}}},
			out: Square{Var{0}},
		},
		{
			in:  Square{Var{0}},
			out: Square{Var{0}},
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
