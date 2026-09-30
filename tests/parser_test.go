package tests

import (
	"reflect"
	"testing"

	"github.com/dmitdub/shell-emulator/src/parser"
)

// TestParseSimple checks that plain arguments are split correctly.
func TestParseSimple(t *testing.T) {
	got, err := parser.Parse("ls -la /tmp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"ls", "-la", "/tmp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestParseDoubleQuoted checks that double-quoted arguments are kept together.
func TestParseDoubleQuoted(t *testing.T) {
	got, err := parser.Parse(`cd "/home/my folder"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"cd", "/home/my folder"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestParseSingleQuoted checks that single-quoted arguments are kept together.
func TestParseSingleQuoted(t *testing.T) {
	got, err := parser.Parse(`echo 'hello world'`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"echo", "hello world"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestParseEnvExpansion checks that environment variables are expanded.
func TestParseEnvExpansion(t *testing.T) {
	const homeValue = "/home/tester"
	t.Setenv("HOME", homeValue)

	got, err := parser.Parse("cd $HOME")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"cd", homeValue}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestParseUnbalancedQuote checks error handling for unclosed quotes.
func TestParseUnbalancedQuote(t *testing.T) {
	_, err := parser.Parse(`ls "unclosed`)
	if err == nil {
		t.Fatal("expected error for unbalanced quote, got nil")
	}
}

// TestParseEmpty checks that a whitespace-only line yields no arguments.
func TestParseEmpty(t *testing.T) {
	got, err := parser.Parse("   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

// TestParseMixedQuotes checks a line with several quoted arguments.
func TestParseMixedQuotes(t *testing.T) {
	got, err := parser.Parse(`cp 'src file' "dst file"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"cp", "src file", "dst file"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
