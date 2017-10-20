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
func (displayer DirDisplay) Apply(tree *Tree) string {
	var (
		disp       func(tree *Tree, str string, depth int, carriage bool) string
		whitespace = strings.Repeat(" ", displayer.TabSize)
	)

	disp = func(tree *Tree, str string, depth int, carriage bool) string {
		str += strings.Repeat(whitespace, depth) + tree.Operator.String()
		if carriage {
			str += "\n"
		}
		for i := len(tree.Branches) - 1; i >= 0; i-- {
			str = disp(tree.Branches[i], str, depth+1, i > 0)
		}
		return str
	}

	return disp(tree, "", 0, len(tree.Branches) > 0)
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
func (displayer GraphvizDisplay) Apply(tree *Tree) string {
	var (
		counter int
		disp    func(tree *Tree, str string) string
	)
	disp = func(tree *Tree, str string) string {
		var idx = counter
		str += fmt.Sprintf("\t%d [label=\"%s\"];\n", idx, tree.Operator.String())
		for _, branch := range tree.Branches {
			counter++
			str += fmt.Sprintf("\t%d -> %d;\n", idx, counter)
			str = disp(branch, str)
		}
		return str
	}
	var str = disp(tree, "digraph G {\n")
	str += "}"
	return str
}

// FormulaDisplay outputs an equation-like representation of a Tree.
type FormulaDisplay struct{}

// Apply FormulaDisplay.
func (displayer FormulaDisplay) Apply(tree *Tree) string {
	switch len(tree.Branches) {
	case 0:
		return tree.Operator.String()
	case 1:
		return fmt.Sprintf("%s(%s)", tree.Operator.String(), displayer.Apply(tree.Branches[0]))
	case 2:
		return fmt.Sprintf(
			"(%s)%s(%s)",
			displayer.Apply(tree.Branches[0]),
			tree.Operator.String(),
			displayer.Apply(tree.Branches[1]),
		)
	default:
		return tree.Operator.String()
	}
}
