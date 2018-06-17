package op

import (
	"fmt"
	"testing"
)

func TestInvSimplify(t *testing.T) {
	var testCases = []struct {
		in  Inv
		out Operator
	}{
		{
			in:  Inv{Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Inv{Mul{Const{3}, Var{0}}},
		},
		{
			in:  Inv{Inv{Var{42}}},
			out: Var{42},
		},
		{
			in:  Inv{Const{42}},
			out: Const{1.0 / 42},
		},
		{
			in:  Inv{Var{0}},
			out: Inv{Var{0}},
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
