package xgp

import (
	"testing"
)

func TestSubTreeCrossover(t *testing.T) {
	var (
		rng   = makeRNG()
		prog1 = randProg(rng)
		prog2 = randProg(rng)
	)
	// fmt.Println(prog1)
	// fmt.Println(prog2)
	SubTreeCrossover{}.Apply(&prog1, &prog2, rng)
	// fmt.Println(prog1)
	// fmt.Println(prog2)
}
