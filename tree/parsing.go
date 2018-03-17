package tree

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaxHalford/xgp/op"
)

// ParseCode takes a code representation of a Tree and parses it into a Tree.
func ParseCode(code string) (Tree, error) {
	var tree Tree
	// The operator is either a Variable or either a Constant
	if !strings.HasSuffix(code, ")") {

		// The operator is a Variable
		if strings.HasSuffix(code, "]") {
			var index, err = strconv.Atoi(code[2 : len(code)-1])
			if err != nil {
				return tree, err
			}
			tree.op = op.Variable{Index: index}
			return tree, nil
		}

		// The operator is a Constant
		var value, err = strconv.ParseFloat(code, 64)
		if err != nil {
			return tree, err
		}
		tree.op = op.Constant{Value: value}
		return tree, nil
	}

	// The operator is a function
	var (
		parts         = regexp.MustCompile("\\(").Split(code, 2)
		operator, err = op.ParseFuncName(parts[0])
		inside        = parts[1][:len(parts[1])-1]
	)
	if err != nil {
		return tree, err
	}

	var (
		operands           = []string{}
		operand            = ""
		parenthesesCounter int
	)
	for _, r := range inside {
		var s = string(r)
		if s == " " {
			continue
		}
		if s == "(" {
			parenthesesCounter++
		}
		if s == "," && parenthesesCounter <= 0 {
			operands = append(operands, operand)
			operand = ""
		} else {
			operand += s
		}
		if s == ")" {
			parenthesesCounter--
		}
	}
	operands = append(operands, operand)

	// Check the number of operands if consistent with the arity of the operator
	if len(operands) != operator.Arity() {
		return tree, errors.New("Number of operands does not match with operator arity")
	}

	tree.op = operator
	tree.branches = make([]*Tree, len(operands))
	for i, operand := range operands {
		// Parse the operand
		var branch, err = ParseCode(operand)
		if err != nil {
			return tree, err
		}
		tree.SetBranch(i, branch)
	}

	return tree, nil
}

// MustParseCode is identical to ParseCode but doesn't return an error. This
// method should only be used for testing purposes.
func MustParseCode(code string) Tree {
	var tree, _ = ParseCode(code)
	return tree
}
