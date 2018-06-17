package op

import (
	"fmt"
	"testing"
)

func TestDivSimplify(t *testing.T) {
	var testCases = []struct {
		in  Div
		out Operator
	}{
		{
			in:  Div{Mul{Const{1}, Const{-2}}, Var{0}},
			out: Div{Const{-2}, Var{0}},
		},
		{
			in:  Div{Var{0}, Mul{Const{1}, Const{-2}}},
			out: Div{Var{0}, Const{-2}},
		},
		{
			in:  Div{Var{1}, Const{0}},
			out: Const{1},
		},
		{
			in:  Div{Var{2}, Const{1}},
			out: Var{2},
		},
		{
			in:  Div{Var{2}, Const{-1}},
			out: Neg{Var{2}},
		},
		{
			in:  Div{Const{27}, Const{9}},
			out: Const{3},
		},
		{
			in:  Div{Var{1}, Var{1}},
			out: Const{1},
		},
		{
			in:  Div{Const{0}, Var{1}},
			out: Const{0},
		},
		{
			in:  Div{Const{1}, Var{1}},
			out: Inv{Var{1}},
		},
		{
			in:  Div{Const{-1}, Var{1}},
			out: Inv{Neg{Var{1}}},
		},
		{
			in:  Div{Var{0}, Var{1}},
			out: Div{Var{0}, Var{1}},
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
