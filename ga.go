package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
)

// Evaluate is required to implement gago.Genome.
func (prog *Program) Evaluate() (float64, error) {
	// Simplify the Program's Tree to avoid unnecessary computations
	prog.Tree.Simplify()
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
		fitness += est.ParsimonyCoeff * float64(prog.Tree.Size())
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
		prog.Estimator.HoistMutation.Apply(&prog.Tree, rng)
		return
	}
	// Apply subtree mutation
	if dice < pHoist+pSubtree {
		prog.Estimator.SubtreeMutation.Apply(&prog.Tree, rng)
		return
	}
	// Apply point mutation
	prog.Estimator.PointMutation.Apply(&prog.Tree, rng)
}

// Crossover is required to implement gago.Genome.
func (prog *Program) Crossover(prog2 gago.Genome, rng *rand.Rand) {
	prog.Estimator.SubtreeCrossover.Apply(&prog.Tree, &prog2.(*Program).Tree, rng)
}

// Clone is required to implement gago.Genome.
func (prog Program) Clone() gago.Genome {
	var clone = prog.clone()
	return &clone
}

type gaModel struct {
	selector   gago.Selector
	pMutate    float64
	pCrossover float64
}

func (mod gaModel) Apply(pop *gago.Population) error {
	var offsprings = make(gago.Individuals, len(pop.Individuals))
	for i := range offsprings {
		// Select an individual
		selected, _, err := mod.selector.Apply(1, pop.Individuals, pop.RNG)
		if err != nil {
			return err
		}
		var offspring = selected[0]
		// Roll a dice to decide what to do
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

func (mod gaModel) Validate() error {
	return nil
}
