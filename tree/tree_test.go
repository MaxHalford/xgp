package tree

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MaxHalford/xgp/op"
)

func TestTreeHeight(t *testing.T) {
	// Initial tr has no branches and thus has a height of 0
	var tr = Tree{}
	if tr.Height() != 0 {
		t.Errorf("Wrong height, expected %d got %d", 0, tr.Height())
	}
	// Add a branch
	tr.Branches = []*Tree{&Tree{}}
	if tr.Height() != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, tr.Height())
	}
	// Add another branch
	tr.Branches = append(tr.Branches, &Tree{})
	if tr.Height() != 1 {
		t.Errorf("Wrong height, expected %d got %d", 1, tr.Height())
	}
	// Add a sub-branch to the first branch
	tr.Branches[0].Branches = []*Tree{&Tree{}}
	if tr.Height() != 2 {
		t.Errorf("Wrong height, expected %d got %d", 2, tr.Height())
	}
}

func TestTreeSize(t *testing.T) {
	// Initial tr only has a root and thus has a single operator
	var tr = &Tree{}
	if tr.Size() != 1 {
		t.Errorf("Expected %d got %d", 1, tr.Size())
	}
	// Add a branch
	tr.Branches = []*Tree{&Tree{}}
	if tr.Size() != 2 {
		t.Errorf("Expected %d got %d", 2, tr.Size())
	}
	// Add a branch child
	tr.Branches = append(tr.Branches, &Tree{})
	if tr.Size() != 3 {
		t.Errorf("Expected %d got %d", 3, tr.Size())
	}
	// Add a sub-branch to the first branch
	tr.Branches[0].Branches = []*Tree{&Tree{}}
	if tr.Size() != 4 {
		t.Errorf("Expected %d got %d", 4, tr.Size())
	}
}

func TestTreeSimplify(t *testing.T) {
	var testCases = []struct {
		tr             Tree
		simplifiedTree Tree
	}{
		{
			tr:             MustParseCode("sum(1, 2)"),
			simplifiedTree: MustParseCode("3"),
		},
		{
			tr:             MustParseCode("sum(X[0], 2)"),
			simplifiedTree: MustParseCode("sum(X[0], 2)"),
		},
		{
			tr:             MustParseCode("sub(X[0], X[0])"),
			simplifiedTree: MustParseCode("0"),
		},
		{
			tr:             MustParseCode("div(X[0], X[0])"),
			simplifiedTree: MustParseCode("1"),
		},
		{
			tr:             MustParseCode("sum(sum(1, 2), sum(3, 4))"),
			simplifiedTree: MustParseCode("10"),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			tc.tr.Simplify()
			if !reflect.DeepEqual(tc.tr, tc.simplifiedTree) {
				t.Errorf("Expected %v, got %v", tc.simplifiedTree, tc.tr)
			}
		})
	}
}

func TestTreeMarshalJSON(t *testing.T) {
	var tr = Tree{
		Op: op.Sum{},
		Branches: []*Tree{
			&Tree{Op: op.Variable{0}},
			&Tree{Op: op.Constant{42}},
		},
	}
	var bytes, err = tr.MarshalJSON()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var tr2 = &Tree{}
	err = tr2.UnmarshalJSON(bytes)
	if tr2.String() != tr.String() {
		t.Errorf("Expected %v, got %v", tr.String(), tr2.String())
	}
}
