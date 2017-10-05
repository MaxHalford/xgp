package xgp

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConstantApply(t *testing.T) {
	var testCases = []struct {
		c Constant
		x []float64
		y float64
	}{
		{
			c: Constant{42},
			x: []float64{},
			y: 42,
		},
		{
			c: Constant{42},
			x: []float64{1, 2, 3},
			y: 42,
		},
		{
			c: Constant{42},
			x: nil,
			y: 42,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var y = tc.c.Apply(tc.x)
			if y != tc.y {
				t.Errorf("Expected %v, got %v", tc.y, y)
			}
		})
	}
}

func TestConstantApplyXT(t *testing.T) {
	var testCases = []struct {
		c  Constant
		XT [][]float64
		Y  []float64
	}{
		{
			c:  Constant{42},
			XT: [][]float64{[]float64{}},
			Y:  []float64{},
		},
		{
			c:  Constant{42},
			XT: [][]float64{[]float64{1, 2, 3}},
			Y:  []float64{42, 42, 42},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var Y = tc.c.ApplyXT(tc.XT)
			if !reflect.DeepEqual(Y, tc.Y) {
				t.Errorf("Expected %v, got %v", tc.Y, Y)
			}
		})
	}
}

func TestConstantArity(t *testing.T) {
	var c = Constant{42}
	if c.Arity() != 0 {
		t.Errorf("Expected %d, got %d", 0, c.Arity())
	}
}

func TestConstantString(t *testing.T) {
	var testCases = []struct {
		c Constant
		s string
	}{
		{
			c: Constant{42},
			s: "42.00",
		},
		{
			c: Constant{-42},
			s: "-42.00",
		},
		{
			c: Constant{42.24},
			s: "42.24",
		},
		{
			c: Constant{},
			s: "0.00",
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
