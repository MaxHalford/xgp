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
			in:  Mul{Const{0}, Var{42}},
			out: Const{0},
		},
		{
			in:  Mul{Var{42}, Const{0}},
			out: Const{0},
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
