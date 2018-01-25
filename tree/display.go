package tree

import (
	"fmt"
	"strings"
)

// A Displayer outputs a string representation a Tree.
type Displayer interface {
	Apply(Tree) string
}

// DirDisplay outputs a directory like representation of a Tree.
//
//  root
//      branch 1
//          sub-branch 1.1
//          sub-branch 1.2
//      branch 2
//          sub-branch 2.1
//          sub-branch 2.2
//
type DirDisplay struct {
	TabSize int
}

// Apply directory-style display.
func (displayer DirDisplay) Apply(tree Tree) string {
	var (
		disp       func(tree *Tree, str string, depth int, carriage bool) string
		whitespace = strings.Repeat(" ", displayer.TabSize)
	)

	disp = func(tree *Tree, str string, depth int, carriage bool) string {
		str += strings.Repeat(whitespace, depth) + tree.op.String()
		if carriage {
			str += "\n"
		}
		for i := len(tree.branches) - 1; i >= 0; i-- {
			str = disp(tree.branches[i], str, depth+1, i > 0)
		}
		return str
	}

	return disp(&tree, "", 0, len(tree.branches) > 0)
}

// GraphvizDisplay outputs a Graphviz representation of a Tree. Each branch is
// indexed with a global counter and is labelled with the Tree's ToString
// method.
//
//  digraph G {
//      0 [label="root"];
//      0 -> 1
//      1 [label="branch 1"];
//      1 -> 2
//      2 [label="sub-branch 1.1
//      1 -> 3
//      3 [label="sub-branch 1.2"];
//      0 -> 4
//      4 [label="branch 2"];
//      4 -> 5
//      5 [label="sub-branch 2.1
//      4 -> 6
//      6 [label="sub-branch 2.2"];
//  }
//
type GraphvizDisplay struct{}

// Apply Graphviz display.
func (displayer GraphvizDisplay) Apply(tree Tree) string {
	var (
		counter int
		disp    func(tree *Tree, str string) string
	)
	disp = func(tree *Tree, str string) string {
		var idx = counter
		str += fmt.Sprintf("\t%d [label=\"%s\"];\n", idx, tree.op.String())
		for _, branch := range tree.branches {
			counter++
			str += fmt.Sprintf("\t%d -> %d;\n", idx, counter)
			str = disp(branch, str)
		}
		return str
	}
	var str = disp(&tree, "digraph G {\n")
	str += "}"
	return str
}

// CodeDisplay outputs an code-like representation of a Tree.
//
// pow(sum(X[0], X[1]), cos(2))
//
type CodeDisplay struct{}

// Apply CodeDisplay.
func (displayer CodeDisplay) Apply(tree Tree) string {
	switch len(tree.branches) {
	case 0:
		return tree.op.String()
	case 1:
		return fmt.Sprintf("%s(%s)", tree.op.String(), displayer.Apply(*tree.branches[0]))
	case 2:
		return fmt.Sprintf(
			"%s(%s, %s)",
			tree.op.String(),
			displayer.Apply(*tree.branches[0]),
			displayer.Apply(*tree.branches[1]),
		)
	default:
		return tree.op.String()
	}
}
