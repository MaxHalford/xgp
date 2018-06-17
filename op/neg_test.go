package op

import (
	"fmt"
	"testing"
)

func TestNegSimplify(t *testing.T) {
	var testCases = []struct {
		in  Neg
		out Operator
	}{
		{
			in:  Neg{Add{Mul{Const{1}, Const{-2}}, Var{0}}},
			out: Neg{Add{Const{-2}, Var{0}}},
		},
		{
			in:  Neg{Const{0}},
			out: Const{0},
		},
		{
			in:  Neg{Neg{Const{42}}},
			out: Const{42},
		},
		{
			in:  Neg{Add{Const{42}, Const{43}}},
			out: Const{-85},
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
