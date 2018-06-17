package op

import (
	"fmt"
	"testing"
)

func TestMinSimplify(t *testing.T) {
	var testCases = []struct {
		in  Min
		out Operator
	}{
		{
			in:  Min{Mul{Add{Const{1}, Const{2}}, Var{0}}, Const{2}},
			out: Min{Mul{Const{3}, Var{0}}, Const{2}},
		},
		{
			in:  Min{Const{2}, Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Min{Const{2}, Mul{Const{3}, Var{0}}},
		},
		{
			in:  Min{Const{2}, Const{3}},
			out: Const{2},
		},
		{
			in:  Min{Var{2}, Const{3}},
			out: Min{Var{2}, Const{3}},
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
