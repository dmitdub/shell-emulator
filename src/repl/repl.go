// Package repl implements the read-eval-print loop of the emulator.
//
// It reads input lines, dispatches them to commands, and prints
// results. In stage 1 the ls and cd commands are stubs that only
// echo their arguments, and the exit command terminates the loop.
package repl

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/dmitdub/shell-emulator/src/parser"
	"github.com/dmitdub/shell-emulator/src/prompt"
)

// exitCommand terminates the emulator.
const exitCommand = "exit"

// Run starts the read-eval-print loop.
//
// Run returns nil after the user types exit or presses Ctrl+D.
// It returns an error if the readline instance cannot be created.
func Run() error {
	rl, err := readline.New(prompt.Build())
	if err != nil {
		return fmt.Errorf("start readline: %w", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println()
			return nil
		}
		if handleLine(line) {
			return nil
		}
	}
}

// handleLine parses and dispatches a single input line.
//
// It returns true when the REPL must terminate.
func handleLine(line string) bool {
	args, err := parser.Parse(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		return false
	}
	if len(args) == 0 {
		return false
	}

	return dispatch(args[0], args[1:])
}

// dispatch routes a command name and its arguments to a handler.
//
// It returns true when the command requests termination.
func dispatch(name string, args []string) bool {
	switch name {
	case exitCommand:
		return true
	case "ls":
		runStub("ls", args)
	case "cd":
		runStub("cd", args)
	default:
		fmt.Fprintf(os.Stderr, "command not found: %s\n", name)
	}
	return false
}

// runStub prints a command name along with its arguments.
//
// It temporarily stands in for a real command implementation.
func runStub(name string, args []string) {
	fmt.Println(name)
	for _, arg := range args {
		fmt.Printf("  %s\n", arg)
	}
}
