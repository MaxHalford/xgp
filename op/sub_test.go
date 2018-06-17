package op

import (
	"fmt"
	"testing"
)

func TestSubSimplify(t *testing.T) {
	var testCases = []struct {
		in  Sub
		out Operator
	}{
		{
			in:  Sub{Mul{Const{1}, Const{-2}}, Var{0}},
			out: Sub{Const{-2}, Var{0}},
		},
		{
			in:  Sub{Var{0}, Mul{Const{1}, Const{-2}}},
			out: Sub{Var{0}, Const{-2}},
		},
		{
			in:  Sub{Var{1}, Const{0}},
			out: Var{1},
		},
		{
			in:  Sub{Const{0}, Var{1}},
			out: Neg{Var{1}},
		},
		{
			in:  Sub{Const{1}, Const{2}},
			out: Const{-1},
		},
		{
			in:  Sub{Var{1}, Var{1}},
			out: Const{0},
		},
		{
			in:  Sub{Var{0}, Var{1}},
			out: Sub{Var{0}, Var{1}},
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
