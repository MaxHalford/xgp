package xgp

import (
	"testing"

	"github.com/MaxHalford/xgp/dataframe"
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
				Root: &Node{
					Operator: Sum{},
					Children: []*Node{
						&Node{Operator: Variable{0}},
						&Node{Operator: Variable{1}},
					},
				},
			},
			output: 2,
		},
	}

	for i, tc := range testCases {
		var output = tc.program.PredictRow(tc.row)
		if output != tc.output {
			t.Errorf("Error in test case number %d: got %.2f instead of %.2f", i, output, tc.output)
		}
	}
}

func TestPredictDataFrame(t *testing.T) {
	var testCases = []struct {
		dataframe *dataframe.DataFrame
		program   Program
		y         []float64
	}{
		{
			dataframe: &dataframe.DataFrame{
				X: [][]float64{
					[]float64{1, 1},
					[]float64{1, 2},
					[]float64{1, 3},
				},
			},
			program: Program{
				Root: &Node{
					Operator: Sum{},
					Children: []*Node{
						&Node{Operator: Variable{0}},
						&Node{Operator: Variable{1}},
					},
				},
			},
			y: []float64{2, 3, 4},
		},
	}

	for i, tc := range testCases {
		var y = tc.program.PredictDataFrame(tc.dataframe)
		for j := range y {
			if y[j] != tc.y[j] {
				t.Errorf("Error in test case number %d: got %.2f instead of %.2f", i, y[j], tc.y[j])
			}
		}
	}
}
