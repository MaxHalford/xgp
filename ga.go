package koza

import (
	"math"
	"math/rand"
	"sort"

	"github.com/MaxHalford/gago"
	"gonum.org/v1/gonum/floats"
)

// Evaluate is required to implement gago.Genome.
func (prog *Program) Evaluate() float64 {
	// Run the training set through the Program
	var yPred, err = prog.Predict(prog.Estimator.train.X, prog.Task.Metric.NeedsProbabilities())
	if err != nil {
		return math.Inf(1)
	}
	// Use the Metric defined in the Estimator
	fitness, err := prog.Task.Metric.Apply(prog.Estimator.train.Y, yPred, nil)
	if err != nil || math.IsNaN(fitness) {
		return math.Inf(1)
	}
	prog.Estimator.setBest(prog, fitness)
	if prog.Estimator.ParsimonyCoeff != 0 {
		fitness += prog.Estimator.ParsimonyCoeff * float64(prog.Tree.Height())
	}
	return fitness
}

// Mutate is required to implement gago.Genome.
func (prog *Program) Mutate(rng *rand.Rand) {
	// Select the mutation method through roulette wheel selection
	var (
		probs  = make([]float64, len(prog.Estimator.mutators))
		cumSum = make([]float64, len(probs))
	)
	for i, mutator := range prog.Estimator.mutators {
		probs[i] = mutator.p
	}
	floats.CumSum(cumSum, probs)
	var (
		p      = rng.Float64() * cumSum[len(cumSum)-1]
		i      = sort.SearchFloat64s(cumSum, p)
		method = prog.Estimator.mutators[i].method
	)
	// Apply the selected mutation method to the Tree
	method.Apply(prog.Tree, rng)
}

// Crossover is required to implement gago.Genome.
func (prog *Program) Crossover(prog2 gago.Genome, rng *rand.Rand) {
	// Select the mutation method through roulette wheel selection
	var (
		probs  = make([]float64, len(prog.Estimator.crossovers))
		cumSum = make([]float64, len(probs))
	)
	for i, crossover := range prog.Estimator.crossovers {
		probs[i] = crossover.p
	}
	floats.CumSum(cumSum, probs)
	var (
		p      = rng.Float64() * cumSum[len(cumSum)-1]
		i      = sort.SearchFloat64s(cumSum, p)
		method = prog.Estimator.crossovers[i].method
	)
	// Apply the selected mutation method to the Tree
	method.Apply(prog.Tree, prog2.(*Program).Tree, rng)
}

// Clone is required to implement gago.Genome.
func (prog Program) Clone() gago.Genome {
	var clone = prog.clone()
	return clone
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
		// Apply mutation
		if dice < mod.pMutate {
			offspring.Mutate(pop.RNG)
		}
		// Apply crossover
		if dice < (mod.pMutate + mod.pCrossover) {
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
