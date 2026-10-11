package commands

import (
	"fmt"
	"strings"

	"github.com/dmitdub/shell-emulator/src/vfs"
)

// cpFlags holds the parsed flags of the cp command.
type cpFlags struct {
	recursive bool
}

// Cp implements the cp command.
//
// It copies a file or a directory inside the virtual file system.
// Directories require the -r (or -R) flag. Copying a directory
// into itself or into its own descendant is rejected. Existing
// targets are overwritten.
func (c *Context) Cp(args []string) {
	flags, paths := parseCpArgs(args)
	if len(paths) != 2 {
		fmt.Fprintln(c.ErrOut, "cp: usage: cp [-r] source target")
		return
	}

	src := c.resolve(paths[0])
	if src == nil {
		fmt.Fprintf(c.ErrOut, "cp: %s: No such file or directory\n", paths[0])
		return
	}
	if src.Dir && !flags.recursive {
		fmt.Fprintf(c.ErrOut, "cp: -r not specified; omitting directory %q\n", paths[0])
		return
	}

	dst := c.resolve(paths[1])
	if dst == nil {
		c.copyToNew(src, paths[1], flags.recursive)
		return
	}
	c.copyInto(src, dst, flags.recursive)
}

// copyToNew creates a new node with the target path.
func (c *Context) copyToNew(src *vfs.Node, path string, recursive bool) {
	if src.Dir && !recursive {
		return
	}

	parent, name := c.splitTarget(path)
	if parent == nil {
		fmt.Fprintf(c.ErrOut, "cp: %s: No such file or directory\n", path)
		return
	}

	if src.Dir && isDescendant(parent, src) {
		fmt.Fprintln(c.ErrOut, "cp: cannot copy a directory into itself")
		return
	}

	next := vfs.NextID(c.Root)
	clone := cloneNode(src, &next)
	clone.Name = name

	parent.AddChild(clone)
}

// copyInto copies src into the existing target node dst.
func (c *Context) copyInto(src, dst *vfs.Node, recursive bool) {
	if !dst.Dir {
		c.overwrite(src, dst)
		return
	}

	if src.Dir && isDescendant(dst, src) {
		fmt.Fprintln(c.ErrOut, "cp: cannot copy a directory into itself")
		return
	}

	if existing := dst.FindChild(src.Name); existing != nil {
		c.overwrite(src, existing)
		return
	}

	next := vfs.NextID(c.Root)
	clone := cloneNode(src, &next)
	clone.Name = src.Name

	dst.AddChild(clone)
}

// overwrite replaces the content of dst with a copy of src.
func (c *Context) overwrite(src, dst *vfs.Node) {
	if src.Dir || dst.Dir {
		fmt.Fprintln(c.ErrOut, "cp: cannot overwrite directory")
		return
	}
	dst.Content = src.Content
}

// splitTarget returns the parent directory and the final name
// for a target path that does not exist yet.
func (c *Context) splitTarget(path string) (*vfs.Node, string) {
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return c.Cwd, path
	}

	dirPath := path[:idx]
	name := path[idx+1:]
	if dirPath == "" {
		dirPath = "/"
	}

	parent := c.resolve(dirPath)
	if parent == nil || !parent.Dir {
		return nil, ""
	}
	return parent, name
}

// cloneNode returns a deep copy of node.
//
// New nodes receive identifiers starting from next, which is
// incremented after every allocation.
func cloneNode(node *vfs.Node, next *int) *vfs.Node {
	clone := &vfs.Node{
		ID:      *next,
		Name:    node.Name,
		Dir:     node.Dir,
		Content: node.Content,
	}
	*next = *next + 1

	for _, child := range node.Children {
		childClone := cloneNode(child, next)
		clone.AddChild(childClone)
	}
	return clone
}

// isDescendant reports whether node is a strict descendant of ancestor.
func isDescendant(node, ancestor *vfs.Node) bool {
	current := node.Parent
	for current != nil {
		if current == ancestor {
			return true
		}
		current = current.Parent
	}
	return false
}

// parseCpArgs splits args into flags and paths.
//
// Supported flags are -r and -R, both meaning recursive copy.
func parseCpArgs(args []string) (cpFlags, []string) {
	var flags cpFlags
	var paths []string

	for _, arg := range args {
		if isFlag(arg) {
			applyCpFlag(&flags, arg)
			continue
		}
		paths = append(paths, arg)
	}
	return flags, paths
}

// applyCpFlag sets the corresponding flags for a combined flag.
func applyCpFlag(flags *cpFlags, arg string) {
	for _, ch := range arg[1:] {
		switch ch {
		case 'r', 'R':
			flags.recursive = true
		}
	}
}
