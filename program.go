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
	Estimator *Estimator             `json:"-"`
	DRS       *DynamicRangeSelection `json:"drs"`
}

// String representation of a Program.
func (prog Program) String() string {
	return prog.Tree.String()
}

// Clone a Program.
func (prog Program) clone() *Program {
	var clone = &Program{
		Tree:      prog.Tree.Clone(),
		Estimator: prog.Estimator,
	}
	if prog.DRS != nil {
		clone.DRS = prog.DRS.clone()
	}
	return clone
}

// PredictRow predicts the output of some features.
func (prog Program) PredictRow(x []float64) (float64, error) {
	var y = prog.Tree.EvaluateRow(x)
	if prog.DRS != nil {
		return prog.DRS.PredictRow(y), nil
	}
	return y, nil
}

// Predict predicts the output of a slice of features.
func (prog Program) Predict(X [][]float64) (yPred []float64, err error) {
	yPred, err = prog.Tree.EvaluateCols(X)
	if err != nil {
		return nil, err
	}
	if prog.DRS != nil {
		return prog.DRS.Predict(yPred), nil
	}
	return yPred, nil
}

// Evaluate is required to implement gago.Genome.
func (prog *Program) Evaluate() float64 {
	// Run the training set through the Program
	var yPred, err = prog.Tree.EvaluateCols(prog.Estimator.train.X)
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
