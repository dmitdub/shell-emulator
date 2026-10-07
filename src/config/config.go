// Package config parses the command line configuration of the emulator.
//
// The emulator accepts two optional flags:
//
//	--vfs    path to the virtual file system
//	--script path to the startup script
//
// Both flags are optional. When a flag is not provided the
// corresponding field of Config is left empty.
package config

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Config holds the runtime configuration of the emulator.
type Config struct {
	// VFSPath is the path to the virtual file system.
	// It is empty when the --vfs flag is not provided.
	VFSPath string

	// ScriptPath is the path to the startup script.
	// It is empty when the --script flag is not provided.
	ScriptPath string
}

// Parse reads the given command line arguments and returns a Config.
//
// Parse returns an error when the startup script is specified but
// does not exist. The VFS path is not validated here; the virtual
// file system is loaded later.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("shell-emulator", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	vfsPath := fs.String("vfs", "", "path to the virtual file system")
	scriptPath := fs.String("script", "", "path to the startup script")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		VFSPath:    *vfsPath,
		ScriptPath: *scriptPath,
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// validate checks that the configuration is internally consistent.
func (c *Config) validate() error {
	if c.ScriptPath == "" {
		return nil
	}
	if _, err := os.Stat(c.ScriptPath); err != nil {
		return fmt.Errorf("startup script: %w", err)
	}
	return nil
}

// Dump writes all configuration parameters to the given writer.
//
// The output uses the "key = value" format, one parameter per line.
func (c *Config) Dump(w io.Writer) {
	fmt.Fprintf(w, "vfs.path = %s\n", c.VFSPath)
	fmt.Fprintf(w, "script.path = %s\n", c.ScriptPath)
}
