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
func (displayer DirDisplay) Apply(tr Tree) string {
	var (
		disp       func(tr *Tree, str string, depth int, carriage bool) string
		whitespace = strings.Repeat(" ", displayer.TabSize)
	)

	disp = func(tr *Tree, str string, depth int, carriage bool) string {
		str += strings.Repeat(whitespace, depth) + tr.Op.String()
		if carriage {
			str += "\n"
		}
		for i := len(tr.Branches) - 1; i >= 0; i-- {
			str = disp(tr.Branches[i], str, depth+1, len(tr.Branches) > 0)
		}
		return str
	}

	return disp(&tr, "", 0, len(tr.Branches) > 0)
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
func (displayer GraphvizDisplay) Apply(tr Tree) string {
	var (
		counter int
		disp    func(tr *Tree, str string) string
	)
	disp = func(tr *Tree, str string) string {
		var idx = counter
		str += fmt.Sprintf("  %d [label=\"%s\"];\n", idx, tr.Op.String())
		for _, branch := range tr.Branches {
			counter++
			str += fmt.Sprintf("  %d -> %d;\n", idx, counter)
			str = disp(branch, str)
		}
		return str
	}
	var str = disp(&tr, "digraph G {\n")
	str += "}"
	return str
}

// CodeDisplay outputs an code-like representation of a Tree.
//
// pow(sum(X[0], X[1]), cos(2))
//
type CodeDisplay struct{}

// Apply CodeDisplay.
func (displayer CodeDisplay) Apply(tr Tree) string {
	// Start with the Tree's Operator
	var str = tr.Op.String()

	// Nothing else has to be done if there are no operands
	if len(tr.Branches) == 0 {
		return str
	}

	// Open function
	str += "("

	// Add operands recursively
	for _, br := range tr.Branches {
		str += fmt.Sprintf("%s, ", displayer.Apply(*br))
	}

	// Remove last and close function
	return strings.TrimRight(str, ", ") + ")"
}
