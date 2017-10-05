// Implementation of the Genome interface from the gago package.

package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/tree"
)

// Evaluate method required to implement gago.Genome.
func (prog *Program) Evaluate() float64 {
	// Run the training set through the Program
	var yPred, err = prog.Root.evaluateXT(
		prog.Estimator.train.XT(),
		prog.Estimator.nodeCache,
	)
	// If an error occurred during evaluation return +∞
	if err != nil {
		return math.Inf(1)
	}
	// Use dynamic range selection if applicable
	if prog.DRS != nil {
		prog.DRS.Fit(prog.Estimator.train.Y, yPred)
		yPred = prog.DRS.Predict(yPred)
	}
	// Use the Metric defined in the Estimator
	var fitness, _ = prog.Estimator.Metric.Apply(prog.Estimator.train.Y, yPred, nil)
	// If the Metric returned a NaN return +∞
	if math.IsNaN(fitness) {
		return math.Inf(1)
	}
	// Apply the parsimony coefficient
	if prog.Estimator.ParsimonyCoeff != 0 {
		fitness += prog.Estimator.ParsimonyCoeff * float64(tree.GetHeight(prog.Root))
	}
	return fitness
}

// Mutate method required to implement gago.Genome.
func (prog *Program) Mutate(rng *rand.Rand) {

	var (
		mutRate    = 0.3
		mutateNode func(node *Node)
	)

	mutateNode = func(node *Node) {
		if rng.Float64() < mutRate {
			switch node.Operator.(type) {
			case Constant:
				node.setOperator(prog.Estimator.newConstant(rng), rng)
			case Variable:
				node.setOperator(prog.Estimator.newVariable(rng), rng)
			default:
				node.setOperator(prog.Estimator.newFunctionOfArity(node.Operator.Arity(), rng), rng)
			}
		}
		for _, child := range node.Children {
			mutateNode(child)
		}
	}

	mutateNode(prog.Root)
}

// Crossover method required to implement gago.Genome.
func (prog Program) Crossover(prog2 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var (
		offspring1 = prog.Clone().(*Program)
		offspring2 = prog2.Clone().(*Program)
		picker     = tree.BernoulliPicker{P: 0.2}
		crossover  = tree.SubtreeCrossover{Picker: picker}
	)

	crossover.Apply(offspring1.Root, offspring2.Root, rng)

	return offspring1, offspring2
}

// Clone method required to implement gago.Genome.
func (prog Program) Clone() gago.Genome {
	var clone = prog.clone()
	return &clone
}
