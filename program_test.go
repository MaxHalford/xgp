package xgp

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
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

func TestProgramString(t *testing.T) {
	var (
		tree = randTree(newRand())
		prog = Program{Tree: tree}
	)
	if prog.String() != tree.String() {
		t.Error("The program's string representation should be the same as it's tree")
	}
}

func TestProgramPredict(t *testing.T) {
	var testCases = []struct {
		X       [][]float64
		program Program
		proba   bool
		y       []float64
	}{
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{Tree: tree.MustParseCode("sum(X[0], X[1])")},
			proba:   false,
			y:       []float64{-0.2, 0.1, 0.6, 3},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{Tree: tree.MustParseCode("sum(X[0], X[1])")},
			proba:   true,
			y:       []float64{-0.2, 0.1, 0.6, 3},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{},
			},
			proba: false,
			y:     []float64{-0.2, 0.1, 0.6, 3},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{},
			},
			proba: true,
			y:     []float64{-0.2, 0.1, 0.6, 3},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.MeanSquaredError{}},
			},
			proba: false,
			y:     []float64{-0.2, 0.1, 0.6, 3},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.MeanSquaredError{}},
			},
			proba: true,
			y:     []float64{-0.2, 0.1, 0.6, 3},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.Accuracy{}},
			},
			proba: false,
			y:     []float64{0, 0, 1, 1},
		},
		{
			X: [][]float64{
				[]float64{0.1, -0.3, 0.4, 1},
				[]float64{-0.3, 0.4, 0.2, 2},
			},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.Accuracy{}},
			},
			proba: true,
			y:     []float64{0.45017, 0.52498, 0.64566, 0.95257},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var y, _ = tc.program.Predict(tc.X, tc.proba)
			for j := range y {
				if math.Abs(y[j]-tc.y[j]) > 10e-5 {
					t.Errorf("Expected %.5f, got %.5f", tc.y[j], y[j])
				}
			}
		})
	}
}

func TestProgramPredictPartial(t *testing.T) {
	var testCases = []struct {
		x       []float64
		program Program
		proba   bool
		y       float64
	}{
		{
			x:       []float64{0.1, -0.3},
			program: Program{Tree: tree.MustParseCode("sum(X[0], X[1])")},
			proba:   false,
			y:       -0.2,
		},
		{
			x:       []float64{0.1, -0.3},
			program: Program{Tree: tree.MustParseCode("sum(X[0], X[1])")},
			proba:   true,
			y:       -0.2,
		},
		{
			x: []float64{0.1, -0.3},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{},
			},
			proba: false,
			y:     -0.2,
		},
		{
			x: []float64{0.1, -0.3},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{},
			},
			proba: true,
			y:     -0.2,
		},
		{
			x: []float64{0.1, -0.3},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.MeanSquaredError{}},
			},
			proba: false,
			y:     -0.2,
		},
		{
			x: []float64{0.1, -0.3},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.MeanSquaredError{}},
			},
			proba: true,
			y:     -0.2,
		},
		{
			x: []float64{0.1, -0.3},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.Accuracy{}},
			},
			proba: false,
			y:     0,
		},
		{
			x: []float64{0.1, -0.3},
			program: Program{
				Tree:      tree.MustParseCode("sum(X[0], X[1])"),
				Estimator: &Estimator{LossMetric: metrics.Accuracy{}},
			},
			proba: true,
			y:     0.45017,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var y, _ = tc.program.PredictPartial(tc.x, tc.proba)
			if math.Abs(y-tc.y) > 10e-5 {
				t.Errorf("Expected %.5f, got %.5f", tc.y, y)
			}
		})
	}
}

func TestProgramMarshalJSON(t *testing.T) {
	var prog = &Program{
		Tree:      randTree(newRand()),
		Estimator: &Estimator{LossMetric: metrics.BinaryLogLoss{}},
	}
	var bytes, err = prog.MarshalJSON()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var prog2 = &Program{}
	prog2.UnmarshalJSON(bytes)
	if !reflect.DeepEqual(prog2, prog) {
		t.Errorf("Expected %v, got %v", prog, prog2)
	}
}
