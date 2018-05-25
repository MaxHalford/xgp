package xgp

import (
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
)

// A ConstantSetter can replace a tree's Operator with a Constant.
type ConstantSetter func(value float64)

// NewConstantSetter returns a ConstantSetter that can be used as a callback
// to set a Program's Operator to a given Constant.
func newConstantSetter(tr *tree.Tree) ConstantSetter {
	return func(value float64) {
		tr.Op = op.Constant{Value: value}
	}
}

// A ProgramPolish optimizes a Program by tuning it's Constants.
type ProgramPolish struct {
	Program      Program
	ConstValues  []float64
	ConstSetters []ConstantSetter
}

// String representation of a ProgramPolish.
func (progPolish ProgramPolish) String() string {
	return progPolish.Program.String()
}

// newProgramPolish returns a ProgramPolish from a Program.
func newProgramPolish(prog Program) ProgramPolish {
	var (
		nConsts      = prog.Tree.NConstants()
		consts       = make([]float64, nConsts)
		constSetters = make([]ConstantSetter, nConsts)
		i            int
		addConst     = func(tr *tree.Tree, depth int) (stop bool) {
			if c, ok := tr.Op.(op.Constant); ok {
				consts[i] = c.Value
				constSetters[i] = newConstantSetter(tr)
				i++
			}
			return
		}
		progPolish = ProgramPolish{Program: prog.clone()}
	)
	// Extract all the Constants from the Program
	progPolish.Program.Tree.Walk(addConst)
	progPolish.ConstValues = consts
	progPolish.ConstSetters = constSetters
	return progPolish
}

// Clone a ProgramPolish.
func (progPolish ProgramPolish) clone() ProgramPolish {
	var clone = newProgramPolish(progPolish.Program)
	copy(clone.ConstValues, progPolish.ConstValues)
	return clone
}

func (progPolish *ProgramPolish) setProgConstants() {
	for i, constValue := range progPolish.ConstValues {
		progPolish.ConstSetters[i](constValue)
	}
}

func (progPolish *ProgramPolish) jitterConstants(rng *rand.Rand) {
	for i, constValue := range progPolish.ConstValues {
		progPolish.ConstValues[i] += constValue * rng.NormFloat64()
	}
}

// Implementation of the Genome interface from the gago package

// Evaluate method required to implement gago.Genome.
func (progPolish *ProgramPolish) Evaluate() (float64, error) {
	progPolish.setProgConstants()
	return progPolish.Program.Evaluate()
}

// Mutate method required to implement gago.Genome.
func (progPolish *ProgramPolish) Mutate(rng *rand.Rand) {
	gago.MutNormalFloat64(progPolish.ConstValues, 0.8, rng)
}

// Crossover method required to implement gago.Genome.
func (progPolish *ProgramPolish) Crossover(progPolish2 gago.Genome, rng *rand.Rand) {
	gago.CrossUniformFloat64(
		progPolish.ConstValues,
		progPolish2.(*ProgramPolish).ConstValues,
		rng,
	)
}

// Clone method required to implement gago.Genome.
func (progPolish ProgramPolish) Clone() gago.Genome {
	var clone = progPolish.clone()
	return &clone
}
