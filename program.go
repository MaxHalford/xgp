package xgp

import (
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/dataframe"
	"github.com/MaxHalford/xgp/tree"
)

// A Program holds a tree composed of Nodes and also holds references to
// necessary information. A Program is simply an abstraction of top of a Node
// in order to not store the references in each Node of the tree.
type Program struct {
	Root      *Node
	Estimator *Estimator
}

// String representation of a Program.
func (prog Program) String() string {
	return prog.Root.String()
}

// PredictRow predicts the target of a row in a DataFrame.
func (prog Program) PredictRow(row []float64, transform func(float64) float64) float64 {
	var y = prog.Root.evaluate(row)
	if transform != nil {
		return transform(y)
	}
	return y
}

// PredictDataFrame predicts the target of each row in a DataFrame.
func (prog Program) PredictDataFrame(df *dataframe.DataFrame, transform func(float64) float64) []float64 {
	var (
		n, _  = df.Shape()
		yPred = make([]float64, n)
	)
	for i, row := range df.X {
		yPred[i] = prog.PredictRow(row, transform)
	}
	return yPred
}

// Implementation of the Genome interface from the gago package

// Evaluate method required to implement the Genome interface from the gago
// package.
func (prog Program) Evaluate() float64 {
	var (
		yPred      = prog.PredictDataFrame(prog.Estimator.DataFrame, prog.Estimator.Transform)
		fitness, _ = prog.Estimator.Metric.Apply(prog.Estimator.DataFrame.Y, yPred)
	)
	return fitness
}

// Mutate method required to implement the Genome interface from the gago
// package.
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

// Crossover method required to implement the Genome interface from the gago
// package.
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

// Clone method required to implement the Genome interface from the gago
// package.
func (prog Program) Clone() gago.Genome {
	return &Program{
		Root:      prog.Root.clone(),
		Estimator: prog.Estimator,
	}
}
