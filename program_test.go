package xgp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
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
			op.Sub{},
			op.Mul{},
			op.Div{},
		}
		newOp = func(leaf bool, rng *rand.Rand) op.Operator {
			if leaf {
				if rng.Float64() < 0.5 {
					return op.Constant{randFloat64(-5, 5, rng)}
				}
				return op.Variable{randInt(0, 5, rng)}
			}
			return funcs[rng.Intn(len(funcs))]
		}
	)
	return init.Apply(3, 5, newOp, rng)
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
					LossMetric: metrics.MeanAbsoluteError{},
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
