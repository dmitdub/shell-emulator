package vfs

// NextID returns the smallest identifier greater than every
// identifier currently used in the tree.
//
// It is used when new nodes are created, for example by the cp
// command, to keep identifiers unique.
func NextID(root *Node) int {
	maxID := -1
	walkAll(root, func(n *Node) {
		if n.ID > maxID {
			maxID = n.ID
		}
	})
	return maxID + 1
}

// walkAll calls fn for every node in the tree in depth-first order.
func walkAll(node *Node, fn func(*Node)) {
	fn(node)
	for _, child := range node.Children {
		walkAll(child, fn)
	}
}
