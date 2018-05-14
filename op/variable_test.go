package op

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVariableEval(t *testing.T) {
	var testCases = []struct {
		v Variable
		X [][]float64
		Y []float64
	}{
		{
			v: Variable{0},
			X: [][]float64{[]float64{1, 2, 3}, []float64{4, 5, 6}},
			Y: []float64{1, 2, 3},
		},
		{
			v: Variable{1},
			X: [][]float64{[]float64{1, 2, 3}, []float64{4, 5, 6}},
			Y: []float64{4, 5, 6},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var Y = tc.v.Eval(tc.X)
			if !reflect.DeepEqual(Y, tc.Y) {
				t.Errorf("Expected %v, got %v", tc.Y, Y)
			}
		})
	}
}

func TestVariableArity(t *testing.T) {
	var v = Variable{42}
	if v.Arity() != 0 {
		t.Errorf("Expected %d, got %d", 0, v.Arity())
	}
}

func TestVariableString(t *testing.T) {
	var testCases = []struct {
		v Variable
		s string
	}{
		{
			v: Variable{42},
			s: "X[42]",
		},
		{
			v: Variable{},
			s: "X[0]",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var s = tc.v.String()
			if s != tc.s {
				t.Errorf("Expected %v, got %v", tc.s, s)
			}
		})
	}
}
