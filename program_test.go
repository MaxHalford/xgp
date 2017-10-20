package xgp

import (
	"fmt"
	"testing"

	"github.com/MaxHalford/xgp/tree"
)

func TestPredictRow(t *testing.T) {
	var testCases = []struct {
		row     []float64
		program Program
		output  float64
	}{
		{
			row: []float64{1, 1},
			program: Program{
				Tree: &tree.Tree{
					Operator: tree.Sum{},
					Branches: []*tree.Tree{
						&tree.Tree{Operator: tree.Variable{0}},
						&tree.Tree{Operator: tree.Variable{1}},
					},
				},
			},
			output: 2,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var output, _ = tc.program.PredictRow(tc.row)
			if output != tc.output {
				t.Errorf("Expected %.2f, got %.2f", tc.output, output)
			}
		})
	}
}

func TestPredict(t *testing.T) {
	var testCases = []struct {
		X       [][]float64
		program Program
		y       []float64
	}{
		{
			X: [][]float64{
				[]float64{1, 1},
				[]float64{1, 2},
				[]float64{1, 3},
			},
			program: Program{
				Tree: &tree.Tree{
					Operator: tree.Sum{},
					Branches: []*tree.Tree{
						&tree.Tree{Operator: tree.Variable{0}},
						&tree.Tree{Operator: tree.Variable{1}},
					},
				},
			},
			y: []float64{2, 3, 4},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var y, _ = tc.program.Predict(tc.X)
			for j := range y {
				if y[j] != tc.y[j] {
					t.Errorf("Expected %.2f, got %.2f", tc.y[j], y[j])
				}
			}
		})
	}
}
