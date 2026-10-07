package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmitdub/shell-emulator/src/config"
)

// TestParseNoFlags checks that empty arguments yield an empty config.
func TestParseNoFlags(t *testing.T) {
	cfg, err := config.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VFSPath != "" {
		t.Errorf("VFSPath = %q, want empty", cfg.VFSPath)
	}
	if cfg.ScriptPath != "" {
		t.Errorf("ScriptPath = %q, want empty", cfg.ScriptPath)
	}
}

// TestParseVFSOnly checks that the --vfs flag is parsed.
func TestParseVFSOnly(t *testing.T) {
	cfg, err := config.Parse([]string{"--vfs", "/tmp/data.csv"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VFSPath != "/tmp/data.csv" {
		t.Errorf("VFSPath = %q, want /tmp/data.csv", cfg.VFSPath)
	}
	if cfg.ScriptPath != "" {
		t.Errorf("ScriptPath = %q, want empty", cfg.ScriptPath)
	}
}

// TestParseScriptOnly checks that an existing --script file is accepted.
func TestParseScriptOnly(t *testing.T) {
	path := writeTempFile(t, "ls\ncd /tmp\n")

	cfg, err := config.Parse([]string{"--script", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScriptPath != path {
		t.Errorf("ScriptPath = %q, want %q", cfg.ScriptPath, path)
	}
}

// TestParseMissingScript checks that a missing script file is rejected.
func TestParseMissingScript(t *testing.T) {
	_, err := config.Parse([]string{"--script", "/no/such/file.txt"})
	if err == nil {
		t.Fatal("expected error for missing script, got nil")
	}
}

// TestParseBothFlags checks that both flags are parsed together.
func TestParseBothFlags(t *testing.T) {
	path := writeTempFile(t, "ls\n")

	cfg, err := config.Parse([]string{"--vfs", "/tmp/vfs.csv", "--script", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VFSPath != "/tmp/vfs.csv" {
		t.Errorf("VFSPath = %q, want /tmp/vfs.csv", cfg.VFSPath)
	}
	if cfg.ScriptPath != path {
		t.Errorf("ScriptPath = %q, want %q", cfg.ScriptPath, path)
	}
}

// TestParseEqualsSyntax checks the "--vfs=value" form.
func TestParseEqualsSyntax(t *testing.T) {
	cfg, err := config.Parse([]string{"--vfs=/tmp/data.csv"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VFSPath != "/tmp/data.csv" {
		t.Errorf("VFSPath = %q, want /tmp/data.csv", cfg.VFSPath)
	}
}

// TestParseUnknownFlag checks that unknown flags are rejected.
func TestParseUnknownFlag(t *testing.T) {
	_, err := config.Parse([]string{"--unknown"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

// TestDump checks the key = value output format.
func TestDump(t *testing.T) {
	cfg := &config.Config{
		VFSPath:    "/tmp/vfs.csv",
		ScriptPath: "/tmp/startup.txt",
	}

	var buf bytes.Buffer
	cfg.Dump(&buf)

	got := buf.String()
	if !strings.Contains(got, "vfs.path = /tmp/vfs.csv") {
		t.Errorf("dump missing vfs.path line, got:\n%s", got)
	}
	if !strings.Contains(got, "script.path = /tmp/startup.txt") {
		t.Errorf("dump missing script.path line, got:\n%s", got)
	}
}

// writeTempFile creates a temporary file with the given content.
//
// It registers a cleanup function so the file is removed after the test.
func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "startup.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("cannot write temp file: %v", err)
	}
	return path
}
