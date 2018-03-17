package xgp

import (
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
)

// A ProgramTuner optimizes a Program by tuning the Program's Constants.
type ProgramTuner struct {
	Program      Program
	ConstValues  []float64
	ConstSetters []ConstantSetter
}

// String representation of a ProgramTuner.
func (progTuner ProgramTuner) String() string {
	return progTuner.Program.String()
}

// newProgramTuner returns a ProgramTuner from a Program.
func newProgramTuner(prog Program) ProgramTuner {
	var (
		nConsts      = prog.Tree.NConstants()
		consts       = make([]float64, nConsts)
		constSetters = make([]ConstantSetter, nConsts)
		i            int
		addConst     = func(tr *tree.Tree, depth int) (stop bool) {
			if c, ok := tr.Operator().(op.Constant); ok {
				consts[i] = c.Value
				constSetters[i] = newConstantSetter(tr)
				i++
			}
			return
		}
		progTuner = ProgramTuner{Program: prog.clone()}
	)
	// Extract all the Constants from the Program
	progTuner.Program.Tree.Walk(addConst)
	progTuner.ConstValues = consts
	progTuner.ConstSetters = constSetters
	return progTuner
}

// Clone a ProgramTuner.
func (progTuner ProgramTuner) clone() ProgramTuner {
	var clone = newProgramTuner(progTuner.Program)
	copy(clone.ConstValues, progTuner.ConstValues)
	return clone
}

func (progTuner *ProgramTuner) setProgConstants() {
	for i, constValue := range progTuner.ConstValues {
		progTuner.ConstSetters[i](constValue)
	}
}

func (progTuner *ProgramTuner) jitterConstants(rng *rand.Rand) {
	for i, constValue := range progTuner.ConstValues {
		progTuner.ConstValues[i] += constValue * rng.NormFloat64()
	}
}

// Implementation of the Genome interface from the gago package

// Evaluate method required to implement gago.Genome.
func (progTuner *ProgramTuner) Evaluate() (float64, error) {
	progTuner.setProgConstants()
	return progTuner.Program.Evaluate()
}

// Mutate method required to implement gago.Genome.
func (progTuner *ProgramTuner) Mutate(rng *rand.Rand) {
	gago.MutNormalFloat64(progTuner.ConstValues, 0.8, rng)
}

// Crossover method required to implement gago.Genome.
func (progTuner *ProgramTuner) Crossover(progTuner2 gago.Genome, rng *rand.Rand) {
	gago.CrossUniformFloat64(
		progTuner.ConstValues,
		progTuner2.(*ProgramTuner).ConstValues,
		rng,
	)
}

// Clone method required to implement gago.Genome.
func (progTuner ProgramTuner) Clone() gago.Genome {
	var clone = progTuner.clone()
	return &clone
}
