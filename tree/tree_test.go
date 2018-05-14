package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTreeHeight(t *testing.T) {
	// Initial tree has no branches and thus has a height of 0
	var tree = Tree{}
	if tree.Height() != 0 {
		t.Errorf("Wrong height, expected %d got %d", 0, tree.Height())
	}
	// Add a branch
	tree.branches = []*Tree{&Tree{}}
	if tree.Height() != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, tree.Height())
	}
	// Add another branch
	tree.branches = append(tree.branches, &Tree{})
	if tree.Height() != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, tree.Height())
	}
	// Add a sub-branch to the first branch
	tree.branches[0].branches = []*Tree{&Tree{}}
	if tree.Height() != 2 {
		t.Errorf("Wrong height, expected %d got %d", 2, tree.Height())
	}
}

func TestTreeSize(t *testing.T) {
	// Initial tree only has a root and thus has a single operator
	var tree = &Tree{}
	if tree.Size() != 1 {
		t.Errorf("Expected %d got %d", 1, tree.Size())
	}
	// Add a branch
	tree.branches = []*Tree{&Tree{}}
	if tree.Size() != 2 {
		t.Errorf("Expected %d got %d", 2, tree.Size())
	}
	// Add a branch child
	tree.branches = append(tree.branches, &Tree{})
	if tree.Size() != 3 {
		t.Errorf("Expected %d got %d", 3, tree.Size())
	}
	// Add a sub-branch to the first branch
	tree.branches[0].branches = []*Tree{&Tree{}}
	if tree.Size() != 4 {
		t.Errorf("Expected %d got %d", 4, tree.Size())
	}
}

func TestTreeSimplify(t *testing.T) {
	var testCases = []struct {
		tree           Tree
		simplifiedTree Tree
	}{
		{
			tree:           MustParseCode("sum(1, 2)"),
			simplifiedTree: MustParseCode("3"),
		},
		{
			tree:           MustParseCode("sum(X[0], 2)"),
			simplifiedTree: MustParseCode("sum(X[0], 2)"),
		},
		{
			tree:           MustParseCode("sub(X[0], X[0])"),
			simplifiedTree: MustParseCode("0"),
		},
		{
			tree:           MustParseCode("div(X[0], X[0])"),
			simplifiedTree: MustParseCode("1"),
		},
		{
			tree:           MustParseCode("sum(sum(1, 2), sum(3, 4))"),
			simplifiedTree: MustParseCode("10"),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.tree.Simplify()
			if !reflect.DeepEqual(tc.tree, tc.simplifiedTree) {
				t.Errorf("Expected %v, got %v", tc.simplifiedTree, tc.tree)
			}
		})
	}
}
