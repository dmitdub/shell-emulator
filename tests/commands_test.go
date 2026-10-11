package tests

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/dmitdub/shell-emulator/src/commands"
	"github.com/dmitdub/shell-emulator/src/vfs"
)

// newTestContext creates a Context with a small VFS and buffers.
func newTestContext(t *testing.T) (*commands.Context, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	root := &vfs.Node{ID: 0, Name: "/", Dir: true}
	home := &vfs.Node{ID: 1, Name: "home", Dir: true}
	readme := &vfs.Node{ID: 2, Name: "readme.txt", Content: "Hello\n"}
	empty := &vfs.Node{ID: 3, Name: "empty.txt", Content: ""}

	home.AddChild(readme)
	home.AddChild(empty)
	root.AddChild(home)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	ctx := commands.New(root, out, errOut)
	return ctx, out, errOut
}

// ---- ls ----

// TestLsRoot checks that ls lists the root directory contents.
func TestLsRoot(t *testing.T) {
	ctx, out, _ := newTestContext(t)

	ctx.Ls(nil)

	got := out.String()
	if !strings.Contains(got, "home/") {
		t.Errorf("ls output %q does not contain home/", got)
	}
}

// TestLsLongDir checks the long format for a directory.
func TestLsLongDir(t *testing.T) {
	ctx, out, _ := newTestContext(t)

	ctx.Ls([]string{"-l", "/"})

	if !strings.Contains(out.String(), "drwxr-xr-x") {
		t.Errorf("ls -l output %q missing dir permissions", out.String())
	}
}

// TestLsLongFile checks the long format for a file.
func TestLsLongFile(t *testing.T) {
	ctx, out, _ := newTestContext(t)

	ctx.Ls([]string{"-l", "/home"})

	if !strings.Contains(out.String(), "-rw-r--r--") {
		t.Errorf("ls -l output %q missing file permissions", out.String())
	}
}

// TestLsMissingPath checks the error for a missing path.
func TestLsMissingPath(t *testing.T) {
	ctx, _, errOut := newTestContext(t)

	ctx.Ls([]string{"/no/such"})

	if !strings.Contains(errOut.String(), "No such file or directory") {
		t.Errorf("unexpected stderr: %q", errOut.String())
	}
}

// TestLsFile checks that listing a file prints its name.
func TestLsFile(t *testing.T) {
	ctx, out, _ := newTestContext(t)

	ctx.Ls([]string{"/home/readme.txt"})

	if !strings.Contains(out.String(), "readme.txt") {
		t.Errorf("ls file output %q missing filename", out.String())
	}
}

// ---- cd ----

// TestCdRelative checks relative directory change.
func TestCdRelative(t *testing.T) {
	ctx, _, _ := newTestContext(t)

	ctx.Cd([]string{"home"})

	if ctx.Cwd.Name != "home" {
		t.Errorf("Cwd = %q, want home", ctx.Cwd.Name)
	}
}

// TestCdParent checks the ".." path.
func TestCdParent(t *testing.T) {
	ctx, _, _ := newTestContext(t)

	ctx.Cd([]string{"home"})
	ctx.Cd([]string{".."})

	if ctx.Cwd.Name != "/" {
		t.Errorf("Cwd = %q, want /", ctx.Cwd.Name)
	}
}

// TestCdMissing checks the error for a missing path.
func TestCdMissing(t *testing.T) {
	ctx, _, errOut := newTestContext(t)

	ctx.Cd([]string{"/no/such"})

	if !strings.Contains(errOut.String(), "No such file or directory") {
		t.Errorf("unexpected stderr: %q", errOut.String())
	}
}

// TestCdFile checks the error when the target is a file.
func TestCdFile(t *testing.T) {
	ctx, _, errOut := newTestContext(t)

	ctx.Cd([]string{"/home/readme.txt"})

	if !strings.Contains(errOut.String(), "Not a directory") {
		t.Errorf("unexpected stderr: %q", errOut.String())
	}
}

// ---- cat ----

// TestCat checks that cat prints the file content.
func TestCat(t *testing.T) {
	ctx, out, _ := newTestContext(t)

	ctx.Cat([]string{"/home/readme.txt"})

	if !strings.Contains(out.String(), "Hello") {
		t.Errorf("cat output %q missing content", out.String())
	}
}

// TestCatDirectory checks the error when the target is a directory.
func TestCatDirectory(t *testing.T) {
	ctx, _, errOut := newTestContext(t)

	ctx.Cat([]string{"/home"})

	if !strings.Contains(errOut.String(), "Is a directory") {
		t.Errorf("unexpected stderr: %q", errOut.String())
	}
}

// TestCatMissing checks the error for a missing path.
func TestCatMissing(t *testing.T) {
	ctx, _, errOut := newTestContext(t)

	ctx.Cat([]string{"/no/such"})

	if !strings.Contains(errOut.String(), "No such file or directory") {
		t.Errorf("unexpected stderr: %q", errOut.String())
	}
}

// TestCatNoArgs checks the error when no file is provided.
func TestCatNoArgs(t *testing.T) {
	ctx, _, errOut := newTestContext(t)

	ctx.Cat(nil)

	if !strings.Contains(errOut.String(), "missing operand") {
		t.Errorf("unexpected stderr: %q", errOut.String())
	}
}

// ---- uptime ----

// TestUptimeFormat checks that uptime returns a well-formed line.
func TestUptimeFormat(t *testing.T) {
	ctx, out, _ := newTestContext(t)
	ctx.Started = time.Now()

	ctx.Uptime(nil)

	got := out.String()
	if !strings.HasPrefix(got, "up ") {
		t.Errorf("uptime output %q does not start with 'up '", got)
	}
	if !strings.Contains(got, "second") {
		t.Errorf("uptime output %q missing 'second'", got)
	}
}
