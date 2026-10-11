// Package vfs implements an in-memory virtual file system.
//
// The file system is loaded from a CSV file and kept in memory.
// The original CSV file is never modified, except by the vfs-init
// command, which rewrites it with default content.
package vfs

import "sort"

// Node represents a file or directory in the virtual file system.
type Node struct {
	// ID is the unique identifier of the node.
	ID int

	// Name is the name of the file or directory.
	Name string

	// Dir reports whether the node is a directory.
	Dir bool

	// Content holds the file content. It is empty for directories.
	Content string

	// Parent points to the parent node. It is nil for the root.
	Parent *Node

	// Children holds the child nodes. It is empty for files.
	Children []*Node
}

// AddChild appends a child node, sets its parent, and keeps
// children sorted by name.
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Name < n.Children[j].Name
	})
}

// FindChild returns the child with the given name, or nil.
func (n *Node) FindChild(name string) *Node {
	for _, child := range n.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}
