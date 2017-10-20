package xgp

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"

	"github.com/MaxHalford/gago"

	"github.com/MaxHalford/xgp/boosting"
	"github.com/MaxHalford/xgp/dataset"
	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/tree"
)

// An Estimator links all the different components together and can be used to
// train Programs on a dataset.
type Estimator struct {
	Metric            metrics.Metric
	Functions         []tree.Operator
	TreeInitializer   tree.Initializer
	ConstMin          float64
	ConstMax          float64
	PVariable         float64
	GA                *gago.GA
	TuningGA          *gago.GA
	Generations       int
	TuningGenerations int
	ParsimonyCoeff    float64

	fm    map[int][]tree.Operator
	train *dataset.Dataset
}

// BestProgram returns the best program an Estimator has produced.
func (est Estimator) BestProgram() (*Program, error) {
	var (
		GAOK       = !(est.GA == nil) && est.GA.Initialized()
		tuningGAOK = !(est.TuningGA == nil) && est.TuningGA.Initialized()
	)
	if !GAOK && !tuningGAOK {
		return nil, errors.New("No GA has been set")
	}
	if GAOK && !tuningGAOK {
		return est.GA.Best.Genome.(*Program), nil
	}
	if !GAOK && tuningGAOK {
		return est.TuningGA.Best.Genome.(*Program), nil
	}
	if est.GA.Best.Fitness < est.TuningGA.Best.Fitness {
		return est.GA.Best.Genome.(*Program), nil
	}
	return est.TuningGA.Best.Genome.(*ProgramTuner).Program, nil
}

// Fit an Estimator to a dataset.Dataset.
func (est *Estimator) Fit(X [][]float64, Y []float64, verbose bool) error {
	// Set the train dataset so that the initial GA can be initialized
	var train, err = dataset.NewDatasetXY(X, Y, est.Metric.Classification())
	if err != nil {
		return nil
	}
	est.train = train

	// // ((X[0])-(X[0]))/((X[0])*(X[1]))
	// var n = &Node{
	// 	Operator: Division{},
	// 	Children: []*Node{
	// 		&Node{
	// 			Operator: Difference{},
	// 			Children: []*Node{
	// 				&Node{Operator: Variable{0}},
	// 				&Node{Operator: Variable{0}},
	// 			},
	// 		},
	// 		&Node{
	// 			Operator: Product{},
	// 			Children: []*Node{
	// 				&Node{Operator: Variable{0}},
	// 				&Node{Operator: Variable{1}},
	// 			},
	// 		},
	// 	},
	// }
	// fmt.Println(n.evaluateXT(est.train.XT(), nil))

	var writer = tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	// Initialize the GA
	est.GA.Initialize()
	fmt.Println(est.GA.Populations[0].Individuals)

	// Display initial statistics
	var message = "[%d]\tBest fitness: %.5f\tMean size: %.2f\n"
	if verbose {
		var stats = collectStatistics(est.GA)
		fmt.Fprintf(writer, message, 0, est.GA.Best.Fitness, stats.avgHeight)
	}

	// // Enhance the GA est.Generations times
	// for i := 0; i < est.Generations; i++ {
	// 	est.GA.Enhance()
	// 	// Display statistics
	// 	if verbose {
	// 		var stats = CollectStatistics(est.GA)
	// 		fmt.Fprintf(
	// 			writer,
	// 			"[%d]\tBest fitness: %.5f\tMean size: %.2f\n",
	// 			i+1,
	// 			est.GA.Best.Fitness,
	// 			stats.AvgNNodes,
	// 		)
	// 		writer.Flush()
	// 	}
	// }

	// // Display the best equation
	// best, err := est.BestProgram()
	// fmt.Printf("Best equation: %s\n", best)

	// fmt.Println(pop.Individual[0])

	// // No need to continue if no tuning GA has been provided
	// if est.TuningGA == nil {
	// 	return nil
	// }

	// // Initialize the tuning GA
	// est.TuningGA.Initialize()

	// if verbose {
	// 	fmt.Printf("Number of constants to tune: %d\n", best.Root.NConstants())
	// }

	// for i := 0; i < est.TuningGenerations; i++ {
	// 	est.TuningGA.Enhance()
	// }

	return nil
}

// Predict the output of a slice of features.
func (est Estimator) Predict(X [][]float64) ([]float64, error) {
	var bestProg, err = est.BestProgram()
	if err != nil {
		return nil, err
	}
	yPred, err := bestProg.Predict(X)
	if err != nil {
		return nil, err
	}
	return yPred, nil
}

// functionMap returns a map mapping an integer to Operators that are in the
// Estimator's Functions and whose arity is equal to the integer. The result
// is memoized for subsequent calls.
func (est Estimator) functionMap() map[int][]tree.Operator {
	// Check if functionMap has already been computed
	if est.fm != nil {
		return est.fm
	}
	// Convert the slice of Operators into a map of Operators based on the
	// arities
	est.fm = make(map[int][]tree.Operator)
	for _, f := range est.Functions {
		var arity = f.Arity()
		if _, ok := est.fm[arity]; ok {
			est.fm[arity] = append(est.fm[arity], f)
		} else {
			est.fm[arity] = []tree.Operator{f}
		}
	}
	return est.fm
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
func (est *Estimator) NewProgram(rng *rand.Rand) gago.Genome {
	var (
		opFactory = tree.OperatorFactory{
			PVariable:   est.PVariable,
			NewConstant: est.newConstant,
			NewVariable: est.newVariable,
			NewFunction: est.newFunction,
		}
		prog = Program{
			Tree:      est.TreeInitializer.Apply(opFactory, rng),
			Estimator: est,
		}
	)
	if est.Metric.Classification() {
		prog.DRS = &DynamicRangeSelection{}
	}
	return &prog
}

// NewProgramTuner can be used by gago to produce a new Genome.
func (est *Estimator) NewProgramTuner(rng *rand.Rand) gago.Genome {
	var (
		bestProg  = est.GA.Best.Genome.(*Program)
		progTuner = newProgramTuner(bestProg)
	)
	progTuner.jitterConstants(rng)
	return &progTuner
}

// Learn method required to implement boosting.Learner.
func (est *Estimator) Learn(X [][]float64, Y []float64) (boosting.Predictor, error) {
	est.Fit(X, Y, false)
	var prog, err = est.BestProgram()
	if err != nil {
		return nil, err
	}
	return prog, nil
}
