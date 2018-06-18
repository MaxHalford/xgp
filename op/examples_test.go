package op

import "fmt"

// From https://stats.stackexchange.com/questions/224140/step-by-step-example-of-reverse-mode-automatic-differentiation.
func Example() {
	var z = Add{Mul{Var{0}, Var{1}}, Sin{Var{0}}}
	fmt.Println("z:", z)
	fmt.Println("∂z/∂x0:", z.Diff(0).Simplify())
	fmt.Println("∂z/∂x1:", z.Diff(1).Simplify())
	// Output:
	// z: x0*x1+sin(x0)
	// ∂z/∂x0: x1+cos(x0)
	// ∂z/∂x1: x0
}
