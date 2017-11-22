package tree

import (
	"encoding/json"
	"strings"
)

// MarshalJSON serializes a *Tree into JSON bytes.
func (tree *Tree) MarshalJSON() ([]byte, error) {
	var code = CodeDisplay{}.Apply(tree)
	return json.Marshal(&code)
}

// UnmarshalJSON parses JSON bytes into a *Tree.
func (tree *Tree) UnmarshalJSON(bytes []byte) error {
	var (
		code            = strings.TrimLeft(strings.TrimRight(string(bytes), `"`), `"`)
		parsedtree, err = ParseCode(code)
	)
	if err != nil {
		return err
	}
	*tree = *parsedtree
	return nil
}
