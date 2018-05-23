package op

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseFunc(t *testing.T) {
	var testCases = []struct {
		name string
		op   Operator
	}{
		{"cos", Cos{}},
		{"sin", Sin{}},
		{"exp", Exp{}},
		{"max", Max{}},
		{"min", Min{}},
		{"sum", Sum{}},
		{"sub", Sub{}},
		{"div", Div{}},
		{"mul", Mul{}},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var op, _ = ParseFunc(tc.name)
			if !reflect.DeepEqual(op, tc.op) {
				t.Errorf("Expected %v, got %v", tc.op, op)
			}
		})
	}
}

func TestParseFuncs(t *testing.T) {
	var testCases = []struct {
		names       string
		sep         string
		ops         []Operator
		raisesError bool
	}{
		{"cos,sin,exp", ",", []Operator{Cos{}, Sin{}, Exp{}}, false},
		{"cos,sin,,exp", ",", nil, true},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var ops, err = ParseFuncs(tc.names, ",")
			if !reflect.DeepEqual(ops, tc.ops) {
				t.Errorf("Expected %v, got %v", tc.ops, ops)
			}
			if (err != nil) != tc.raisesError {
				t.Errorf("Expected an error to be raised")
			}
		})
	}
}
