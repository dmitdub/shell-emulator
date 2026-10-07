package tests

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dmitdub/shell-emulator/src/script"
)

// TestLoadSimple checks that plain commands are read as-is.
func TestLoadSimple(t *testing.T) {
	path := writeScript(t, "ls\ncd /tmp\n")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"ls", "cd /tmp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestLoadSkipsComments checks that comment lines are ignored.
func TestLoadSkipsComments(t *testing.T) {
	path := writeScript(t, "// comment\nls\n// another\ncd /tmp\n")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"ls", "cd /tmp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestLoadSkipsBlankLines checks that blank lines are ignored.
func TestLoadSkipsBlankLines(t *testing.T) {
	path := writeScript(t, "ls\n\n   \ncd /tmp\n")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"ls", "cd /tmp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestLoadTrimsWhitespace checks that surrounding spaces are removed.
func TestLoadTrimsWhitespace(t *testing.T) {
	path := writeScript(t, "   ls   \n\tcd /tmp\t\n")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"ls", "cd /tmp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestLoadEmpty checks that an empty file yields no commands.
func TestLoadEmpty(t *testing.T) {
	path := writeScript(t, "")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

// TestLoadCommentsOnly checks a file that contains only comments.
func TestLoadCommentsOnly(t *testing.T) {
	path := writeScript(t, "// one\n// two\n")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

// TestLoadMissingFile checks that a missing file yields an error.
func TestLoadMissingFile(t *testing.T) {
	_, err := script.Load("/no/such/file.txt")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

// TestLoadWithQuotedArgs checks that quotes are preserved in commands.
//
// Quotes are handled later by the parser; the script loader must
// pass them through unchanged.
func TestLoadWithQuotedArgs(t *testing.T) {
	path := writeScript(t, `cd "/home/my folder"`+"\n")

	got, err := script.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{`cd "/home/my folder"`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// writeScript creates a temporary script file with the given content.
//
// It registers a cleanup function so the file is removed after the test.
func writeScript(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "startup.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("cannot write temp script: %v", err)
	}
	return path
}
