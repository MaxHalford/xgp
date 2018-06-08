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
)

// A Config contains all the information needed to instantiate an Estimator.
type Config struct {
	// Learning parameters
	LossMetric     metrics.Metric
	EvalMetric     metrics.Metric
	ParsimonyCoeff float64
	// Function parameters
	Funcs     string
	ConstMin  float64
	ConstMax  float64
	PConst    float64
	PFull     float64
	PLeaf     float64
	MinHeight uint
	MaxHeight uint
	// Genetic algorithm parameters
	NPopulations       uint
	NIndividuals       uint
	NGenerations       uint
	NPolishGenerations uint
	PHoistMutation     float64
	PSubtreeMutation   float64
	PPointMutation     float64
	PointMutationRate  float64
	PSubtreeCrossover  float64
	// Other
	RNG *rand.Rand
}

// String representation of a Config. It returns a string containing the
// parameters line by line.
func (c Config) String() string {
	var (
		buffer     = new(bytes.Buffer)
		parameters = [][]string{
			[]string{"Loss metric", c.LossMetric.String()},
			[]string{"Evaluation metric", c.EvalMetric.String()},
			[]string{"Parsimony coefficient", strconv.FormatFloat(c.ParsimonyCoeff, 'g', -1, 64)},

			[]string{"Functions", c.Funcs},
			[]string{"Constant minimum", strconv.FormatFloat(c.ConstMin, 'g', -1, 64)},
			[]string{"Constant maximum", strconv.FormatFloat(c.ConstMax, 'g', -1, 64)},
			[]string{"Constant probability", strconv.FormatFloat(c.PConst, 'g', -1, 64)},
			[]string{"Full initialization probability", strconv.FormatFloat(c.PFull, 'g', -1, 64)},
			[]string{"Terminal probability", strconv.FormatFloat(c.PLeaf, 'g', -1, 64)},
			[]string{"Minimum height", strconv.Itoa(int(c.MinHeight))},
			[]string{"Maximum height", strconv.Itoa(int(c.MaxHeight))},

			[]string{"Number of populations", strconv.Itoa(int(c.NPopulations))},
			[]string{"Number of individuals per population", strconv.Itoa(int(c.NIndividuals))},
			[]string{"Number of generations", strconv.Itoa(int(c.NGenerations))},
			[]string{"Number of tuning generations", strconv.Itoa(int(c.NPolishGenerations))},
			[]string{"Hoist mutation probability", strconv.FormatFloat(c.PHoistMutation, 'g', -1, 64)},
			[]string{"Subtree mutation probability", strconv.FormatFloat(c.PSubtreeMutation, 'g', -1, 64)},
			[]string{"Point mutation probability", strconv.FormatFloat(c.PPointMutation, 'g', -1, 64)},
			[]string{"Point mutation rate", strconv.FormatFloat(c.PointMutationRate, 'g', -1, 64)},
			[]string{"Subtree crossover probability", strconv.FormatFloat(c.PSubtreeCrossover, 'g', -1, 64)},
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
	functions, err := op.ParseFuncs(c.Funcs, ",")
	if err != nil {
		return nil, err
	}

	// Instantiate an Estimator
	var estimator = &Estimator{
		Config:     c,
		Functions:  functions,
		EvalMetric: c.EvalMetric,
		LossMetric: c.LossMetric,
		Initializer: RampedHaldAndHalfInit{
			PFull:    c.PFull,
			FullInit: FullInit{},
			GrowInit: GrowInit{
				PLeaf: c.PLeaf,
			},
		},
	}

	// Set the initial GA
	estimator.GA = &gago.GA{
		NewGenome: func(rng *rand.Rand) gago.Genome {
			var prog = estimator.newProgram(rng)
			return &prog
		},
		NPops:   int(c.NPopulations),
		PopSize: int(c.NIndividuals),
		Model: gaModel{
			selector: gago.SelTournament{
				NContestants: 3,
			},
			pMutate:    c.PHoistMutation + c.PPointMutation + c.PSubtreeMutation,
			pCrossover: c.PSubtreeCrossover,
		},
		RNG:          c.RNG,
		ParallelEval: true,
	}

	// Build fm which maps arities to functions
	estimator.fm = make(map[uint][]op.Operator)
	for _, f := range estimator.Functions {
		var arity = f.Arity()
		if _, ok := estimator.fm[arity]; ok {
			estimator.fm[arity] = append(estimator.fm[arity], f)
		} else {
			estimator.fm[arity] = []op.Operator{f}
		}
	}

	// Set subtree crossover
	estimator.SubtreeCrossover = SubtreeCrossover{
		Weight: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
			if operator.Arity() == 0 {
				return 0.1 // MAGIC
			}
			return 0.9 // MAGIC
		},
	}

	// Set point mutation
	estimator.PointMutation = PointMutation{
		Rate: c.PointMutationRate,
		Mutate: func(operator op.Operator, rng *rand.Rand) op.Operator {
			return estimator.mutateOperator(operator, rng)
		},
	}

	// Set hoist mutation
	estimator.HoistMutation = HoistMutation{
		Weight1: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
			if operator.Arity() == 0 {
				return 0.1 // MAGIC
			}
			return 0.9 // MAGIC
		},
		Weight2: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
			return 1 // MAGIC
		},
	}

	// Set subtree mutation
	estimator.SubtreeMutation = SubtreeMutation{
		Weight: func(operator op.Operator, depth uint, rng *rand.Rand) float64 {
			if operator.Arity() == 0 {
				return 0.1 // MAGIC
			}
			return 0.9 // MAGIC
		},
		NewOperator: func(rng *rand.Rand) op.Operator {
			return estimator.newOperator(rng)
		},
	}

	return estimator, nil
}

// NewConfigWithDefaults returns a Config with default values.
func NewConfigWithDefaults() Config {
	return Config{
		LossMetric: metrics.MeanSquaredError{},
		EvalMetric: metrics.MeanSquaredError{},

		Funcs:     "add,sub,mul,div",
		ConstMin:  -5,
		ConstMax:  5,
		MinHeight: 3,
		MaxHeight: 5,
		PConst:    0.5,
		PFull:     0.5,
		PLeaf:     0.3,

		NPopulations:       1,
		NIndividuals:       100,
		NGenerations:       30,
		NPolishGenerations: 0,
		PHoistMutation:     0.1,
		PPointMutation:     0.1,
		PSubtreeMutation:   0.1,
		PointMutationRate:  0.3,
		PSubtreeCrossover:  0.5,

		ParsimonyCoeff: 0,
	}
}
