package xgp

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/MaxHalford/gago"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
)

// A Config contains all the information needed to instantiate an Estimator.
type Config struct {
	ConstMax float64
	ConstMin float64

	EvalMetric metrics.Metric
	LossMetric metrics.Metric

	Funcs string

	MinHeight int
	MaxHeight int

	NPopulations       int
	NIndividuals       int
	NGenerations       int
	NTuningGenerations int

	PConstant float64
	PFull     float64
	PTerminal float64

	PHoistMutation    float64
	PSubTreeMutation  float64
	PPointMutation    float64
	PointMutationRate float64

	PSubTreeCrossover float64

	ParsimonyCoeff float64

	RNG *rand.Rand
}

// String representation of a Config.
func (c Config) String() string {
	var (
		buffer     = new(bytes.Buffer)
		parameters = [][]string{
			[]string{"Constant minimum", strconv.FormatFloat(c.ConstMin, 'g', -1, 64)},
			[]string{"Constant maximum", strconv.FormatFloat(c.ConstMax, 'g', -1, 64)},

			[]string{"Evaluation metric", c.EvalMetric.String()},
			[]string{"Loss metric", c.LossMetric.String()},

			[]string{"Functions", c.Funcs},

			[]string{"Minimum height", strconv.Itoa(c.MinHeight)},
			[]string{"Maximum height", strconv.Itoa(c.MaxHeight)},

			[]string{"Number of populations", strconv.Itoa(c.NPopulations)},
			[]string{"Number of individuals per population", strconv.Itoa(c.NIndividuals)},
			[]string{"Number of generations", strconv.Itoa(c.NGenerations)},
			[]string{"Number of tuning generations", strconv.Itoa(c.NGenerations)},

			[]string{"Constant probability", strconv.FormatFloat(c.PConstant, 'g', -1, 64)},
			[]string{"Full initialization probability", strconv.FormatFloat(c.PFull, 'g', -1, 64)},
			[]string{"Terminal probability", strconv.FormatFloat(c.PTerminal, 'g', -1, 64)},

			[]string{"Hoist mutation probability", strconv.FormatFloat(c.PHoistMutation, 'g', -1, 64)},
			[]string{"Sub-tree mutation probability", strconv.FormatFloat(c.PSubTreeMutation, 'g', -1, 64)},
			[]string{"Point mutation probability", strconv.FormatFloat(c.PPointMutation, 'g', -1, 64)},
			[]string{"Point mutation rate", strconv.FormatFloat(c.PointMutationRate, 'g', -1, 64)},

			[]string{"Sub-tree crossover probability", strconv.FormatFloat(c.PSubTreeCrossover, 'g', -1, 64)},

			[]string{"Parsimony coefficient", strconv.FormatFloat(c.ParsimonyCoeff, 'g', -1, 64)},
		}
	)
	for _, param := range parameters {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", param[0], param[1]))
	}
	return strings.Trim(buffer.String(), "\n")
}

// NewEstimator returns an Estimator from a Config.
func (c Config) NewEstimator() (*Estimator, error) {

	// Default the evaluation metric to the fitness metric if it's nil
	if c.EvalMetric == nil {
		c.EvalMetric = c.LossMetric
	}

	// The convention is to use a fitness metric which has to be minimized
	if c.LossMetric.BiggerIsBetter() {
		c.LossMetric = metrics.NegativeMetric{Metric: c.LossMetric}
	}

	// Determine the functions to use
	functions, err := op.ParseStringFuncs(c.Funcs)
	if err != nil {
		return nil, err
	}

	// Instantiate an Estimator
	var estimator = &Estimator{
		Config:     c,
		Functions:  functions,
		EvalMetric: c.EvalMetric,
		LossMetric: c.LossMetric,
		Initializer: RampedHaldAndHalfInitializer{
			PFull:           c.PFull,
			FullInitializer: FullInitializer{},
			GrowInitializer: GrowInitializer{
				PTerminal: c.PTerminal,
			},
		},
	}

	// Set the initial GA
	estimator.GA = &gago.GA{
		NewGenome: estimator.newProgram,
		NPops:     c.NPopulations,
		PopSize:   c.NIndividuals,
		Model: gaModel{
			selector: gago.SelTournament{
				NContestants: 3,
			},
			pMutate:    c.PHoistMutation + c.PPointMutation + c.PSubTreeMutation,
			pCrossover: c.PSubTreeCrossover,
		},
		RNG:          c.RNG,
		ParallelEval: true,
	}

	// Build fm which maps arities to functions
	estimator.fm = make(map[int][]op.Operator)
	for _, f := range estimator.Functions {
		var arity = f.Arity()
		if _, ok := estimator.fm[arity]; ok {
			estimator.fm[arity] = append(estimator.fm[arity], f)
		} else {
			estimator.fm[arity] = []op.Operator{f}
		}
	}

	// Set crossover methods

	estimator.SubTreeCrossover = SubTreeCrossover{
		Picker: WeightedPicker{
			Weighting: Weighting{
				PConstant: 0.1, // MAGIC
				PVariable: 0.1, // MAGIC
				PFunction: 0.8, // MAGIC
			},
		},
	}

	// Set mutation methods

	estimator.PointMutation = PointMutation{
		Weighting: Weighting{
			PConstant: c.PointMutationRate,
			PVariable: c.PointMutationRate,
			PFunction: c.PointMutationRate,
		},
		MutateOperator: func(operator op.Operator, rng *rand.Rand) op.Operator {
			return estimator.mutateOperator(operator, rng)
		},
	}

	estimator.HoistMutation = HoistMutation{
		Picker: WeightedPicker{
			Weighting: Weighting{
				PConstant: 0.1, // MAGIC
				PVariable: 0.1, // MAGIC
				PFunction: 0.8, // MAGIC
			},
		},
	}

	estimator.SubTreeMutation = SubTreeMutation{
		Crossover: SubTreeCrossover{
			Picker: WeightedPicker{
				Weighting: Weighting{
					PConstant: 0.1, // MAGIC
					PVariable: 0.1, // MAGIC
					PFunction: 0.8, // MAGIC
				},
			},
		},
		NewTree: func(rng *rand.Rand) tree.Tree {
			return estimator.newTree(rng)
		},
	}

	return estimator, nil
}

// NewConfigWithDefaults returns a Config with default values.
func NewConfigWithDefaults() Config {
	return Config{
		ConstMin: -5,
		ConstMax: 5,

		EvalMetric: metrics.MeanSquaredError{},
		LossMetric: metrics.MeanSquaredError{},

		Funcs: "sum,sub,mul,div",

		MinHeight: 3,
		MaxHeight: 5,

		NPopulations:       1,
		NIndividuals:       50,
		NGenerations:       30,
		NTuningGenerations: 0,

		PConstant: 0.5,
		PFull:     0.5,
		PTerminal: 0.3,

		PHoistMutation:    0.1,
		PPointMutation:    0.1,
		PSubTreeMutation:  0.1,
		PointMutationRate: 0.3,

		PSubTreeCrossover: 0.5,

		ParsimonyCoeff: 0,
	}
}
