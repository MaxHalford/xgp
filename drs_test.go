package koza

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClosestFloat64(t *testing.T) {
	var testCases = []struct {
		float64s []float64
		f        float64
		closest  float64
	}{
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -11,
			closest:  -10,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -10,
			closest:  -10,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -9,
			closest:  -10,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -6,
			closest:  -5,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -5,
			closest:  -5,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -4,
			closest:  -5,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        -1,
			closest:  0,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        0,
			closest:  0,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        1,
			closest:  0,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        4,
			closest:  5,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        5,
			closest:  5,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        6,
			closest:  5,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        9,
			closest:  10,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        10,
			closest:  10,
		},
		{
			float64s: []float64{-10, -5, 0, 5, 10},
			f:        11,
			closest:  10,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var closest = closestFloat64(tc.float64s, tc.f)
			if closest != tc.closest {
				t.Errorf("Expected %f, got %f", tc.closest, closest)
			}
		})
	}
}

func TestDynamicRangeSelection(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		yVal    []float64
		classes []float64
		drs     DynamicRangeSelection
	}{
		{
			yTrue:   []float64{0, 0, 1, 2, 2},
			yPred:   []float64{-60, -30, 0, 30, 60},
			yVal:    []float64{-70, -60, -50, -40, -30, -20, -10, 0, 10, 20, 30, 40, 50, 60, 70},
			classes: []float64{0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 2, 2, 2, 2, 2},
			drs:     DynamicRangeSelection{},
		},
		{
			yTrue:   []float64{0, 1, 2, 1, 0},
			yPred:   []float64{-60, -30, 0, 30, 60},
			yVal:    []float64{-70, -60, -50, -40, -30, -20, -10, 0, 10, 20, 30, 40, 50, 60, 70},
			classes: []float64{0, 0, 0, 1, 1, 1, 2, 2, 2, 1, 1, 1, 0, 0, 0},
			drs:     DynamicRangeSelection{},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.drs.Fit(tc.yTrue, tc.yPred)
			var classes = tc.drs.Predict(tc.yVal)
			if !reflect.DeepEqual(classes, tc.classes) {
				t.Errorf("Expected %f, got %f", tc.classes, classes)
			}
		})
	}
}
