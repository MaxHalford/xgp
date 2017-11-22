package tree

import (
	"encoding/json"
	"testing"
)

func TestTreeJSONEncodeDecode(t *testing.T) {
	var initialtree = mustParseCode("sum(42, X[1])")

	// Serialize the initial Tree
	var bytes, err = json.Marshal(initialtree)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Parse the bytes into a new Tree
	var newtree *Tree
	err = json.Unmarshal(bytes, &newtree)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	// Compare the new Tree with the initial Tree
	var check func(n1, n2 *Tree)
	check = func(n1, n2 *Tree) {
		if n1.Operator.String() != n2.Operator.String() {
			t.Errorf("Operator mismatch: %s != %s", n1.Operator.String(), n2.Operator.String())
		}
		for i := range n1.Branches {
			check(n1.Branches[i], n2.Branches[i])
		}
	}
	check(newtree, initialtree)
}
