package command

import (
	"fmt"
	"strings"
	"unicode"
)

// ParseLine splits a command line into individual commands and their arguments.
// It handles quoted strings and escape characters.
func ParseLine(line string) ([][]string, error) {
	var commands [][]string
	var currentArgs []string
	var currentToken strings.Builder
	inQuote := false
	escaped := false

	for _, r := range line {
		if escaped {
			currentToken.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if inQuote {
			if r == '"' {
				inQuote = false
			} else {
				currentToken.WriteRune(r)
			}
			continue
		}

		if r == '"' {
			inQuote = true
			continue
		}

		if r == ';' {
			if currentToken.Len() > 0 {
				currentArgs = append(currentArgs, currentToken.String())
				currentToken.Reset()
			}
			if len(currentArgs) > 0 {
				commands = append(commands, currentArgs)
				currentArgs = nil
			}
			continue
		}

		if unicode.IsSpace(r) {
			if currentToken.Len() > 0 {
				currentArgs = append(currentArgs, currentToken.String())
				currentToken.Reset()
			}
			continue
		}

		currentToken.WriteRune(r)
	}

	if inQuote {
		return nil, fmt.Errorf("unclosed quote")
	}

	if currentToken.Len() > 0 {
		currentArgs = append(currentArgs, currentToken.String())
	}
	if len(currentArgs) > 0 {
		commands = append(commands, currentArgs)
	}

	return commands, nil
}
