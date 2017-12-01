package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConstantApplyRow(t *testing.T) {
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
			var y = tc.c.ApplyRow(tc.x)
			if y != tc.y {
				t.Errorf("Expected %v, got %v", tc.y, y)
			}
		})
	}
}

func TestConstantApplyCols(t *testing.T) {
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
			var Y = tc.c.ApplyCols(tc.X)
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
			s: "42.000",
		},
		{
			c: Constant{-42},
			s: "-42.000",
		},
		{
			c: Constant{42.24},
			s: "42.240",
		},
		{
			c: Constant{},
			s: "0.000",
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
