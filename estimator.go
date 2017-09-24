package xgp

import (
	"errors"
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/boosting"
	"github.com/MaxHalford/xgp/dataset"
	"github.com/MaxHalford/xgp/metrics"
)

// An Estimator links all the different components together and can be used to
// train Programs on a dataset.
type Estimator struct {
	Metric            metrics.Metric
	Transform         Transform
	PVariable         float64         // Probability of producing a Variable when creating a terminal Node
	NodeInitializer   NodeInitializer // Method for producing new Program trees
	Functions         []Operator
	GA                *gago.GA
	TuningGA          *gago.GA
	Generations       int
	TuningGenerations int
	ProgressChan      chan float64

	fm           map[int][]Operator // Function map
	dataset      *dataset.Dataset
	targetMean   float64 // Used for generating new Constants
	targetStdDev float64 // Used for generating new Constants
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
	return &est.TuningGA.Best.Genome.(*ProgramTuner).Program, nil
}

// Fit an Estimator to a dataset.Dataset.
func (est *Estimator) Fit(X [][]float64, Y []float64) error {
	// Set the dataset so that the initial GA can be initialized
	var dataset, err = dataset.NewDatasetXY(X, Y, est.Metric.Classification())
	if err != nil {
		return nil
	}
	est.dataset = dataset

	// Compute the target average and standard deviation to help produce
	// meaningful Constants
	est.targetMean = meanFloat64s(est.dataset.Y)
	est.targetStdDev = math.Pow(varianceFloat64s(est.dataset.Y), 0.5)

	// Run the initial GA
	est.GA.Initialize()
	for i := 0; i < est.Generations; i++ {
		est.GA.Enhance()
		if est.ProgressChan != nil {
			est.ProgressChan <- est.GA.CurrentBest.Fitness
		}
	}

	// Run the tuning GA
	if est.TuningGA == nil {
		return nil
	}
	est.TuningGA.Initialize()
	for i := 0; i < est.TuningGenerations; i++ {
		est.TuningGA.Enhance()
		if est.ProgressChan != nil {
			est.ProgressChan <- est.TuningGA.CurrentBest.Fitness
		}
	}

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

// newConstant returns a Constant whose value is sampled from a normal
// distribution based on the Estimator's dataset's y slice.
func (est Estimator) newConstant(rng *rand.Rand) Constant {
	return Constant{
		Value: est.targetMean + rng.NormFloat64()*est.targetStdDev,
	}
}

// functionMap returns a map mapping an integer to Operators that are in the
// Estimator's Functions and whose arity is equal to the integer. The result
// is memoized for subsequent calls.
func (est Estimator) functionMap() map[int][]Operator {
	// Check if functionMap has already been computed
	if est.fm != nil {
		return est.fm
	}
	// Convert the slice of Operators into a map of Operators based on the
	// arities
	est.fm = make(map[int][]Operator)
	for _, f := range est.Functions {
		var arity = f.Arity()
		if _, ok := est.fm[arity]; ok {
			est.fm[arity] = append(est.fm[arity], f)
		} else {
			est.fm[arity] = []Operator{f}
		}
	}
	return est.fm
}

// newVariable returns a Variable with an index in range [0, p) where p is the
// number of explanatory variables in the Estimator's dataset.dataset.
func (est Estimator) newVariable(rng *rand.Rand) Variable {
	return Variable{
		Index: rng.Intn(est.dataset.NFeatures()),
	}
}

func (est Estimator) newFunctionOfArity(arity int, rng *rand.Rand) Operator {
	var fm = est.functionMap()
	return fm[arity][rng.Intn(len(fm[arity]))]
}

// newOperator generates a random *Node. If terminal is true then a Constant or
// a Variable is returned. If not a Function is returned.
func (est Estimator) newOperator(terminal bool, rng *rand.Rand) Operator {
	if terminal {
		if rng.Float64() < est.PVariable {
			return est.newVariable(rng)
		}
		return est.newConstant(rng)
	}
	return est.newFunctionOfArity(2, rng)
}

// NewProgram can be used by gago to produce a new Genome.
func (est *Estimator) NewProgram(rng *rand.Rand) gago.Genome {
	return &Program{
		Root:      est.NodeInitializer.Apply(est.newOperator, rng),
		Estimator: est,
	}
}

// NewProgramTuner can be used by gago to produce a new Genome.
func (est *Estimator) NewProgramTuner(rng *rand.Rand) gago.Genome {
	var (
		bestProg  = est.GA.Best.Genome.(*Program)
		progTuner = newProgramTuner(*bestProg)
	)
	progTuner.jitterConstants(rng)
	return &progTuner
}

// Learn method required to implement boosting.Learner.
func (est *Estimator) Learn(X [][]float64, Y []float64) (boosting.Predictor, error) {
	est.Fit(X, Y)
	var prog, err = est.BestProgram()
	if err != nil {
		return nil, err
	}
	return prog, nil
}
