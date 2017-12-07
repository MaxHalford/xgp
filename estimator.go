package koza

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"

	"github.com/MaxHalford/gago"

	"github.com/MaxHalford/koza/metrics"
	"github.com/MaxHalford/koza/tree"
)

// An Estimator links all the different components together and can be used to
// train Programs on a dataset. You shouldn't instantiate this struct directly;
// instead you should use the Config struct.
type Estimator struct {
	Config

	EvalMetric       metrics.Metric
	LossMetric       metrics.Metric
	Functions        []tree.Operator
	TreeInitializer  tree.Initializer
	GA               *gago.GA
	TuningGA         *gago.GA
	PointMutation    tree.PointMutation
	SubTreeMutation  tree.SubTreeMutation
	HoistMutation    tree.HoistMutation
	SubTreeCrossover tree.SubTreeCrossover

	cache    *tree.Cache
	fm       map[int][]tree.Operator
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

func notifyProgress(est *Estimator, generation int, duration time.Duration, w io.Writer) error {
	var (
		writer = tabwriter.NewWriter(w, 0, 0, 4, ' ', 0)
		best   = est.BestProgram()
		stats  = collectStats(est.GA)
	)

	// Get the score on the training set
	yTrainPred, err := best.Predict(est.XTrain, est.EvalMetric.NeedsProbabilities())
	if err != nil {
		return err
	}
	trainScore, err := est.EvalMetric.Apply(est.YTrain, yTrainPred, nil)
	if err != nil {
		return err
	}

	// Get the score on the evaluation set
	var evalScoreStr = "/"
	if est.XVal != nil && est.YVal != nil {
		yEvalPred, err := best.Predict(est.XVal, est.EvalMetric.NeedsProbabilities())
		if err != nil {
			return err
		}
		evalScore, err := est.EvalMetric.Apply(est.YVal, yEvalPred, est.WVal)
		if err != nil {
			return err
		}
		evalScoreStr = fmt.Sprintf("%.5f", evalScore)
	}

	// Produce the output message
	var message = "[%d]\ttrain %s: %.5f\tval %s: %s\tbest size: %d\tmean size: %.2f\tduration: %s\n"
	fmt.Fprintf(
		writer,
		message,
		generation,
		est.EvalMetric.String(),
		trainScore,
		est.EvalMetric.String(),
		evalScoreStr,
		best.Tree.NOperators(),
		stats.avgNOperators,
		fmtDuration(duration),
	)
	writer.Flush()

	return nil
}

// Fit an Estimator to a dataset.Dataset.
func (est *Estimator) Fit(
	XTrain [][]float64,
	YTrain []float64,
	WTrain []float64,
	XVal [][]float64,
	YVal []float64,
	WVal []float64,
	notifyEvery uint,
) error {

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
			return errors.New("Multi-class classification is not supported")
		}
	}

	// Initialize the GA
	est.GA.Initialize()

	// Initialize the cache
	if est.CacheDuration > 0 {
		est.cache = tree.NewCache()
	}

	fmt.Println(est)

	// Display initial statistics
	if notifyEvery > 0 {
		err := notifyProgress(est, 0, 0, os.Stdout)
		if err != nil {
			return err
		}
	}

	// Keep track of the time spent evolving
	var start = time.Now()

	// Evolve the GA
	for i := 0; i < est.NGenerations; i++ {

		// Make sure each tree has at least a height of 2
		for _, pop := range est.GA.Populations {
			for _, indi := range pop.Individuals {
				var prog = indi.Genome.(*Program)
				if prog.Tree.Height() < 2 {
					est.SubTreeMutation.Apply(prog.Tree, pop.RNG)
					prog.Evaluate()
				}
			}
		}

		est.GA.Evolve()

		// Display current statistics
		if notifyEvery > 0 && uint(i+1)%notifyEvery == 0 {
			err := notifyProgress(est, i+1, time.Since(start), os.Stdout)
			if err != nil {
				return err
			}
		}
	}

	// Display the best program
	fmt.Printf("Best program: %s\n", est.BestProgram())

	return nil
}

// Predict the output of a slice of features.
func (est Estimator) Predict(X [][]float64, predictProba bool) ([]float64, error) {
	var yPred, err = est.BestProgram().Predict(X, predictProba)
	if err != nil {
		return nil, err
	}
	return yPred, nil
}

func (est Estimator) newConstant(rng *rand.Rand) tree.Constant {
	return tree.Constant{
		Value: est.Config.ConstMin + rng.Float64()*(est.Config.ConstMax-est.ConstMin),
	}
}

func (est Estimator) newVariable(rng *rand.Rand) tree.Variable {
	return tree.Variable{Index: rng.Intn(len(est.XTrain))}
}

func (est Estimator) newFunction(rng *rand.Rand) tree.Operator {
	return est.Functions[rng.Intn(len(est.Functions))]
}

func (est Estimator) newFunctionOfArity(arity int, rng *rand.Rand) tree.Operator {
	return est.fm[arity][rng.Intn(len(est.fm[arity]))]
}

func (est Estimator) mutateOperator(op tree.Operator, rng *rand.Rand) tree.Operator {
	switch op.(type) {
	case tree.Constant:
		return tree.Constant{Value: op.(tree.Constant).Value * rng.NormFloat64()}
	case tree.Variable:
		return est.newVariable(rng)
	default:
		return est.newFunctionOfArity(op.Arity(), rng)
	}
}

func (est Estimator) newTree(rng *rand.Rand) *tree.Tree {
	return est.TreeInitializer.Apply(
		est.MinHeight,
		est.MaxHeight,
		tree.OperatorFactory{
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
	if est.LossMetric.Classification() {
		prog.DRS = &DynamicRangeSelection{}
	}
	return &prog
}

// newProgramTuner can be used by gago to produce a new Genome.
func (est *Estimator) newProgramTuner(rng *rand.Rand) gago.Genome {
	var (
		bestProg  = est.GA.HallOfFame[0].Genome.(*Program)
		progTuner = newProgramTuner(bestProg)
	)
	progTuner.jitterConstants(rng)
	return &progTuner
}
