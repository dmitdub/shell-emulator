// Package script reads and executes startup scripts of the emulator.
//
// A startup script is a plain text file with one emulator command
// per line. Lines that start with "//" are treated as comments and
// skipped. Empty lines are ignored.
package script

import (
	"bufio"
	"os"
	"strings"
)

// commentPrefix marks a line as a comment.
const commentPrefix = "//"

// Load reads the given file and returns its commands.
//
// Comments and blank lines are removed. Surrounding whitespace
// is trimmed from every command. Load returns an error if the
// file cannot be opened or read.
func Load(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if command := cleanLine(scanner.Text()); command != "" {
			commands = append(commands, command)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return commands, nil
}

// cleanLine trims whitespace and drops comment lines.
//
// It returns an empty string when the line contains no command.
func cleanLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, commentPrefix) {
		return ""
	}
	return trimmed
}
