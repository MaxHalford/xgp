package tree

import (
	"fmt"

	"github.com/MaxHalford/xgp/op"
)

func ExampleDirDisplay() {
	var (
		tree = Tree{
			Op: op.Cos{},
			Branches: []*Tree{
				&Tree{
					Op: op.Sum{},
					Branches: []*Tree{
						&Tree{Op: op.Variable{0}},
						&Tree{Op: op.Constant{42}},
					},
				},
			},
		}
		disp = DirDisplay{TabSize: 4}
	)
	fmt.Println(disp.Apply(tree))
	// Output:
	// cos
	//     sum
	//         42
	//         X[0]
}

func ExampleGraphvizDisplay() {
	var (
		tree = Tree{
			Op: op.Cos{},
			Branches: []*Tree{
				&Tree{
					Op: op.Sum{},
					Branches: []*Tree{
						&Tree{Op: op.Variable{0}},
						&Tree{Op: op.Constant{42}},
					},
				},
			},
		}
		disp = GraphvizDisplay{}
	)
	fmt.Println(disp.Apply(tree))
	// Output:
	// digraph G {
	//   0 [label="cos"];
	//   0 -> 1;
	//   1 [label="sum"];
	//   1 -> 2;
	//   2 [label="X[0]"];
	//   1 -> 3;
	//   3 [label="42"];
	// }
}

func ExampleCodeDisplay() {
	var (
		tree = Tree{
			Op: op.Cos{},
			Branches: []*Tree{
				&Tree{
					Op: op.Sum{},
					Branches: []*Tree{
						&Tree{Op: op.Variable{0}},
						&Tree{Op: op.Constant{42}},
					},
				},
			},
		}
		disp = CodeDisplay{}
	)
	fmt.Println(disp.Apply(tree))
	// Output: cos(sum(X[0], 42))
}
