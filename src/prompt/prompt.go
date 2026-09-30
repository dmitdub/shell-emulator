// Package prompt builds the interactive prompt of the emulator.
//
// The prompt is derived from real operating system data and looks
// like "username@hostname:~$ ".
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
// On Windows the username returned by the OS includes a domain
// prefix, which is stripped to keep the prompt readable.
func Build() string {
	return fmt.Sprintf("%s@%s:~$ ", username(), hostname())
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
