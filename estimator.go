package main

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/tiago/dataframe"
	"github.com/MaxHalford/tiago/metric"
)

type Estimator struct {
	DataFrame       *dataframe.DataFrame
	Metric          metric.Metric
	Activation      func(float64) float64
	PVariable       float64         // Probability of producing a Variable when creating a terminal Node
	NodeInitializer NodeInitializer // Method for producing new Program trees
	FunctionSet     map[int][]Operator
	GA              gago.GA

	// Fields that are generated at runtime
	targetMean                  float64 // Used for generating new Constants
	targetStdDev                float64 // Used for generating new Constants
	bestScore                   float64 // Used for determining early stopping
	generationsSinceImprovement int     // Used for determining early stopping
}

func (est *Estimator) Initialize() {
	// Compute the target average and standard deviation to help produce
	// meaningful Constants
	est.targetMean = meanFloat64s(est.DataFrame.Y)
	est.targetStdDev = math.Pow(varianceFloat64s(est.DataFrame.Y), 0.5)

	// Initialize the genetic algorithm
	est.GA.Initialize()
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

// newNode generates a random *Node. If leaf is true then newNode returns a
// Constant or a Variable by flipping a coin against PVariable. If leaf is not
// true then newNode returns a Function.
func (est Estimator) newNode(terminal bool, rng *rand.Rand) *Node {
	var operator Operator
	if terminal {
		if rng.Float64() < est.PVariable {
			operator = est.newVariable(rng)
		} else {
			operator = est.newConstant(rng)
		}
		return &Node{Operator: operator}
	} else {
		operator = est.newFunctionOfArity(2, rng)
		return &Node{
			Operator: operator,
			Children: make([]*Node, operator.Arity()),
		}
	}
}

// NewProgram can be used by gago to produce a new Genome.
func (est *Estimator) NewProgram(rng *rand.Rand) gago.Genome {
	return &Program{
		Root:      est.NodeInitializer.Apply(est.newNode, rng),
		Estimator: est,
	}
}
