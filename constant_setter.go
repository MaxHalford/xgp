package xgp

// A ConstantSetter can replace a Node's Operator with a Constant.
type ConstantSetter func(value float64)

func (node *Node) newConstantSetter() ConstantSetter {
	return func(value float64) {
		node.Operator = Constant{Value: value}
	}
}
