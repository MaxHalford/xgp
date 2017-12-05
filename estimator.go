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
// train Programs on a dataset. You shouldn't instantiate this struct directly,
// instead you should use the NewEstimator method.
type Estimator struct {
	ConstMin          float64
	ConstMax          float64
	EvalMetric        metrics.Metric // Defaults to LossMetric if nil
	LossMetric        metrics.Metric
	Functions         []tree.Operator // Should be kept in sync with fm
	Generations       int
	ParsimonyCoeff    float64
	TreeInitializer   tree.Initializer
	MaxHeight         int
	MinHeight         int
	PConstant         float64
	GA                *gago.GA
	TuningGA          *gago.GA
	TuningGenerations int
	CacheDuration     int

	PointMutation    tree.PointMutation
	PPointMutation   float64
	SubTreeMutation  tree.SubTreeMutation
	PSubTreeMutation float64
	HoistMutation    tree.HoistMutation
	PHoistMutation   float64
	SubTreeCrossover tree.SubTreeCrossover

	cache    *tree.Cache
	fm       map[int][]tree.Operator
	XTrain   [][]float64
	YTrain   []float64
	WTrain   []float64
	XVal     [][]float64
	YVal     []float64
	nClasses int
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
		evalScore, err := est.EvalMetric.Apply(est.YVal, yEvalPred, nil)
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
	XNames []string,
	verbose bool,
) error {

	// Set the training set
	est.XTrain = XTrain
	est.YTrain = YTrain
	est.WTrain = WTrain

	// Set the validation set
	est.XVal = XVal
	est.YVal = YVal

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

	// Display initial statistics
	if verbose {
		err := notifyProgress(est, 0, 0, os.Stdout)
		if err != nil {
			return err
		}
	}

	// Keep track of the time spent evolving
	var start = time.Now()

	// Evolve the GA
	for i := 0; i < est.Generations; i++ {

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
		if verbose {
			err := notifyProgress(est, i+1, time.Since(start), os.Stdout)
			if err != nil {
				return err
			}
		}
	}

	// Display the best program
	if verbose {
		fmt.Printf("Best program: %s\n", est.BestProgram())
	}

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
	return tree.Constant{Value: est.ConstMin + rng.Float64()*(est.ConstMax-est.ConstMin)}
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

// NewEstimator instantiates an Estimator. This method should be prefered over
// directly instantiating an Estimator.
func NewEstimator(
	constMax float64,
	constMin float64,
	evalMetric string,
	funcs string,
	generations int,
	lossMetric string,
	maxHeight int,
	minHeight int,
	nPops int,
	pConstant float64,
	pCrossover float64,
	pFull float64,
	pHoistMutation float64,
	pPointMutation float64,
	pSubTreeMutation float64,
	pTerminal float64,
	parsimonyCoeff float64,
	pointMutationRate float64,
	popSize int,
	seed int64,
	tuningGenerations int,
) (*Estimator, error) {

	// Determine the loss metric to use
	loss, err := metrics.GetMetric(lossMetric, 1)
	if err != nil {
		return nil, err
	}

	// Default the evaluation metric to the fitness metric if it's nil
	var eval metrics.Metric
	if evalMetric == "" {
		eval = loss
	} else {
		metric, err := metrics.GetMetric(evalMetric, 1)
		if err != nil {
			return nil, err
		}
		eval = metric
	}

	// The convention is to use a fitness metric which has to be minimized
	if loss.BiggerIsBetter() {
		loss = metrics.NegativeMetric{Metric: loss}
	}

	// Determine the functions to use
	fs, err := tree.ParseStringFuncs(funcs)
	if err != nil {
		return nil, err
	}

	// Instantiate an Estimator
	var estimator = &Estimator{
		ConstMin:       constMin,
		ConstMax:       constMax,
		EvalMetric:     eval,
		Functions:      fs,
		Generations:    generations,
		LossMetric:     loss,
		MaxHeight:      maxHeight,
		MinHeight:      minHeight,
		ParsimonyCoeff: parsimonyCoeff,
		PConstant:      pConstant,
		TreeInitializer: tree.RampedHaldAndHalfInitializer{
			PFull:           pFull,
			FullInitializer: tree.FullInitializer{},
			GrowInitializer: tree.GrowInitializer{
				PTerminal: pTerminal,
			},
		},
		TuningGenerations: tuningGenerations,
	}

	// Determine the random number generator of the GA
	var rng *rand.Rand
	if seed == 0 {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = rand.New(rand.NewSource(seed))
	}

	// Set the initial GA
	estimator.GA = &gago.GA{
		NewGenome: estimator.newProgram,
		NPops:     nPops,
		PopSize:   popSize,
		Model: gaModel{
			selector: gago.SelTournament{
				NContestants: 3,
			},
			pMutate:    pHoistMutation + pPointMutation + pSubTreeMutation,
			pCrossover: pCrossover,
		},
		RNG: rng,
	}

	// Build fm which maps arities to functions
	estimator.fm = make(map[int][]tree.Operator)
	for _, f := range estimator.Functions {
		var arity = f.Arity()
		if _, ok := estimator.fm[arity]; ok {
			estimator.fm[arity] = append(estimator.fm[arity], f)
		} else {
			estimator.fm[arity] = []tree.Operator{f}
		}
	}

	// Set crossover methods
	estimator.SubTreeCrossover = tree.SubTreeCrossover{
		Picker: tree.WeightedPicker{
			Weighting: tree.Weighting{
				PConstant: 0.1, // MAGIC
				PVariable: 0.1, // MAGIC
				PFunction: 0.8, // MAGIC
			},
		},
	}

	// Set mutation methods
	estimator.PointMutation = tree.PointMutation{
		Weighting: tree.Weighting{
			PConstant: pointMutationRate,
			PVariable: pointMutationRate,
			PFunction: pointMutationRate,
		},
		MutateOperator: func(op tree.Operator, rng *rand.Rand) tree.Operator {
			return estimator.mutateOperator(op, rng)
		},
	}
	estimator.PPointMutation = pPointMutation

	estimator.HoistMutation = tree.HoistMutation{
		Picker: tree.WeightedPicker{
			Weighting: tree.Weighting{
				PConstant: 0.1, // MAGIC
				PVariable: 0.1, // MAGIC
				PFunction: 0.8, // MAGIC
			},
		},
	}
	estimator.PHoistMutation = pHoistMutation

	estimator.SubTreeMutation = tree.SubTreeMutation{
		Crossover: tree.SubTreeCrossover{
			Picker: tree.WeightedPicker{
				Weighting: tree.Weighting{
					PConstant: 0.1, // MAGIC
					PVariable: 0.1, // MAGIC
					PFunction: 0.8, // MAGIC
				},
			},
		},
		NewTree: func(rng *rand.Rand) *tree.Tree {
			return estimator.newTree(rng)
		},
	}
	estimator.PSubTreeMutation = pSubTreeMutation

	return estimator, nil
}

// NewEstimatorWithDefaults call NewEstimator with default values.
func NewEstimatorWithDefaults() (*Estimator, error) {
	var estimator, err = NewEstimator(
		10,                // constMax
		-10,               // constMin
		"mae",             // evalMetric
		"sum,sub,mul,div", // funcs
		10,                // generations
		"mae",             // lossMetric
		6,                 // maxHeight
		3,                 // minHeight
		1,                 // nPops
		0.5,               // pConstant
		0.6,               // pCrossover
		0.5,               // pFull
		0.1,               // pHoistMutation
		0.1,               // pPointMutation
		0.1,               // PSubTreeMutation
		0.3,               // pTerminal
		0,                 // parsimonyCoeff
		0.1,               // pointMutationRate
		30,                // popSize
		42,                // seed
		0,                 // tuningGenerations
	)
	if err != nil {
		return nil, err
	}
	return estimator, nil
}
