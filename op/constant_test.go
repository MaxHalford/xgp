package op

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConstantEval(t *testing.T) {
	var testCases = []struct {
		c Constant
		X [][]float64
		Y []float64
	}{
		{
			c: Constant{42},
			X: [][]float64{[]float64{}},
			Y: []float64{},
		},
		{
			c: Constant{42},
			X: [][]float64{[]float64{1, 2, 3}},
			Y: []float64{42, 42, 42},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var Y = tc.c.Eval(tc.X)
			if !reflect.DeepEqual(Y, tc.Y) {
				t.Errorf("Expected %v, got %v", tc.Y, Y)
			}
		})
	}
}

func TestConstantString(t *testing.T) {
	var testCases = []struct {
		c Constant
		s string
	}{
		{
			c: Constant{42},
			s: "42",
		},
		{
			c: Constant{-42},
			s: "-42",
		},
		{
			c: Constant{42.24},
			s: "42.24",
		},
		{
			c: Constant{},
			s: "0",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var s = tc.c.String()
			if s != tc.s {
				t.Errorf("Expected %v, got %v", tc.s, s)
			}
		})
	}
}
