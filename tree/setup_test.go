package tree

import "fmt"

type TestTree struct {
	Value    float64
	Branches []*TestTree
}

func (tree TestTree) sum() float64 {
	var sum = tree.Value
	for _, branch := range tree.Branches {
		sum += branch.sum()
	}
	return sum
}

// Implement Tree interface

func (tree *TestTree) NBranches() int {
	return len(tree.Branches)
}

func (tree *TestTree) GetBranch(i int) Tree {
	return tree.Branches[i]
}

func (tree *TestTree) Swap(otherTree Tree) {
	*tree, *otherTree.(*TestTree) = *otherTree.(*TestTree), *tree
}

func (tree *TestTree) ToString() string {
	return fmt.Sprintf("%.2f", tree.Value)
}
