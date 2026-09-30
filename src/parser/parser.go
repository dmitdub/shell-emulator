// Package parser turns a raw input line into a list of arguments.
//
// It expands environment variables and respects single and double
// quotes, so that arguments containing spaces are kept together.
package parser

import (
	"fmt"
	"os"
	"strings"
)

// Parse splits the given line into arguments.
//
// Environment variables such as $HOME are expanded before splitting.
// Arguments wrapped in single or double quotes are kept as a single
// argument. Parse returns an error if quotes are unbalanced.
func Parse(line string) ([]string, error) {
	expanded := os.ExpandEnv(line)

	var (
		args    []string
		current strings.Builder
		quote   rune
	)

	flush := func() {
		if current.Len() > 0 {
			args = append(args, current.String())
			current.Reset()
		}
	}

	for _, r := range expanded {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
		case r == ' ' || r == '\t':
			flush()
		default:
			current.WriteRune(r)
		}
	}

	if quote != 0 {
		return nil, fmt.Errorf("unbalanced quote")
	}

	flush()
	return args, nil
}
