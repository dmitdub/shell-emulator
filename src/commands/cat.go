package commands

import (
	"fmt"
)

// Cat implements the cat command.
//
// It prints the content of a file. Passing a directory or a
// missing path is reported as an error.
func (c *Context) Cat(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(c.ErrOut, "cat: missing operand")
		return
	}

	for _, path := range args {
		node := c.resolve(path)
		if node == nil {
			fmt.Fprintf(c.ErrOut, "cat: %s: No such file or directory\n", path)
			continue
		}
		if node.Dir {
			fmt.Fprintf(c.ErrOut, "cat: %s: Is a directory\n", path)
			continue
		}
		fmt.Fprintln(c.Out, node.Content)
	}
}
