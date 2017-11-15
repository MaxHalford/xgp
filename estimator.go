package xgp

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/MaxHalford/gago"

	"github.com/MaxHalford/xgp/dataset"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/tree"
)

// An Estimator links all the different components together and can be used to
// train Programs on a dataset.
type Estimator struct {
	ConstMin          float64
	ConstMax          float64
	EvalMetric        metrics.Metric // Defaults to LossMetric if nil
	LossMetric        metrics.Metric
	Functions         []tree.Operator
	Generations       int
	ParsimonyCoeff    float64
	TreeInitializer   tree.Initializer
	PConstant         float64
	GA                *gago.GA
	TuningGA          *gago.GA
	TuningGenerations int
	CacheDuration     int

	bestProgram *Program
	bestLoss    float64
	mutex       sync.Mutex
	cache       *tree.Cache
	fm          map[int][]tree.Operator
	train       *dataset.Dataset
}

// BestProgram set an Estimator's bestProgram and bestLoss in a safe way.
func (est *Estimator) setBest(prog *Program, score float64) {
	est.mutex.Lock()
	if score < est.bestLoss {
		est.bestProgram = prog.clone()
		est.bestLoss = score
	}
	est.mutex.Unlock()
}

// BestProgram returns an Estimator's bestProgram in a safe way.
func (est Estimator) BestProgram() (*Program, error) {
	est.mutex.Lock()
	defer est.mutex.Unlock()
	if est.bestProgram == nil {
		return nil, errors.New("No best program has been set yet")
	}
	return est.bestProgram, nil
}

// Fit an Estimator to a dataset.Dataset.
func (est *Estimator) Fit(X [][]float64, Y []float64, XNames []string, verbose bool) error {

	// Set the train dataset so that the initial GA can be initialized
	var train, err = dataset.NewFromXY(X, Y, XNames, est.LossMetric.Classification())
	if err != nil {
		return err
	}
	est.train = train

	// Check the task is not multi-class classification
	if est.LossMetric.Classification() && est.train.NClasses() > 2 {
		return errors.New("Multi-class classification is not supported")
	}

	// Initialize the Estimator
	est.initialize()

	// Measure the progress through time
	var (
		totalDuration time.Duration
		writer        = tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		notify        = func(generation int, duration time.Duration) error {
			totalDuration += duration
			var best, err = est.BestProgram()
			if err != nil {
				return err
			}
			var (
				stats        = collectStats(est.GA)
				yPred, _     = best.Predict(est.train.X, est.EvalMetric.NeedsProbabilities())
				evalScore, _ = est.EvalMetric.Apply(est.train.Y, yPred, nil)
				message      = "[%d]\t%s: %.5f\tbest height: %d\tmean height: %.2f\tt_gen: %s\tt_total: %s\n"
			)
			fmt.Fprintf(
				writer,
				message,
				generation,
				est.EvalMetric.String(),
				evalScore,
				best.Tree.Height(),
				stats.avgHeight,
				duration,
				totalDuration,
			)
			writer.Flush()
			return nil
		}
	)

	// Display initial statistics
	if verbose {
		err := notify(0, 0)
		if err != nil {
			return err
		}
	}

	// Enhance the GA est.Generations times
	for i := 0; i < est.Generations; i++ {
		var start = time.Now()
		est.GA.Enhance()
		// Display current statistics
		if verbose {
			err := notify(i+1, time.Since(start))
			if err != nil {
				return err
			}
		}
	}

	// Display the best program
	best, err := est.BestProgram()
	if err != nil {
		return err
	}
	if verbose {
		fmt.Printf("Best program: %s\n", best)
	}

	return nil
}

// Predict the output of a slice of features.
func (est Estimator) Predict(X [][]float64, predictProba bool) ([]float64, error) {
	var bestProg, err = est.BestProgram()
	if err != nil {
		return nil, err
	}
	yPred, err := bestProg.Predict(X, predictProba)
	if err != nil {
		return nil, err
	}
	return yPred, nil
}

// Initialize an Estimator.
func (est *Estimator) initialize() error {
	// Initialize the best loss to +∞
	est.bestLoss = math.Inf(1)
	// Build fm which maps arities to functions
	est.fm = make(map[int][]tree.Operator)
	for _, f := range est.Functions {
		var arity = f.Arity()
		if _, ok := est.fm[arity]; ok {
			est.fm[arity] = append(est.fm[arity], f)
		} else {
			est.fm[arity] = []tree.Operator{f}
		}
	}
	// Initialize the GA
	est.GA.Initialize()
	// Initialize the cache
	if est.CacheDuration > 0 {
		est.cache = tree.NewCache()
	}
	return nil
}

func (est Estimator) newConstant(rng *rand.Rand) tree.Constant {
	return tree.Constant{Value: est.ConstMin + rng.Float64()*(est.ConstMax-est.ConstMin)}
}

func (est Estimator) newVariable(rng *rand.Rand) tree.Variable {
	return tree.Variable{Index: rng.Intn(est.train.NFeatures())}
}

func (est Estimator) newFunction(rng *rand.Rand) tree.Operator {
	return est.Functions[rng.Intn(len(est.Functions))]
}

func (est Estimator) newFunctionOfArity(arity int, rng *rand.Rand) tree.Operator {
	return est.fm[arity][rng.Intn(len(est.fm[arity]))]
}

// NewProgram can be used by gago to produce a new Genome.
func (est *Estimator) newProgram(rng *rand.Rand) gago.Genome {
	var (
		opFactory = tree.OperatorFactory{
			PConstant:   est.PConstant,
			NewConstant: est.newConstant,
			NewVariable: est.newVariable,
			NewFunction: est.newFunction,
		}
		prog = Program{
			Tree: est.TreeInitializer.Apply(opFactory, rng),
			Task: Task{
				Metric:   est.LossMetric,
				NClasses: est.train.NClasses(),
			},
			Estimator: est,
		}
	)
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
	lossMetric string,
	maxHeight int,
	minHeight int,
	generations int,
	parsimonyCoeff float64,
	pConstant float64,
	pTerminal float64,
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

	// Determine the random number generator
	var rng *rand.Rand
	if seed == 0 {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = rand.New(rand.NewSource(seed))
	}

	// Instantiate an Estimator
	var estimator = &Estimator{
		EvalMetric: eval,
		LossMetric: loss,
		ConstMin:   constMin,
		ConstMax:   constMax,
		PConstant:  pConstant,
		TreeInitializer: tree.RampedHaldAndHalfInitializer{
			MinHeight: minHeight,
			MaxHeight: maxHeight,
			PTerminal: pTerminal,
		},
		Functions:         fs,
		Generations:       generations,
		TuningGenerations: tuningGenerations,
		ParsimonyCoeff:    parsimonyCoeff,
	}

	// Set the initial GA
	estimator.GA = &gago.GA{
		NewGenome: estimator.newProgram,
		NPops:     1,
		PopSize:   1000,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
		RNG: rng,
	}

	return estimator, nil
}
