package tree

type Tree interface {
	NBranches() int
	GetBranch(i int) Tree
	Swap(otherTree Tree)
	ToString() string
}

// rApply recursively applies a function to a Tree and it's branches.
func rApply(tree Tree, f func(tree Tree, depth int) (stop bool)) {

	var apply func(tree Tree, depth int) bool
	apply = func(tree Tree, depth int) bool {
		// Apply the function to the Tree and check if the recursion should stop
		var stop = f(tree, depth)
		if stop {
			return stop
		}
		// Apply recursion to each branch
		for i := 0; i < tree.NBranches(); i++ {
			stop = apply(tree.GetBranch(i), depth+1)
			if stop {
				break
			}
		}
		return false
	}

	apply(tree, 0)
}

// GetHeight returns the height of a tree. The height of a tree is the height of
// its root node. The height of a node is the number of edges on the longest
// path between that node and a leaf.
func GetHeight(tree Tree) int {
	var maxHeight = -1
	for i := 0; i < tree.NBranches(); i++ {
		var height = GetHeight(tree.GetBranch(i))
		if height > maxHeight {
			maxHeight = height
		}
	}
	return maxHeight + 1
}

// GetNNodes returns the number of nodes in a Tree.
func GetNNodes(tree Tree) (n int) {
	rApply(tree, func(Tree, int) bool { n++; return false })
	return
}
