package op

import "fmt"

// ParseFuncName returns a functional Operator from it's String representation.
func ParseFuncName(funcName string) (Operator, error) {
	var f, ok = map[string]Operator{
		Cos{}.String():        Cos{},
		Sin{}.String():        Sin{},
		Log{}.String():        Log{},
		Exp{}.String():        Exp{},
		Max{}.String():        Max{},
		Min{}.String():        Min{},
		Sum{}.String():        Sum{},
		Difference{}.String(): Difference{},
		Division{}.String():   Division{},
		Product{}.String():    Product{},
		Power{}.String():      Power{},
	}[funcName]
	if !ok {
		return nil, fmt.Errorf("Unknown function name '%s'", funcName)
	}
	return f, nil
}
