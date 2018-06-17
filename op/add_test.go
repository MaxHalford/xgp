package op

import (
	"fmt"
	"testing"
)

func TestAddSimplify(t *testing.T) {
	var testCases = []struct {
		in  Add
		out Operator
	}{
		{
			in:  Add{Mul{Const{1}, Const{-2}}, Var{0}},
			out: Add{Const{-2}, Var{0}},
		},
		{
			in:  Add{Var{0}, Mul{Const{1}, Const{-2}}},
			out: Add{Var{0}, Const{-2}},
		},
		{
			in:  Add{Const{0}, Var{1}},
			out: Var{1},
		},
		{
			in:  Add{Var{1}, Const{0}},
			out: Var{1},
		},
		{
			in:  Add{Const{42}, Const{43}},
			out: Const{85},
		},
		{
			in:  Add{Var{0}, Const{42}},
			out: Add{Var{0}, Const{42}},
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
