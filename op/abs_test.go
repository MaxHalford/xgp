package op

import (
	"fmt"
	"testing"
)

func TestAbsSimplify(t *testing.T) {
	var testCases = []struct {
		in  Abs
		out Operator
	}{
		{
			in:  Abs{Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Abs{Mul{Const{3}, Var{0}}},
		},
		{
			in:  Abs{Abs{Var{1}}},
			out: Abs{Var{1}},
		},
		{
			in:  Abs{Neg{Var{1}}},
			out: Var{1},
		},
		{
			in:  Abs{Const{-42}},
			out: Const{42},
		},
		{
			in:  Abs{Var{0}},
			out: Abs{Var{0}},
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
