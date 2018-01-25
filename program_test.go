package koza

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/koza/metrics"
	"github.com/MaxHalford/koza/op"
	"github.com/MaxHalford/koza/tree"
)

// randTree is a convenience method that produces a random Tree for testing
// purposes.
func randTree(rng *rand.Rand) tree.Tree {
	var (
		init = RampedHaldAndHalfInitializer{
			PFull:           0.5,
			FullInitializer: FullInitializer{},
			GrowInitializer: GrowInitializer{0.5},
		}
		funcs = []op.Operator{
			op.Cos{},
			op.Sin{},
			op.Sum{},
			op.Difference{},
			op.Product{},
			op.Division{},
		}
		of = OperatorFactory{
			PConstant:   0.5,
			NewConstant: func(rng *rand.Rand) op.Constant { return op.Constant{randFloat64(-5, 5, rng)} },
			NewVariable: func(rng *rand.Rand) op.Variable { return op.Variable{randInt(0, 5, rng)} },
			NewFunction: func(rng *rand.Rand) op.Operator { return funcs[rng.Intn(len(funcs))] },
		}
	)
	return init.Apply(3, 5, of, rng)
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
				Tree: tree.MustParseCode("sum(X[0], X[1])"),
				Task: Task{
					Metric: metrics.MeanAbsoluteError{},
				},
			},
			y: []float64{2, 3, 4},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var y, _ = tc.program.Predict(tc.X, false)
			for j := range y {
				if y[j] != tc.y[j] {
					t.Errorf("Expected %.2f, got %.2f", tc.y[j], y[j])
				}
			}
		})
	}
}
