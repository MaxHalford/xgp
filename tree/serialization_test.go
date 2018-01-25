package tree

/*func TestTreeJSONEncodeDecode(t *testing.T) {
	var initialtree = MustParseCode("sum(42, X[1])")

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
	var check func(n1, n2 Tree)
	check = func(n1, n2 Tree) {
		if n1.op.String() != n2.op.String() {
			t.Errorf("Operator mismatch: %s != %s", n1.op.String(), n2.op.String())
		}
		for i := range n1.branches {
			check(n1.branches[i], n2.branches[i])
		}
	}
	check(newtree, initialtree)
}
*/
