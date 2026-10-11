package commands

import (
	"fmt"
	"strings"

	"github.com/dmitdub/shell-emulator/src/vfs"
)

// lsFlags holds the parsed flags of the ls command.
type lsFlags struct {
	long  bool
	all   bool
	human bool
}

// Ls implements the ls command.
//
// It prints the contents of a directory. When no path is given,
// the current working directory is listed. Supported flags are
// -l (long format), -a (include hidden files), and -h (human
// readable sizes, only meaningful with -l).
func (c *Context) Ls(args []string) {
	flags, paths := parseLsArgs(args)

	if len(paths) == 0 {
		c.printDir(c.Cwd, flags)
		return
	}
	for _, path := range paths {
		node := c.resolve(path)
		if node == nil {
			fmt.Fprintf(c.ErrOut, "ls: %s: No such file or directory\n", path)
			continue
		}
		c.printNode(node, flags)
	}
}

// printNode prints a single node.
//
// Directories are listed by their content; files are printed
// by their name (or full row in long mode).
func (c *Context) printNode(node *vfs.Node, flags lsFlags) {
	if node.Dir {
		c.printDir(node, flags)
		return
	}
	if flags.long {
		fmt.Fprintln(c.Out, formatLong(node, flags.human))
		return
	}
	fmt.Fprintln(c.Out, node.Name)
}

// printDir prints the contents of a directory.
//
// In long mode each entry is printed on its own line.
// In short mode entries are separated by two spaces.
func (c *Context) printDir(dir *vfs.Node, flags lsFlags) {
	entries := visibleChildren(dir, flags.all)

	if flags.long {
		for _, child := range entries {
			fmt.Fprintln(c.Out, formatLong(child, flags.human))
		}
		return
	}

	names := make([]string, 0, len(entries))
	for _, child := range entries {
		names = append(names, formatName(child))
	}
	fmt.Fprintln(c.Out, strings.Join(names, "  "))
}

// visibleChildren returns the children of dir.
//
// When all is false, entries whose name starts with a dot are
// filtered out.
func visibleChildren(dir *vfs.Node, all bool) []*vfs.Node {
	if all {
		return dir.Children
	}
	visible := make([]*vfs.Node, 0, len(dir.Children))
	for _, child := range dir.Children {
		if !strings.HasPrefix(child.Name, ".") {
			visible = append(visible, child)
		}
	}
	return visible
}

// formatName returns the display name of a node.
//
// Directories are suffixed with a slash.
func formatName(node *vfs.Node) string {
	if node.Dir {
		return node.Name + "/"
	}
	return node.Name
}

// formatLong returns a single long-format row for a node.
func formatLong(node *vfs.Node, human bool) string {
	perms := "-rw-r--r--"
	size := fmt.Sprintf("%d", len(node.Content))
	if node.Dir {
		perms = "drwxr-xr-x"
		size = "-"
	} else if human {
		size = humanSize(len(node.Content))
	}
	return fmt.Sprintf("%s  %s  %s", perms, size, node.Name)
}

// humanSize formats a byte count as a human-readable string.
//
// Values below one kilobyte are printed as-is.
func humanSize(n int) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case n >= gb:
		return fmt.Sprintf("%.1fG", float64(n)/gb)
	case n >= mb:
		return fmt.Sprintf("%.1fM", float64(n)/mb)
	case n >= kb:
		return fmt.Sprintf("%.1fK", float64(n)/kb)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// parseLsArgs splits args into flags and paths.
//
// Flags may be combined, as in "-lah". Unknown flags are ignored.
func parseLsArgs(args []string) (lsFlags, []string) {
	var flags lsFlags
	var paths []string

	for _, arg := range args {
		if isFlag(arg) {
			applyLsFlag(&flags, arg)
			continue
		}
		paths = append(paths, arg)
	}
	return flags, paths
}

// isFlag reports whether arg looks like a flag.
func isFlag(arg string) bool {
	return len(arg) > 1 && arg[0] == '-'
}

// applyLsFlag sets the corresponding flags for a combined flag.
func applyLsFlag(flags *lsFlags, arg string) {
	for _, ch := range arg[1:] {
		switch ch {
		case 'l':
			flags.long = true
		case 'a':
			flags.all = true
		case 'h':
			flags.human = true
		}
	}
}
