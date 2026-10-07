// Command shell-emulator implements a UNIX-like shell emulator.
//
// Stage 2 accepts a runtime configuration from the command line:
// the path to the virtual file system and the path to a startup
// script. The configuration is printed to stderr at startup. When
// a startup script is provided, its commands are executed before
// the interactive REPL begins.
package main

import (
	"fmt"
	"os"

	"github.com/dmitdub/shell-emulator/src/config"
	"github.com/dmitdub/shell-emulator/src/repl"
)

// main starts the shell emulator.
//
// It reports a fatal error and exits with a non-zero status if
// the configuration is invalid or the REPL cannot be started.
func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "shell-emulator: %v\n", err)
		os.Exit(1)
	}

	runner := repl.New(cfg, os.Stdout, os.Stderr)
	if err := runner.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "shell-emulator: %v\n", err)
		os.Exit(1)
	}
}
