package xgp

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"github.com/MaxHalford/xgp/op"
)

// Evaluate is required to implement eaopt.Genome.
func (prog Program) Evaluate() (float64, error) {
	// For convenience
	gp := prog.GP
	// Run the training set through the Program
	var yPred, err = prog.Predict(gp.X, gp.LossMetric.NeedsProbabilities())
	if err != nil {
		return math.Inf(1), err
	}
	// Use the Metric defined in the GP
	fitness, err := gp.LossMetric.Apply(gp.Y, yPred, gp.W)
	if err != nil {
		return math.Inf(1), err
	}
	if math.IsNaN(fitness) {
		return math.Inf(1), nil
	}
	// Apply the parsimony coefficient
	if gp.ParsimonyCoeff != 0 {
		fitness += gp.ParsimonyCoeff * float64(op.CountOps(prog.Op))
	}
	return fitness, nil
}

// Mutate is required to implement eaopt.Genome.
func (prog *Program) Mutate(rng *rand.Rand) {
	defer func() { prog.Op = prog.Op.Simplify() }()
	var (
		pHoist   = prog.GP.PHoistMutation
		pSubtree = prog.GP.PSubtreeMutation
		pPoint   = prog.GP.PPointMutation
		dice     = rng.Float64() * (pHoist + pSubtree + pPoint)
	)
	// Apply hoist mutation
	if dice < pHoist {
		prog.Op = prog.GP.HoistMutation.Apply(prog.Op, rng)
		return
	}
	// Apply subtree mutation
	if dice < pHoist+pSubtree {
		prog.Op = prog.GP.SubtreeMutation.Apply(prog.Op, rng)
		return
	}
	// Apply point mutation
	prog.Op = prog.GP.PointMutation.Apply(prog.Op, rng)
}

// Crossover is required to implement eaopt.Genome.
func (prog *Program) Crossover(prog2 eaopt.Genome, rng *rand.Rand) {
	newOp1, newOp2 := prog.GP.SubtreeCrossover.Apply(prog.Op, prog2.(*Program).Op, rng)
	prog.Op = newOp1.Simplify()
	prog2.(*Program).Op = newOp2.Simplify()
}

// Clone is required to implement eaopt.Genome.
func (prog Program) Clone() eaopt.Genome {
	return &Program{
		Op: prog.Op,
		GP: prog.GP,
	}
}

// Custom genetic algorithm model.
type gaModel struct {
	selector   eaopt.Selector
	pMutate    float64
	pCrossover float64
}

// Apply is necessary to implement eaopt.Model.
func (mod gaModel) Apply(pop *eaopt.Population) error {
	var offsprings = make(eaopt.Individuals, len(pop.Individuals))
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

// Validate is necessary to implement eaopt.Model.
func (mod gaModel) Validate() error {
	return nil
}
