package commands

import (
	"fmt"
	"strings"

	"github.com/dmitdub/shell-emulator/src/vfs"
)

// Cd implements the cd command.
//
// It changes the current working directory. Without arguments
// the root directory is selected. The special path ".." moves
// to the parent directory. Passing a file or a missing path is
// reported as an error.
func (c *Context) Cd(args []string) {
	if len(args) == 0 {
		c.Cwd = c.Root
		return
	}

	target := c.resolve(args[0])
	if target == nil {
		fmt.Fprintf(c.ErrOut, "cd: %s: No such file or directory\n", args[0])
		return
	}
	if !target.Dir {
		fmt.Fprintf(c.ErrOut, "cd: %s: Not a directory\n", args[0])
		return
	}
	c.Cwd = target
}

// resolve returns the node addressed by path, or nil if it does
// not exist. Paths may be absolute (starting with a slash) or
// relative to the current working directory.
func (c *Context) resolve(path string) *vfs.Node {
	if path == "" {
		return nil
	}
	if path == "." {
		return c.Cwd
	}
	if path == ".." {
		return c.parent(c.Cwd)
	}

	if strings.HasPrefix(path, "/") {
		return c.walk(c.Root, splitPath(path))
	}
	return c.walk(c.Cwd, splitPath(path))
}

// walk follows the given path components starting from node.
//
// It returns nil when any component is missing or when an
// intermediate component is not a directory.
func (c *Context) walk(node *vfs.Node, parts []string) *vfs.Node {
	current := node
	for _, part := range parts {
		if !current.Dir {
			return nil
		}
		next := current.FindChild(part)
		if next == nil {
			return nil
		}
		current = next
	}
	return current
}

// parent returns the parent of node, or node itself for the root.
func (c *Context) parent(node *vfs.Node) *vfs.Node {
	if node.Parent == nil {
		return node
	}
	return node.Parent
}

// splitPath splits path into non-empty components.
func splitPath(path string) []string {
	raw := strings.Split(path, "/")
	parts := make([]string, 0, len(raw))
	for _, p := range raw {
		if p != "" && p != "." {
			parts = append(parts, p)
		}
	}
	return parts
}
