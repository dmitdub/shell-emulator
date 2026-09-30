package tests

import (
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/dmitdub/shell-emulator/src/prompt"
)

// TestBuildFormat checks that the prompt has the expected shape.
//
// The format must be "username@hostname:~$ ".
func TestBuildFormat(t *testing.T) {
	got := prompt.Build()

	if !strings.HasSuffix(got, ":~$ ") {
		t.Errorf("prompt %q does not end with :~$ ", got)
	}
	if !strings.Contains(got, "@") {
		t.Errorf("prompt %q does not contain @", got)
	}
}

// TestBuildContainsHostname checks that the hostname is present.
func TestBuildContainsHostname(t *testing.T) {
	host, err := os.Hostname()
	if err != nil {
		t.Skipf("cannot get hostname: %v", err)
	}

	got := prompt.Build()
	if !strings.Contains(got, host) {
		t.Errorf("prompt %q does not contain hostname %q", got, host)
	}
}

// TestBuildContainsUsername checks that the username is present.
//
// On Windows the OS returns "DOMAIN\User"; the prompt must show
// only the part after the last backslash.
func TestBuildContainsUsername(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Skipf("cannot get user: %v", err)
	}

	name := current.Username
	if idx := strings.LastIndex(name, `\`); idx >= 0 {
		name = name[idx+1:]
	}

	got := prompt.Build()
	if !strings.Contains(got, name) {
		t.Errorf("prompt %q does not contain username %q", got, name)
	}
}

// TestBuildNoDomainPrefix checks that no backslash leaks into the prompt.
func TestBuildNoDomainPrefix(t *testing.T) {
	got := prompt.Build()

	if strings.Contains(got, `\`) {
		t.Errorf("prompt %q contains a domain prefix", got)
	}
}
