// Command shell-emulator implements a UNIX-like shell emulator.
//
// Stage 1 provides a CLI REPL prototype: it reads user input
// interactively, expands environment variables, and dispatches
// commands. The ls and cd commands are stubs that only echo their
// arguments. The exit command terminates the emulator.
package main

import (
	"fmt"
	"os"

	"github.com/dmitdub/shell-emulator/src/repl"
)

// main starts the shell emulator.
//
// It reports a fatal error and exits with a non-zero status if
// the REPL cannot be started or terminates abnormally.
func main() {
	if err := repl.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "shell-emulator: %v\n", err)
		os.Exit(1)
	}
}
