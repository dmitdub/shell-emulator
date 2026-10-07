// Package repl implements the read-eval-print loop of the emulator.
//
// It reads input lines, dispatches them to commands, and prints
// results. In stage 2 the emulator accepts a runtime configuration,
// executes an optional startup script, and provides the conf-dump
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
)

// exitCommand terminates the emulator.
const exitCommand = "exit"

// confDumpCommand prints the current configuration.
const confDumpCommand = "conf-dump"

// Runner holds the state of a single emulator session.
type Runner struct {
	cfg    *config.Config
	out    io.Writer
	errOut io.Writer
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

// Run executes the startup script (if any) and then starts the REPL.
//
// Run returns nil after the user types exit or presses Ctrl+D.
// It does not start the interactive loop when the startup script
// terminates the session.
func (r *Runner) Run() error {
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
		r.cfg.Dump(r.out)
	case "ls":
		r.stub("ls", args)
	case "cd":
		r.stub("cd", args)
	default:
		fmt.Fprintf(r.errOut, "command not found: %s\n", name)
	}
	return false
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
