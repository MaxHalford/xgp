package koza

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/gosuri/uiprogress"

	"github.com/MaxHalford/koza/ensemble"
	"github.com/MaxHalford/koza/metrics"
	"github.com/MaxHalford/koza/op"
	"github.com/MaxHalford/koza/tree"
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
	TuningGA         *gago.GA
	PointMutation    PointMutation
	SubTreeMutation  SubTreeMutation
	HoistMutation    HoistMutation
	SubTreeCrossover SubTreeCrossover

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
) (ensemble.Predictor, error) {

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
		est.nClasses = countDistinct(YTrain)
		// Check that the task to perform is not multi-class classification
		if est.nClasses > 2 {
			return nil, errors.New("Multi-class classification is not supported")
		}
	}

	// Initialize the GA
	est.GA.Initialize()

	// Keep track of the time spent evolving
	//var start = time.Now()

	// Evolve the GA
	var (
		bar      *uiprogress.Bar
		progress *uiprogress.Progress
	)
	if verbose {
		progress = uiprogress.New()
		progress.Start()
		bar = progress.AddBar(est.NGenerations)
		bar.PrependCompleted()
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			var (
				best            = est.BestProgram()
				yTrainPred, err = best.Predict(est.XTrain, est.EvalMetric.NeedsProbabilities())
			)
			if err != nil {
				return "ERROR"
			}
			trainScore, err := est.EvalMetric.Apply(est.YTrain, yTrainPred, nil)
			if err != nil {
				return "ERROR"
			}
			var message = fmt.Sprintf("train %s: %.5f", est.EvalMetric.String(), trainScore)
			if est.XVal != nil && est.YVal != nil {
				yEvalPred, err := best.Predict(est.XVal, est.EvalMetric.NeedsProbabilities())
				if err != nil {
					return "ERROR"
				}
				evalScore, err := est.EvalMetric.Apply(est.YVal, yEvalPred, est.WVal)
				if err != nil {
					return "ERROR"
				}
				message += fmt.Sprintf(" val %s: %.5f", est.EvalMetric.String(), evalScore)
			}
			return message
		})
	}

	for i := 0; i < est.NGenerations; i++ {
		if verbose {
			bar.Incr()
		}

		// Make sure each tree has at least a height of 2
		for _, pop := range est.GA.Populations {
			for _, indi := range pop.Individuals {
				var prog = indi.Genome.(*Program)
				if prog.Tree.Height() < 2 {
					est.SubTreeMutation.Apply(&prog.Tree, pop.RNG)
					prog.Evaluate()
				}
			}
		}

		est.GA.Evolve()
	}

	if verbose {
		progress.Stop()
	}

	var best = est.BestProgram()

	return *best, nil
}

// Predict the output of a slice of features.
func (est Estimator) Predict(X [][]float64, predictProba bool) ([]float64, error) {
	var yPred, err = est.BestProgram().Predict(X, predictProba)
	if err != nil {
		return nil, err
	}
	return yPred, nil
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
		OperatorFactory{
			PConstant:   est.PConstant,
			NewConstant: est.newConstant,
			NewVariable: est.newVariable,
			NewFunction: est.newFunction,
		},
		rng,
	)
}

// newProgram can be used by gago to produce a new Genome.
func (est *Estimator) newProgram(rng *rand.Rand) gago.Genome {
	var prog = Program{
		Tree: est.newTree(rng),
		Task: Task{
			Metric:   est.LossMetric,
			NClasses: est.nClasses,
		},
		Estimator: est,
	}
	if prog.Task.multiClassification() {
		prog.DRS = &DynamicRangeSelection{}
	}
	return &prog
}

// newProgramTuner can be used by gago to produce a new Genome.
func (est *Estimator) newProgramTuner(rng *rand.Rand) gago.Genome {
	var (
		bestProg  = est.GA.HallOfFame[0].Genome.(*Program)
		progTuner = newProgramTuner(*bestProg)
	)
	progTuner.jitterConstants(rng)
	return &progTuner
}
