package tree

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func ParseCode(code string) (*Tree, error) {
	var tree Tree

	// The operator is either a Variable or either a Constant
	if !strings.HasSuffix(code, ")") {

		// The operator is a Variable
		if strings.HasSuffix(code, "]") {
			var index, err = strconv.Atoi(code[2 : len(code)-1])
			if err != nil {
				return nil, err
			}
			tree.Operator = Variable{Index: index}
			return &tree, nil
		}

		// The operator is a Constant
		var value, err = strconv.ParseFloat(code, 64)
		if err != nil {
			return nil, err
		}
		tree.Operator = Constant{Value: value}
		return &tree, nil
	}

	// The operator is a function
	var (
		parts         = regexp.MustCompile("\\(").Split(code, 2)
		operator, err = parseFuncName(parts[0])
	)
	if err != nil {
		return nil, err
	}

	// Remove the trailing closing parenthesis
	parts[len(parts)-1] = parts[len(parts)-1][:len(parts[len(parts)-1])-1]

	var (
		operands         = []string{}
		operand          = ""
		parenthisCounter int
	)
	for _, r := range parts[1] {
		var s = string(r)
		if s == " " {
			continue
		}
		if s == "(" {
			parenthisCounter++
		}
		if s == "," && parenthisCounter <= 0 {
			operands = append(operands, operand)
			operand = ""
		} else {
			operand += s
		}
		if s == ")" {
			parenthisCounter--
		}
	}
	operands = append(operands, operand)

	// Check the number of operands if consistent with the arity of the operator
	if len(operands) != operator.Arity() {
		return nil, errors.New("Number of operands does not match with operator arity")
	}

	tree.Operator = operator
	tree.Branches = make([]*Tree, len(operands))
	for i, operand := range operands {
		// Parse the operand
		tree.Branches[i], err = ParseCode(operand)
		if err != nil {
			return nil, err
		}
	}

	return &tree, nil
}
