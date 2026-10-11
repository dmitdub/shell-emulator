// Package repl implements the read-eval-print loop of the emulator.
//
// It reads input lines, dispatches them to commands, and prints
// results. In stage 5 the cp command is added to the set of
// commands implemented on top of the virtual file system.
package repl

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/dmitdub/shell-emulator/src/commands"
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
	cfg *config.Config
	ctx *commands.Context
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

	ctx := commands.New(&vfs.Node{Name: "/", Dir: true}, out, errOut)
	return &Runner{cfg: cfg, ctx: ctx}
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
	var root *vfs.Node
	if r.cfg.VFSPath == "" {
		root = vfs.Default()
	} else {
		loaded, err := vfs.Load(r.cfg.VFSPath)
		if err != nil {
			return fmt.Errorf("load VFS: %w", err)
		}
		root = loaded
	}

	out := r.ctx.Out
	errOut := r.ctx.ErrOut
	r.ctx = commands.New(root, out, errOut)
	return nil
}

// loop runs the interactive read-eval-print loop.
func (r *Runner) loop() error {
	rl, err := readline.New(prompt.Build(vfs.Path(r.ctx.Cwd)))
	if err != nil {
		return fmt.Errorf("start readline: %w", err)
	}
	defer rl.Close()

	for {
		rl.SetPrompt(prompt.Build(vfs.Path(r.ctx.Cwd)))

		line, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(r.ctx.Out)
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
	fmt.Fprintln(r.ctx.ErrOut, "configuration:")
	r.cfg.Dump(r.ctx.ErrOut)
	fmt.Fprintf(r.ctx.ErrOut, "vfs.root = %s\n", r.ctx.Root.Name)
	fmt.Fprintf(r.ctx.ErrOut, "vfs.children = %d\n", len(r.ctx.Root.Children))
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

	lines, err := script.Load(r.cfg.ScriptPath)
	if err != nil {
		return false, fmt.Errorf("load startup script: %w", err)
	}

	for _, line := range lines {
		fmt.Fprintf(r.ctx.Out, "%s%s\n", prompt.Build(vfs.Path(r.ctx.Cwd)), line)
		if r.handle(line) {
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
		fmt.Fprintf(r.ctx.ErrOut, "parse error: %v\n", err)
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
		r.ctx.Ls(args)
	case "cd":
		r.ctx.Cd(args)
	case "uptime":
		r.ctx.Uptime(args)
	case "cat":
		r.ctx.Cat(args)
	case "cp":
		r.ctx.Cp(args)
	default:
		fmt.Fprintf(r.ctx.ErrOut, "command not found: %s\n", name)
	}
	return false
}

// dumpVFS prints the current configuration and VFS state.
func (r *Runner) dumpVFS() {
	r.cfg.Dump(r.ctx.Out)
	fmt.Fprintf(r.ctx.Out, "vfs.root = %s\n", r.ctx.Root.Name)
	fmt.Fprintf(r.ctx.Out, "vfs.cwd = %s\n", vfs.Path(r.ctx.Cwd))
	fmt.Fprintf(r.ctx.Out, "vfs.children = %d\n", len(r.ctx.Root.Children))
}

// initVFS replaces the current VFS with the default one.
//
// When a VFS path is configured, the default tree is also written
// back to that file. The original CSV file is overwritten.
func (r *Runner) initVFS() {
	root := vfs.Default()
	out := r.ctx.Out
	errOut := r.ctx.ErrOut
	r.ctx = commands.New(root, out, errOut)

	if r.cfg.VFSPath == "" {
		fmt.Fprintln(r.ctx.Out, "VFS reset to default")
		return
	}

	if err := vfs.Save(r.cfg.VFSPath, root); err != nil {
		fmt.Fprintf(r.ctx.ErrOut, "vfs-init: %v\n", err)
		return
	}
	fmt.Fprintf(r.ctx.Out, "VFS reset to default and saved to %s\n", r.cfg.VFSPath)
}
