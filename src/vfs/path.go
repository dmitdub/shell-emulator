package vfs

import "strings"

// Path returns the absolute path of the given node.
//
// The path is built by walking up the parent chain until the
// root is reached. The root itself renders as a single slash.
func Path(node *Node) string {
	if node.Parent == nil {
		return "/"
	}

	var parts []string
	current := node
	for current.Parent != nil {
		parts = append([]string{current.Name}, parts...)
		current = current.Parent
	}
	return "/" + strings.Join(parts, "/")
}
