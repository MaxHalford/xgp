package op

import "fmt"

func ExampleDirDisplay() {
	var (
		op   = Cos{Add{Const{42}, Var{0}}}
		disp = DirDisplay{TabSize: 4}
	)
	fmt.Println(disp.Apply(op))
	// Output:
	// cos
	//     add
	//         42
	//         x0
}

func ExampleGraphvizDisplay() {
	var (
		op   = Cos{Add{Const{42}, Var{0}}}
		disp = GraphvizDisplay{}
	)
	fmt.Println(disp.Apply(op))
	// Output:
	// digraph G {
	//   0 [label="cos"];
	//   0 -> 1;
	//   1 [label="add"];
	//   1 -> 2;
	//   2 [label="42"];
	//   1 -> 3;
	//   3 [label="x0"];
	// }
}

func ExampleCodeDisplay() {
	var (
		op   = Cos{Add{Const{42}, Var{0}}}
		disp = CodeDisplay{}
	)
	fmt.Println(disp.Apply(op))
	// Output: cos(add(42, x0))
}
