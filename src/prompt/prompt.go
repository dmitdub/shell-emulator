// Package prompt builds the interactive prompt of the emulator.
//
// The prompt is derived from real operating system data and the
// current working directory of the virtual file system. It looks
// like "username@hostname:/path$ ".
package prompt

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

// Default values are used when the operating system does not
// provide the corresponding information.
const (
	defaultUser = "user"
	defaultHost = "localhost"
)

// Build returns the prompt string for the emulator.
//
// The cwd parameter is the absolute path of the current working
// directory in the virtual file system. On Windows the username
// returned by the OS includes a domain prefix, which is stripped
// to keep the prompt readable.
func Build(cwd string) string {
	return fmt.Sprintf("%s@%s:%s$ ", username(), hostname(), cwd)
}

// username returns the current user name without a domain prefix.
func username() string {
	current, err := user.Current()
	if err != nil {
		return defaultUser
	}

	name := current.Username
	if idx := strings.LastIndex(name, `\`); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

// hostname returns the name of the current host.
func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return defaultHost
	}
	return name
}
