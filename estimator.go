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
)

// An Estimator links all the different components together and can be used to
// train Programs on a dataset. You shouldn't instantiate this struct directly;
// instead you should use the Config struct and call it's NewEstimator method.
type Estimator struct {
	Config

	EvalMetric       metrics.Metric
	LossMetric       metrics.Metric
	Functions        []op.Operator
	Initializer      Initializer
	GA               *gago.GA
	PointMutation    PointMutation
	SubtreeMutation  SubtreeMutation
	HoistMutation    HoistMutation
	SubtreeCrossover SubtreeCrossover

	fm       map[uint][]op.Operator
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

// BestProgram returns the Estimator's best obtained Program.
func (est Estimator) BestProgram() Program {
	return *est.GA.HallOfFame[0].Genome.(*Program)
}

func (est Estimator) progress(start time.Time) string {
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
}

// polishBest takes the best Program and polishes it.
func (est *Estimator) polishBest() error {
	var (
		best          = *est.GA.HallOfFame[0].Genome.(*Program)
		polished, err = polishProgram(best, est.RNG)
	)
	if err != nil {
		return err
	}
	fitness, err := polished.Evaluate()
	if err != nil {
		return err
	}
	if fitness < est.GA.HallOfFame[0].Fitness {
		est.GA.HallOfFame[0].Genome = &polished
	}
	return nil
}

// Fit an Estimator to a dataset.
func (est *Estimator) Fit(
	// Required arguments
	XTrain [][]float64,
	YTrain []float64,
	// Optional arguments (can safely be nil)
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
		var steps = int(est.NGenerations)
		if est.PolishBest {
			steps++
		}
		bar = progress.AddBar(steps)
		bar.PrependCompleted()
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			return est.progress(start)
		})
	}

	// Make sure the progress bar will stop
	if verbose {
		defer func() { progress.Stop() }()
	}

	for i := uint(0); i < est.NGenerations; i++ {
		// Update progress
		if verbose {
			bar.Incr()
		}
		// Evolve a new generation
		err = est.GA.Evolve()
		if err != nil {
			return
		}
	}

	// Polish the best Program
	if est.PolishBest {
		err = est.polishBest()
		if err != nil {
			return
		}
		if verbose {
			bar.Incr()
		}
	}

	// Extract the best Program
	var best = est.BestProgram()

	return best, nil
}

func (est Estimator) newConst(rng *rand.Rand) op.Const {
	return op.Const{
		Value: est.Config.ConstMin + rng.Float64()*(est.Config.ConstMax-est.ConstMin),
	}
}

func (est Estimator) newVar(rng *rand.Rand) op.Var {
	return op.Var{Index: uint(rng.Intn(len(est.XTrain)))}
}

func (est Estimator) newFunction(rng *rand.Rand) op.Operator {
	return est.Functions[rng.Intn(len(est.Functions))]
}

func (est Estimator) newFunctionOfArity(arity uint, rng *rand.Rand) op.Operator {
	n := len(est.fm[arity])
	if n == 0 {
		return nil
	}
	return est.fm[arity][rng.Intn(n)]
}

func (est Estimator) newOperator(rng *rand.Rand) op.Operator {
	return est.Initializer.Apply(
		est.MinHeight,
		est.MaxHeight,
		func(leaf bool, rng *rand.Rand) op.Operator {
			if leaf {
				if rng.Float64() < est.PConst {
					return est.newConst(rng)
				}
				return est.newVar(rng)
			}
			return est.newFunction(rng)
		},
		rng,
	)
}

func (est Estimator) newProgram(rng *rand.Rand) Program {
	return Program{
		Op:        est.newOperator(rng),
		Estimator: &est,
	}
}

func (est Estimator) mutateOperator(operator op.Operator, rng *rand.Rand) op.Operator {
	switch operator.(type) {
	case op.Const:
		return op.Const{Value: operator.(op.Const).Value * rng.NormFloat64()}
	case op.Var:
		return est.newVar(rng)
	default:
		newOp := est.newFunctionOfArity(operator.Arity(), rng)
		// newFunctionOfArity might return nil if there are no available
		// operators of the given arity
		if newOp == nil {
			return operator
		}
		// Don't forget to set the new Operator's operands
		for i := uint(0); i < operator.Arity(); i++ {
			newOp = newOp.SetOperand(i, operator.Operand(i))
		}
		return newOp
	}
}
