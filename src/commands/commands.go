// Package commands implements the built-in commands of the emulator.
//
// Every command operates on the virtual file system through a
// shared Context. The Context holds the current working directory
// and the writers used for output.
package commands

import (
	"io"
	"time"

	"github.com/dmitdub/shell-emulator/src/vfs"
)

// Context carries the shared state of a single emulator session.
type Context struct {
	// Root is the root of the virtual file system.
	Root *vfs.Node

	// Cwd is the current working directory.
	Cwd *vfs.Node

	// Started is the moment the emulator was started.
	Started time.Time

	// Out receives normal command output.
	Out io.Writer

	// ErrOut receives error messages.
	ErrOut io.Writer
}

// New creates a Context for the given VFS root.
//
// The current working directory is initially set to the root,
// and the start time is captured at the moment of the call.
func New(root *vfs.Node, out, errOut io.Writer) *Context {
	return &Context{
		Root:    root,
		Cwd:     root,
		Started: time.Now(),
		Out:     out,
		ErrOut:  errOut,
	}
}
