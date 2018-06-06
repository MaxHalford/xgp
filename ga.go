package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/op"
)

// Evaluate is required to implement gago.Genome.
func (prog Program) Evaluate() (float64, error) {
	// For convenience
	est := prog.Estimator
	// Run the training set through the Program
	var yPred, err = prog.Predict(est.XTrain, est.LossMetric.NeedsProbabilities())
	if err != nil {
		return math.Inf(1), err
	}
	// Use the Metric defined in the Estimator
	fitness, err := est.LossMetric.Apply(est.YTrain, yPred, est.WTrain)
	if err != nil {
		return math.Inf(1), err
	}
	if math.IsNaN(fitness) {
		return math.Inf(1), nil
	}
	// Apply the parsimony coefficient
	if est.ParsimonyCoeff != 0 {
		fitness += est.ParsimonyCoeff * float64(op.CountOps(prog.Op))
	}
	return fitness, nil
}

// Mutate is required to implement gago.Genome.
func (prog *Program) Mutate(rng *rand.Rand) {
	var (
		pHoist   = prog.Estimator.PHoistMutation
		pSubtree = prog.Estimator.PSubtreeMutation
		pPoint   = prog.Estimator.PPointMutation
		dice     = rng.Float64() * (pHoist + pSubtree + pPoint)
	)
	// Apply hoist mutation
	if dice < pHoist {
		prog.Op = prog.Estimator.HoistMutation.Apply(prog.Op, rng)
		return
	}
	// Apply subtree mutation
	if dice < pHoist+pSubtree {
		prog.Op = prog.Estimator.SubtreeMutation.Apply(prog.Op, rng)
		return
	}
	// Apply point mutation
	prog.Op = prog.Estimator.PointMutation.Apply(prog.Op, rng).Simplify()
}

// Crossover is required to implement gago.Genome.
func (prog *Program) Crossover(prog2 gago.Genome, rng *rand.Rand) {
	newOp1, newOp2 := prog.Estimator.SubtreeCrossover.Apply(prog.Op, prog2.(*Program).Op, rng)
	prog.Op = newOp1.Simplify()
	prog2.(*Program).Op = newOp2.Simplify()
}

// Clone is required to implement gago.Genome.
func (prog Program) Clone() gago.Genome {
	return &Program{
		Op:        prog.Op,
		Estimator: prog.Estimator,
	}
}

// Custom genetic algorithm model.
type gaModel struct {
	selector   gago.Selector
	pMutate    float64
	pCrossover float64
}

// Apply is necessary to implement gago.Model.
func (mod gaModel) Apply(pop *gago.Population) error {
	var offsprings = make(gago.Individuals, len(pop.Individuals))
	for i := range offsprings {
		// Select an individual
		selected, _, err := mod.selector.Apply(1, pop.Individuals, pop.RNG)
		if err != nil {
			return err
		}
		var offspring = selected[0]
		// Roll a dice and decide what to do
		var dice = pop.RNG.Float64()
		if dice < mod.pMutate {
			// Mutation
			offspring.Mutate(pop.RNG)
		} else if dice < (mod.pMutate + mod.pCrossover) {
			// Crossover
			selected, _, err := mod.selector.Apply(1, pop.Individuals, pop.RNG)
			if err != nil {
				return err
			}
			offspring.Crossover(selected[0], pop.RNG)
		}
		// Insert the offsprings into the new population
		offsprings[i] = offspring
	}
	// Replace the population's individuals with the offsprings
	copy(pop.Individuals, offsprings)
	return nil
}

// Validate is necessary to implement gago.Model.
func (mod gaModel) Validate() error {
	return nil
}
