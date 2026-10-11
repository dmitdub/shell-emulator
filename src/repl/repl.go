// Package repl implements the read-eval-print loop of the emulator.
//
// It reads input lines, dispatches them to commands, and prints
// results. In stage 3 the emulator loads a virtual file system
// from a CSV file, keeps it in memory, and provides the vfs-init
// command. The ls and cd commands remain stubs.
package repl

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/dmitdub/shell-emulator/src/config"
	"github.com/dmitdub/shell-emulator/src/parser"
	"github.com/dmitdub/shell-emulator/src/prompt"
	"github.com/dmitdub/shell-emulator/src/script"
	"github.com/dmitdub/shell-emulator/src/vfs"
)

// Command names recognized by the emulator.
const (
	exitCommand     = "exit"
	confDumpCommand = "conf-dump"
	vfsInitCommand  = "vfs-init"
)

// Runner holds the state of a single emulator session.
type Runner struct {
	cfg    *config.Config
	out    io.Writer
	errOut io.Writer
	root   *vfs.Node
}

// New creates a Runner with the given configuration.
//
// The out writer receives command output; errOut receives errors.
// When out or errOut is nil, os.Stdout and os.Stderr are used.
func New(cfg *config.Config, out, errOut io.Writer) *Runner {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return &Runner{cfg: cfg, out: out, errOut: errOut}
}

// Run loads the VFS, executes the startup script, and starts the REPL.
//
// Run returns nil after the user types exit or presses Ctrl+D.
// It does not start the interactive loop when the startup script
// terminates the session.
func (r *Runner) Run() error {
	if err := r.loadVFS(); err != nil {
		return err
	}
	r.dumpConfig()

	done, err := r.runStartupScript()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return r.loop()
}

// loadVFS loads the virtual file system into memory.
//
// When the VFS path is empty, a default in-memory tree is used.
func (r *Runner) loadVFS() error {
	if r.cfg.VFSPath == "" {
		r.root = vfs.Default()
		return nil
	}

	root, err := vfs.Load(r.cfg.VFSPath)
	if err != nil {
		return fmt.Errorf("load VFS: %w", err)
	}
	r.root = root
	return nil
}

// loop runs the interactive read-eval-print loop.
func (r *Runner) loop() error {
	rl, err := readline.New(prompt.Build())
	if err != nil {
		return fmt.Errorf("start readline: %w", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(r.out)
			return nil
		}
		if r.handle(line) {
			return nil
		}
	}
}

// dumpConfig writes all configuration parameters to the error output.
//
// The dump is produced once at startup for debugging purposes.
func (r *Runner) dumpConfig() {
	fmt.Fprintln(r.errOut, "configuration:")
	r.cfg.Dump(r.errOut)
	fmt.Fprintf(r.errOut, "vfs.root = %s\n", r.root.Name)
	fmt.Fprintf(r.errOut, "vfs.children = %d\n", len(r.root.Children))
}

// runStartupScript executes commands from the startup script, if any.
//
// Every command is printed with the prompt before it is executed,
// imitating an interactive session. It returns true when a command
// requests termination, so the caller can skip the interactive loop.
func (r *Runner) runStartupScript() (bool, error) {
	if r.cfg.ScriptPath == "" {
		return false, nil
	}

	commands, err := script.Load(r.cfg.ScriptPath)
	if err != nil {
		return false, fmt.Errorf("load startup script: %w", err)
	}

	p := prompt.Build()
	for _, cmd := range commands {
		fmt.Fprintf(r.out, "%s%s\n", p, cmd)
		if r.handle(cmd) {
			return true, nil
		}
	}
	return false, nil
}

// handle parses and dispatches a single input line.
//
// It returns true when the REPL must terminate.
func (r *Runner) handle(line string) bool {
	args, err := parser.Parse(line)
	if err != nil {
		fmt.Fprintf(r.errOut, "parse error: %v\n", err)
		return false
	}
	if len(args) == 0 {
		return false
	}
	return r.dispatch(args[0], args[1:])
}

// dispatch routes a command name and its arguments to a handler.
//
// It returns true when the command requests termination.
func (r *Runner) dispatch(name string, args []string) bool {
	switch name {
	case exitCommand:
		return true
	case confDumpCommand:
		r.dumpVFS()
	case vfsInitCommand:
		r.initVFS()
	case "ls":
		r.stub("ls", args)
	case "cd":
		r.stub("cd", args)
	default:
		fmt.Fprintf(r.errOut, "command not found: %s\n", name)
	}
	return false
}

// dumpVFS prints the current configuration and VFS state.
func (r *Runner) dumpVFS() {
	r.cfg.Dump(r.out)
	fmt.Fprintf(r.out, "vfs.root = %s\n", r.root.Name)
	fmt.Fprintf(r.out, "vfs.children = %d\n", len(r.root.Children))
}

// initVFS replaces the current VFS with the default one.
//
// When a VFS path is configured, the default tree is also written
// back to that file. The original CSV file is overwritten.
func (r *Runner) initVFS() {
	r.root = vfs.Default()

	if r.cfg.VFSPath == "" {
		fmt.Fprintln(r.out, "VFS reset to default")
		return
	}

	if err := vfs.Save(r.cfg.VFSPath, r.root); err != nil {
		fmt.Fprintf(r.errOut, "vfs-init: %v\n", err)
		return
	}
	fmt.Fprintf(r.out, "VFS reset to default and saved to %s\n", r.cfg.VFSPath)
}

// stub prints a command name along with its arguments.
//
// It temporarily stands in for a real command implementation.
func (r *Runner) stub(name string, args []string) {
	fmt.Fprintln(r.out, name)
	for _, arg := range args {
		fmt.Fprintf(r.out, "  %s\n", arg)
	}
}
