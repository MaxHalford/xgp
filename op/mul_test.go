package op

import (
	"fmt"
	"testing"
)

func TestMulSimplify(t *testing.T) {
	var testCases = []struct {
		in  Mul
		out Operator
	}{
		{
			in:  Mul{Mul{Add{Const{1}, Const{2}}, Var{0}}, Const{2}},
			out: Mul{Mul{Const{3}, Var{0}}, Const{2}},
		},
		{
			in:  Mul{Const{2}, Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Mul{Const{2}, Mul{Const{3}, Var{0}}},
		},
		{
			in:  Mul{Const{0}, Var{42}},
			out: Const{0},
		},
		{
			in:  Mul{Var{42}, Const{0}},
			out: Const{0},
		},
		{
			in:  Mul{Const{1}, Var{42}},
			out: Var{42},
		},
		{
			in:  Mul{Var{42}, Const{1}},
			out: Var{42},
		},
		{
			in:  Mul{Const{-1}, Var{42}},
			out: Neg{Var{42}},
		},
		{
			in:  Mul{Var{42}, Const{-1}},
			out: Neg{Var{42}},
		},
		{
			in:  Mul{Const{-42}, Const{-1}},
			out: Const{42},
		},
		{
			in:  Mul{Var{0}, Var{1}},
			out: Mul{Var{0}, Var{1}},
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
