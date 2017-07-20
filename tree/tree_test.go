package tree

import (
	"testing"
)

func TestGetHeight(t *testing.T) {
	// Initial tree has no branches and thus has a height of 0
	var tree = &TestTree{}
	if GetHeight(tree) != 0 {
		t.Errorf("Wrong height, expected %d got %d", 0, GetHeight(tree))
	}
	// Add a branch
	tree.Branches = []*TestTree{&TestTree{}}
	if GetHeight(tree) != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, GetHeight(tree))
	}
	// Add a branch child
	tree.Branches = append(tree.Branches, &TestTree{})
	if GetHeight(tree) != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, GetHeight(tree))
	}
	// Add a sub-branch to the first branch
	tree.Branches[0].Branches = []*TestTree{&TestTree{}}
	if GetHeight(tree) != 2 {
		t.Errorf("Wrong height, expected %d got %d", 2, GetHeight(tree))
	}
}

func TestGetNNodes(t *testing.T) {
	// Initial tree only has a root and thus has 1 single node
	var tree = &TestTree{}
	if GetNNodes(tree) != 1 {
		t.Errorf("Wrong number of nodes, expected %d got %d", 1, GetNNodes(tree))
	}
	// Add a branch
	tree.Branches = []*TestTree{&TestTree{}}
	if GetNNodes(tree) != 2 {
		t.Errorf("Wrong number of nodes, expected %d got %d", 2, GetNNodes(tree))
	}
	// Add a branch child
	tree.Branches = append(tree.Branches, &TestTree{})
	if GetNNodes(tree) != 3 {
		t.Errorf("Wrong number of nodes, expected %d got %d", 3, GetNNodes(tree))
	}
	// Add a sub-branch to the first branch
	tree.Branches[0].Branches = []*TestTree{&TestTree{}}
	if GetNNodes(tree) != 4 {
		t.Errorf("Wrong number of nodes, expected %d got %d", 4, GetNNodes(tree))
	}
}
