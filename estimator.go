package xgp

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/MaxHalford/gago"
	"github.com/gosuri/uiprogress"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
)

// An Estimator links all the different components together and can be used to
// train Programs on a dataset. You shouldn't instantiate this struct directly;
// instead you should use the Config struct.
type Estimator struct {
	Config

	EvalMetric       metrics.Metric
	LossMetric       metrics.Metric
	Functions        []op.Operator
	Initializer      Initializer
	GA               *gago.GA
	PolishGA         *gago.GA
	PointMutation    PointMutation
	SubtreeMutation  SubtreeMutation
	HoistMutation    HoistMutation
	SubtreeCrossover SubtreeCrossover

	fm       map[int][]op.Operator
	XTrain   [][]float64
	YTrain   []float64
	WTrain   []float64
	XVal     [][]float64
	YVal     []float64
	WVal     []float64
	nClasses int
}

// String representation of an Estimator.
func (est Estimator) String() string {
	return est.Config.String()
}

// BestProgram returns an Estimator's bestProgram in a safe way.
func (est Estimator) BestProgram() *Program {
	return est.GA.HallOfFame[0].Genome.(*Program)
}

// Fit an Estimator to a dataset.Dataset.
func (est *Estimator) Fit(
	XTrain [][]float64,
	YTrain []float64,
	WTrain []float64,
	XVal [][]float64,
	YVal []float64,
	WVal []float64,
	verbose bool,
) (prog Program, err error) {

	// Set the training set
	est.XTrain = XTrain
	est.YTrain = YTrain
	est.WTrain = WTrain

	// Set the validation set
	est.XVal = XVal
	est.YVal = YVal
	est.WVal = WVal

	// Count the number of classes if the task is classification
	if est.LossMetric.Classification() {
		// Check that the task to perform is not multi-class classification
		est.nClasses = countDistinct(YTrain)
		if est.nClasses > 2 {
			return prog, errors.New("Multi-class classification is not supported")
		}
	}

	// Initialize the GA
	err = est.GA.Initialize()
	if err != nil {
		return prog, err
	}

	// Evolve the GA
	var (
		bar      *uiprogress.Bar
		progress *uiprogress.Progress
	)
	if verbose {
		var start = time.Now()
		progress = uiprogress.New()
		progress.Start()
		bar = progress.AddBar(est.NGenerations)
		bar.PrependCompleted()
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			// Add time spent
			var message = fmtDuration(time.Since(start))
			// Add training error
			var (
				best            = est.BestProgram()
				yTrainPred, err = best.Predict(est.XTrain, est.EvalMetric.NeedsProbabilities())
			)
			if err != nil {
				return ""
			}
			trainScore, err := est.EvalMetric.Apply(est.YTrain, yTrainPred, nil)
			if err != nil {
				return ""
			}
			message += fmt.Sprintf(", train %s: %.5f", est.EvalMetric.String(), trainScore)
			// Add validation error
			if est.XVal != nil && est.YVal != nil {
				yEvalPred, err := best.Predict(est.XVal, est.EvalMetric.NeedsProbabilities())
				if err != nil {
					return ""
				}
				evalScore, err := est.EvalMetric.Apply(est.YVal, yEvalPred, est.WVal)
				if err != nil {
					return ""
				}
				message += fmt.Sprintf(", val %s: %.5f", est.EvalMetric.String(), evalScore)
			}
			return message
		})
	}

	for i := 0; i < est.NGenerations; i++ {
		if verbose {
			bar.Incr()
		}

		// Make sure each tree has at least a height of 2
		for i, pop := range est.GA.Populations {
			for j, indi := range pop.Individuals {
				var prog = indi.Genome.(*Program)
				if prog.Tree.Height() < 2 { // MAGIC
					est.SubtreeMutation.Apply(&prog.Tree, pop.RNG)
					est.GA.Populations[i].Individuals[j].Evaluate()
				}
			}
		}

		err = est.GA.Evolve()
		if err != nil {
			return prog, err
		}
	}

	// Close the progress bar
	if verbose {
		progress.Stop()
	}

	// Return the best program
	var best = est.BestProgram()
	return *best, nil
}

func (est Estimator) newConstant(rng *rand.Rand) op.Constant {
	return op.Constant{
		Value: est.Config.ConstMin + rng.Float64()*(est.Config.ConstMax-est.ConstMin),
	}
}

func (est Estimator) newVariable(rng *rand.Rand) op.Variable {
	return op.Variable{Index: rng.Intn(len(est.XTrain))}
}

func (est Estimator) newFunction(rng *rand.Rand) op.Operator {
	return est.Functions[rng.Intn(len(est.Functions))]
}

func (est Estimator) newFunctionOfArity(arity int, rng *rand.Rand) op.Operator {
	return est.fm[arity][rng.Intn(len(est.fm[arity]))]
}

func (est Estimator) mutateOperator(operator op.Operator, rng *rand.Rand) op.Operator {
	switch operator.(type) {
	case op.Constant:
		return op.Constant{Value: operator.(op.Constant).Value * rng.NormFloat64()}
	case op.Variable:
		return est.newVariable(rng)
	default:
		return est.newFunctionOfArity(operator.Arity(), rng)
	}
}

func (est Estimator) newTree(rng *rand.Rand) tree.Tree {
	return est.Initializer.Apply(
		est.MinHeight,
		est.MaxHeight,
		func(leaf bool, rng *rand.Rand) op.Operator {
			if leaf {
				if rng.Float64() < est.PConstant {
					return est.newConstant(rng)
				}
				return est.newVariable(rng)
			}
			return est.newFunction(rng)
		},
		rng,
	)
}

// newProgram can be used by gago to produce a new Genome.
func (est *Estimator) newProgram(rng *rand.Rand) gago.Genome {
	var prog = Program{
		Tree:      est.newTree(rng),
		Estimator: est,
	}
	return &prog
}

// newProgramPolish can be used by gago to produce a new Genome.
func (est *Estimator) newProgramPolish(rng *rand.Rand) gago.Genome {
	var (
		bestProg   = est.GA.HallOfFame[0].Genome.(*Program)
		progPolish = newProgramPolish(*bestProg)
	)
	progPolish.jitterConstants(rng)
	return &progPolish
}
