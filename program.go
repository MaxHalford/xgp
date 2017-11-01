package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/tree"
)

// A Program is simply an abstraction of top of a Tree.
type Program struct {
	Tree      *tree.Tree             `json:"tree"`
	Task      Task                   `json:"task"`
	DRS       *DynamicRangeSelection `json:"drs"`
	Estimator *Estimator             `json:"-"`
}

// String representation of a Program.
func (prog Program) String() string {
	return prog.Tree.String()
}

// Clone a Program.
func (prog Program) clone() *Program {
	var clone = &Program{
		Tree:      prog.Tree.Clone(),
		Task:      prog.Task,
		Estimator: prog.Estimator,
	}
	if prog.DRS != nil {
		clone.DRS = prog.DRS.clone()
	}
	return clone
}

// Predict predicts the output of a slice of features.
func (prog Program) Predict(X [][]float64, predictProba bool) (yPred []float64, err error) {
	yPred, err = prog.Tree.EvaluateCols(X, nil)
	if err != nil {
		return nil, err
	}
	// Binary classification
	if prog.Task.binaryClassification() {
		if predictProba {
			return sigmoid(yPred), nil
		}
		return binary(yPred), nil
	}
	// Multi-class classification
	if prog.Task.multiClassClassification() {
		return prog.DRS.Predict(yPred), nil
	}
	// Regression
	return yPred, nil
}

// Evaluate is required to implement gago.Genome.
func (prog *Program) Evaluate() float64 {
	// Run the training set through the Program
	var yPred, err = prog.Predict(prog.Estimator.train.X, prog.Task.Metric.NeedsProbabilities())
	// If an error occurred during evaluation return +∞
	if err != nil {
		return math.Inf(1)
	}
	// Use the Metric defined in the Estimator
	var fitness, _ = prog.Task.Metric.Apply(prog.Estimator.train.Y, yPred, nil)
	prog.Estimator.setBest(prog, fitness)
	// If the Metric returned a NaN return +∞
	if math.IsNaN(fitness) {
		return math.Inf(1)
	}
	// Apply the parsimony coefficient
	if prog.Estimator.ParsimonyCoeff != 0 {
		fitness += prog.Estimator.ParsimonyCoeff * float64(prog.Tree.Height())
	}
	return fitness
}

// Mutate is required to implement gago.Genome.
func (prog *Program) Mutate(rng *rand.Rand) {
	var newOp = func(op tree.Operator, rng *rand.Rand) tree.Operator {
		switch op.(type) {
		case tree.Constant:
			return prog.Estimator.newConstant(rng)
		case tree.Variable:
			return prog.Estimator.newVariable(rng)
		default:
			return prog.Estimator.newFunctionOfArity(op.Arity(), rng)
		}
	}
	tree.PointMutation{NewOperator: newOp, P: 0.2}.Apply(prog.Tree, rng)
}

// Crossover is required to implement gago.Genome.
func (prog Program) Crossover(prog2 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var (
		subTreeCrossover = tree.SubTreeCrossover{
			PConstant: 0.1,
			PVariable: 0.1,
			PFunction: 0.9,
		}
		offspring1 = prog.clone()
		offspring2 = prog2.(*Program).clone()
	)
	subTreeCrossover.Apply(offspring1.Tree, offspring2.Tree, rng)
	return offspring1, offspring2
}

// Clone is required to implement gago.Genome.
func (prog Program) Clone() gago.Genome {
	var clone = prog.clone()
	return clone
}
