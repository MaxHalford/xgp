package tree

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

// randTree is a convenience method that produces a random Tree for testing
// purposes.
func randTree(rng *rand.Rand) *Tree {
	var (
		init  = FullInitializer{Height: randInt(3, 5, rng)}
		funcs = []Operator{Cos{}, Sin{}, Sum{}, Difference{}, Product{}, Division{}}
		of    = OperatorFactory{
			PVariable:   0.5,
			NewConstant: func(rng *rand.Rand) Constant { return Constant{randFloat64(-5, 5, rng)} },
			NewVariable: func(rng *rand.Rand) Variable { return Variable{randInt(0, 5, rng)} },
			NewFunction: func(rng *rand.Rand) Operator { return funcs[rng.Intn(len(funcs))] },
		}
	)
	return init.Apply(of, rng)
}

func TestHeight(t *testing.T) {
	// Initial tree has no branches and thus has a height of 0
	var tree = Tree{}
	if tree.Height() != 0 {
		t.Errorf("Wrong height, expected %d got %d", 0, tree.Height())
	}
	// Add a branch
	tree.Branches = []*Tree{&Tree{}}
	if tree.Height() != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, tree.Height())
	}
	// Add another branch
	tree.Branches = append(tree.Branches, &Tree{})
	if tree.Height() != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, tree.Height())
	}
	// Add a sub-branch to the first branch
	tree.Branches[0].Branches = []*Tree{&Tree{}}
	if tree.Height() != 2 {
		t.Errorf("Wrong height, expected %d got %d", 2, tree.Height())
	}
}

func TestNOperators(t *testing.T) {
	// Initial tree only has a root and thus has a single operator
	var tree = &Tree{}
	if tree.NOperators() != 1 {
		t.Errorf("Expected %d got %d", 1, tree.NOperators())
	}
	// Add a branch
	tree.Branches = []*Tree{&Tree{}}
	if tree.NOperators() != 2 {
		t.Errorf("Expected %d got %d", 2, tree.NOperators())
	}
	// Add a branch child
	tree.Branches = append(tree.Branches, &Tree{})
	if tree.NOperators() != 3 {
		t.Errorf("Expected %d got %d", 3, tree.NOperators())
	}
	// Add a sub-branch to the first branch
	tree.Branches[0].Branches = []*Tree{&Tree{}}
	if tree.NOperators() != 4 {
		t.Errorf("Expected %d got %d", 4, tree.NOperators())
	}
}

func TestTreeSimplify(t *testing.T) {
	var testCases = []struct {
		tree       *Tree
		prunedTree *Tree
	}{
		{
			tree: &Tree{
				Operator: Sum{},
				Branches: []*Tree{
					&Tree{Operator: Constant{1}},
					&Tree{Operator: Constant{2}},
				},
			},
			prunedTree: &Tree{
				Operator: Constant{3},
			},
		},
		{
			tree: &Tree{
				Operator: Sum{},
				Branches: []*Tree{
					&Tree{Operator: Variable{0}},
					&Tree{Operator: Constant{42}},
				},
			},
			prunedTree: &Tree{
				Operator: Sum{},
				Branches: []*Tree{
					&Tree{Operator: Variable{0}},
					&Tree{Operator: Constant{42}},
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.tree.simplify()
			if !reflect.DeepEqual(tc.tree, tc.prunedTree) {
				t.Errorf("Expected %v, got %v", tc.prunedTree, tc.tree)
			}
		})
	}
}
