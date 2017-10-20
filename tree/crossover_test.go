package tree

import (
	"testing"
)

func TestSubTreeCrossover(t *testing.T) {
	var (
		rng   = newRand()
		left  = randTree(rng)
		right = randTree(rng)
	)
	// fmt.Println(prog1)
	// fmt.Println(prog2)
	SubTreeCrossover{}.Apply(left, right, rng)
	// fmt.Println(prog1)
	// fmt.Println(prog2)
}
