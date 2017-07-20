package tree

type Tree interface {
	NBranches() int
	GetBranch(i int) Tree
	Swap(otherTree Tree)
	ToString() string
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

// NNodes returns the number of nodes in a Tree.
func GetNNodes(tree Tree) int {
	var n = 1
	for i := 0; i < tree.NBranches(); i++ {
		n += GetNNodes(tree.GetBranch(i))
	}
	return n
}
