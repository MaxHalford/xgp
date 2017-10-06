package tree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPickSubTree(t *testing.T) {
	var (
		tree = &TestTree{
			Value: 0,
			Branches: []*TestTree{
				&TestTree{
					Value: 1,
				},
				&TestTree{
					Value: 2,
					Branches: []*TestTree{
						&TestTree{
							Value: 3,
						},
						&TestTree{
							Value: 4,
						},
					},
				},
			},
		}
		subTree, depth = PickSubTree(
			tree,
			func(tree Tree) float64 { return 1 },
			2,
			2,
			rand.New(rand.NewSource(time.Now().UnixNano())),
		)
	)
	fmt.Println(subTree, depth)
}
