package op

import (
	"fmt"
	"testing"
)

func TestMaxSimplify(t *testing.T) {
	var testCases = []struct {
		in  Max
		out Operator
	}{
		{
			in:  Max{Mul{Add{Const{1}, Const{2}}, Var{0}}, Const{2}},
			out: Max{Mul{Const{3}, Var{0}}, Const{2}},
		},
		{
			in:  Max{Const{2}, Mul{Add{Const{1}, Const{2}}, Var{0}}},
			out: Max{Const{2}, Mul{Const{3}, Var{0}}},
		},
		{
			in:  Max{Const{2}, Const{3}},
			out: Const{3},
		},
		{
			in:  Max{Var{2}, Const{3}},
			out: Max{Var{2}, Const{3}},
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
