package op

import (
	"fmt"
	"strings"
)

// A Displayer outputs a string representation of a Operator.
type Displayer interface {
	Apply(Operator) string
}

// DirDisplay outputs a directory like representation of an Operator.
type DirDisplay struct {
	TabSize int
}

// Apply DirDisplay.
func (displayer DirDisplay) Apply(op Operator) string {
	var (
		disp       func(op Operator, str string, depth int) string
		whitespace = strings.Repeat(" ", displayer.TabSize)
	)

	disp = func(op Operator, str string, depth int) string {
		if op == nil {
			return str
		}
		str += strings.Repeat(whitespace, depth) + op.Name() + "\n"
		for i := uint(0); i < op.Arity(); i++ {
			str = disp(op.Operand(i), str, depth+1)
		}
		return str
	}

	return strings.TrimRight(disp(op, "", 0), "\n")
}

// GraphvizDisplay outputs a Graphviz representation of an Operator. Each
// branch is indexed with a global counter and is labelled with the Operator's
// Name method.
type GraphvizDisplay struct{}

// Apply Graphviz display.
func (displayer GraphvizDisplay) Apply(op Operator) string {
	var (
		counter int
		disp    func(op Operator, str string) string
	)
	disp = func(op Operator, str string) string {
		var c = counter
		str += fmt.Sprintf("  %d [label=\"%s\"];\n", c, op.Name())
		for i := uint(0); i < op.Arity(); i++ {
			counter++
			str += fmt.Sprintf("  %d -> %d;\n", c, counter)
			str = disp(op.Operand(i), str)
		}
		return str
	}
	var str = disp(op, "digraph G {\n")
	str += "}"
	return str
}

// CodeDisplay outputs an code-like representation of an Operator.
type CodeDisplay struct{}

// Apply CodeDisplay.
func (displayer CodeDisplay) Apply(op Operator) string {
	// Start with the Operator's name
	var str = op.Name()

	// Nothing else has to be done if there are no operands
	if op.Arity() == 0 {
		return str
	}

	// Open function
	str += "("

	// Add operands recursively
	for i := uint(0); i < op.Arity(); i++ {
		str += fmt.Sprintf("%s, ", displayer.Apply(op.Operand(i)))
	}

	// Remove last and close function
	return strings.TrimRight(str, ", ") + ")"
}

func parenthesize(op Operator) string {
	if op.Arity() > 1 {
		return fmt.Sprintf("(%s)", op.String())
	}
	return op.String()
}
