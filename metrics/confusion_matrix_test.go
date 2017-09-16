package metrics

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConfusionMatrixClasses(t *testing.T) {
	var testCases = []struct {
		cm      ConfusionMatrix
		classes []float64
	}{
		{
			cm: ConfusionMatrix{
				0: map[float64]float64{0: 1, 1: 1, 2: 1},
				1: map[float64]float64{0: 1, 1: 1, 2: 1},
				2: map[float64]float64{0: 1, 1: 1, 2: 1},
			},
			classes: []float64{0, 1, 2},
		},
		{
			cm: ConfusionMatrix{
				1: map[float64]float64{0: 1, 1: 1, 2: 1},
				0: map[float64]float64{0: 1, 1: 1, 2: 1},
				2: map[float64]float64{0: 1, 1: 1, 2: 1},
			},
			classes: []float64{0, 1, 2},
		},
		{
			cm: ConfusionMatrix{
				2: map[float64]float64{0: 1, 1: 1, 2: 1},
				0: map[float64]float64{0: 1, 1: 1, 2: 1},
				1: map[float64]float64{0: 1, 1: 1, 2: 1},
			},
			classes: []float64{0, 1, 2},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var classes = tc.cm.Classes()
			if !reflect.DeepEqual(classes, tc.classes) {
				t.Errorf("Error in test case number %d", i)
			}
		})
	}
}

func TestMakeConfusionMatrix(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		cm      ConfusionMatrix
		err     error
	}{
		{
			yTrue:   []float64{0, 1, 0, 1},
			yPred:   []float64{0, 1, 1, 0},
			weights: nil,
			cm: ConfusionMatrix{
				0: map[float64]float64{0: 1, 1: 1},
				1: map[float64]float64{0: 1, 1: 1},
			},
			err: nil,
		},
		{
			yTrue:   []float64{0, 1, 0, 1},
			yPred:   []float64{0, 1, 1},
			weights: nil,
			cm:      nil,
			err:     &errMismatchedLengths{4, 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var cm, err = MakeConfusionMatrix(tc.yTrue, tc.yPred, tc.weights)
			if !reflect.DeepEqual(cm, tc.cm) || !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Error in test case number %d", i)
			}
		})
	}
}

func TestConfusionMatrixTruePositives(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		TPs     []float64
	}{
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{0, 1, 1},
			weights: nil,
			TPs:     []float64{1, 2},
		},
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{1, 0, 1},
			weights: nil,
			TPs:     []float64{0, 1},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 1, 2, 2, 2},
			weights: nil,
			TPs:     []float64{1, 2, 3},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 2, 1, 1, 2},
			weights: nil,
			TPs:     []float64{1, 1, 1},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred, tc.weights)
			for _, class := range cm.Classes() {
				var TP, _ = cm.TruePositives(class)
				if TP != tc.TPs[int(class)] {
					t.Errorf("Error in test case number %d", i)
				}
			}
		})
	}
}

func TestConfusionMatrixFalsePositives(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		FPs     []float64
	}{
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{0, 1, 1},
			weights: nil,
			FPs:     []float64{0, 0},
		},
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{1, 0, 1},
			weights: nil,
			FPs:     []float64{1, 1},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 1, 2, 2, 2},
			weights: nil,
			FPs:     []float64{0, 0, 0},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 2, 1, 1, 2},
			weights: nil,
			FPs:     []float64{0, 2, 1},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred, tc.weights)
			for _, class := range cm.Classes() {
				var FP, _ = cm.FalsePositives(class)
				if FP != tc.FPs[int(class)] {
					t.Errorf("Error in test case number %d", i)
				}
			}
		})
	}
}

func TestConfusionMatrixFalseNegatives(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		FNs     []float64
	}{
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{0, 1, 1},
			weights: nil,
			FNs:     []float64{0, 0},
		},
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{1, 0, 1},
			weights: nil,
			FNs:     []float64{1, 1},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 1, 2, 2, 2},
			weights: nil,
			FNs:     []float64{0, 0, 0},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 2, 1, 1, 2},
			weights: nil,
			FNs:     []float64{0, 1, 2},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred, tc.weights)
			for _, class := range cm.Classes() {
				var FN, _ = cm.FalseNegatives(class)
				if FN != tc.FNs[int(class)] {
					t.Errorf("Error in test case number %d", i)
				}
			}
		})
	}
}

func TestConfusionMatrixTrueNegatives(t *testing.T) {
	var testCases = []struct {
		yTrue   []float64
		yPred   []float64
		weights []float64
		TNs     []float64
	}{
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{0, 1, 1},
			weights: nil,
			TNs:     []float64{2, 1},
		},
		{
			yTrue:   []float64{0, 1, 1},
			yPred:   []float64{1, 0, 1},
			weights: nil,
			TNs:     []float64{1, 0},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 1, 2, 2, 2},
			weights: nil,
			TNs:     []float64{5, 4, 3},
		},
		{
			yTrue:   []float64{0, 1, 1, 2, 2, 2},
			yPred:   []float64{0, 1, 2, 1, 1, 2},
			weights: nil,
			TNs:     []float64{5, 2, 2},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred, tc.weights)
			for _, class := range cm.Classes() {
				var TN, _ = cm.TrueNegatives(class)
				if TN != tc.TNs[int(class)] {
					t.Errorf("Error in test case number %d", i)
				}
			}
		})
	}
}
