package xgp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
)

func TestEvaluate(t *testing.T) {
	var testCases = []struct {
		prog    Program
		fitness float64
	}{
		{
			prog: Program{
				Op: op.Add{op.Var{0}, op.Var{1}},
				Estimator: &Estimator{
					XTrain: [][]float64{
						[]float64{1, 2, 3},
						[]float64{2, 1, 4},
					},
					YTrain:     []float64{2, 4, 7},
					LossMetric: metrics.MeanAbsoluteError{},
				},
			},
			fitness: 2.0 / 3,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			fitness, err := tc.prog.Evaluate()
			if err != nil {
				t.Errorf("Expected nil error, got %s", err)
			}
			if fitness != tc.fitness {
				t.Errorf("Expected %f, got %f", tc.fitness, fitness)
			}
		})
	}
}

func TestMutate(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			in  Program
			out Program
		}{
			{
				in: Program{
					Op: op.Add{op.Const{42}, op.Var{1}},
					Estimator: &Estimator{
						Config: Config{PHoistMutation: 1},
						HoistMutation: HoistMutation{
							Weight1: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
								switch operator.(type) {
								case op.Add:
									return 1
								default:
									return 0
								}
							},
							Weight2: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
								switch operator.(type) {
								case op.Const:
									return 1
								default:
									return 0
								}
							},
						},
					},
				},
				out: Program{
					Op: op.Const{42},
				},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.in.Mutate(rng)
			if tc.in.String() != tc.out.String() {
				t.Errorf("Expected %s, got %s", tc.out, tc.in)
			}
		})
	}
}

func TestGAValidate(t *testing.T) {
	var err = gaModel{}.Validate()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}
