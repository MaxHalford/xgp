package metric

import (
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
		var classes = tc.cm.Classes()
		if !reflect.DeepEqual(classes, tc.classes) {
			t.Errorf("Error in test case number %d", i)
		}
	}
}

func TestMakeConfusionMatrix(t *testing.T) {
	var testCases = []struct {
		yTrue []float64
		yPred []float64
		cm    ConfusionMatrix
		err   error
	}{
		{
			yTrue: []float64{0, 1, 0, 1},
			yPred: []float64{0, 1, 1, 0},
			cm: ConfusionMatrix{
				0: map[float64]float64{0: 1, 1: 1},
				1: map[float64]float64{0: 1, 1: 1},
			},
			err: nil,
		},
		{
			yTrue: []float64{0, 1, 0, 1},
			yPred: []float64{0, 1, 1},
			cm:    nil,
			err:   &errMismatchedLengths{4, 3},
		},
	}
	for i, tc := range testCases {
		var cm, err = MakeConfusionMatrix(tc.yTrue, tc.yPred)
		if !reflect.DeepEqual(cm, tc.cm) || !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Error in test case number %d", i)
		}
	}
}

func TestConfusionMatrixTruePositives(t *testing.T) {
	var testCases = []struct {
		yTrue []float64
		yPred []float64
		tps   []float64
	}{
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{0, 1, 1},
			tps:   []float64{1, 2},
		},
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{1, 0, 1},
			tps:   []float64{0, 1},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 1, 2, 2, 2},
			tps:   []float64{1, 2, 3},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 2, 1, 1, 2},
			tps:   []float64{1, 1, 1},
		},
	}
	for i, tc := range testCases {
		var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred)
		for _, class := range cm.Classes() {
			var TP, _ = cm.TruePositives(class)
			if TP != tc.tps[int(class)] {
				t.Errorf("Error in test case number %d", i)
			}
		}
	}
}

func TestConfusionMatrixFalsePositives(t *testing.T) {
	var testCases = []struct {
		yTrue []float64
		yPred []float64
		fps   []float64
	}{
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{0, 1, 1},
			fps:   []float64{0, 0},
		},
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{1, 0, 1},
			fps:   []float64{1, 1},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 1, 2, 2, 2},
			fps:   []float64{0, 0, 0},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 2, 1, 1, 2},
			fps:   []float64{0, 2, 1},
		},
	}
	for i, tc := range testCases {
		var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred)
		for _, class := range cm.Classes() {
			var FP, _ = cm.FalsePositives(class)
			if FP != tc.fps[int(class)] {
				t.Errorf("Error in test case number %d", i)
			}
		}
	}
}

func TestConfusionMatrixFalseNegatives(t *testing.T) {
	var testCases = []struct {
		yTrue []float64
		yPred []float64
		fns   []float64
	}{
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{0, 1, 1},
			fns:   []float64{0, 0},
		},
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{1, 0, 1},
			fns:   []float64{1, 1},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 1, 2, 2, 2},
			fns:   []float64{0, 0, 0},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 2, 1, 1, 2},
			fns:   []float64{0, 1, 2},
		},
	}
	for i, tc := range testCases {
		var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred)
		for _, class := range cm.Classes() {
			var FN, _ = cm.FalseNegatives(class)
			if FN != tc.fns[int(class)] {
				t.Errorf("Error in test case number %d", i)
			}
		}
	}
}

func TestConfusionMatrixTrueNegatives(t *testing.T) {
	var testCases = []struct {
		yTrue []float64
		yPred []float64
		tns   []float64
	}{
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{0, 1, 1},
			tns:   []float64{2, 1},
		},
		{
			yTrue: []float64{0, 1, 1},
			yPred: []float64{1, 0, 1},
			tns:   []float64{1, 0},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 1, 2, 2, 2},
			tns:   []float64{5, 4, 3},
		},
		{
			yTrue: []float64{0, 1, 1, 2, 2, 2},
			yPred: []float64{0, 1, 2, 1, 1, 2},
			tns:   []float64{5, 2, 2},
		},
	}
	for i, tc := range testCases {
		var cm, _ = MakeConfusionMatrix(tc.yTrue, tc.yPred)
		for _, class := range cm.Classes() {
			var TN, _ = cm.TrueNegatives(class)
			if TN != tc.tns[int(class)] {
				t.Errorf("Error in test case number %d", i)
			}
		}
	}
}
