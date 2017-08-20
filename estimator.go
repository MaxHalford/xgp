package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/MaxHalford/xgp/metric"
)

// An Estimator links all the different components together and can be used to
// train Programs on a DataFrame.
type Estimator struct {
	DataFrame       *dataframe.DataFrame
	Metric          metric.Metric
	Transform       Transform
	PVariable       float64         // Probability of producing a Variable when creating a terminal Node
	NodeInitializer NodeInitializer // Method for producing new Program trees
	FunctionSet     map[int][]Operator
	GA              *gago.GA
	TuningGA        *gago.GA

	// Fields that are generated at runtime
	targetMean                  float64 // Used for generating new Constants
	targetStdDev                float64 // Used for generating new Constants
	bestScore                   float64 // Used for determining early stopping
	generationsSinceImprovement int     // Used for determining early stopping
}

// Initialize an Estimator.
func (est *Estimator) Initialize() {
	// Compute the target average and standard deviation to help produce
	// meaningful Constants
	est.targetMean = meanFloat64s(est.DataFrame.Y)
	est.targetStdDev = math.Pow(varianceFloat64s(est.DataFrame.Y), 0.5)

	// Initialize the genetic algorithm
	est.GA.Initialize()

	// Initialize the tuning genetic algorithm
	if est.TuningGA != nil {
		est.TuningGA.Initialize()
	}
}

func (est Estimator) newConstant(rng *rand.Rand) Constant {
	return Constant{
		Value: est.targetMean + rng.NormFloat64()*est.targetStdDev,
	}
}

// newVariable returns a Variable with an index in range [0, p) where p is the
// number of explanatory variables in the Estimator's DataFrame.
func (est Estimator) newVariable(rng *rand.Rand) Variable {
	return Variable{
		Index: rng.Intn(est.DataFrame.NFeatures()),
	}
}

func (est Estimator) newFunctionOfArity(arity int, rng *rand.Rand) Operator {
	return est.FunctionSet[arity][rng.Intn(len(est.FunctionSet[arity]))]
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
	return &progTuner
}

// Fit an Estimator to find an optimal Program.
func (est *Estimator) Fit() {
}
