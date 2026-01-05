package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/qbradq/annalon/assets"
)

// ExecFile executes commands from a file.
func (i *Interpreter) ExecFile(path string) error {
	var reader io.Reader

	// "If the file_spec begins with an '@' character, strip that character and load the named file relative to the player's data directory"
	// "For now, the player's data directory will be "./data" relative to the working directory."
	if strings.HasPrefix(path, "@") {
		cleanPath := strings.TrimPrefix(path, "@")
		realPath := filepath.Join("data", cleanPath)
		f, err := os.Open(realPath)
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
		// Update path for error reporting to be clearer
		path = realPath
	} else {
		// "Else, load the file relative to the root of the assests file system"
		f, err := assets.FS.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
	}

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// "Empty lines are ignored"
		// "A line starting with the '#' character is a comment line and should be ignored"
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// "Report errors to both stderr and the console's print function"
		// "Report file names and line numbers starting from 1 in errors"
		if err := i.Execute(line); err != nil {
			errMsg := fmt.Sprintf("Error executing %s:%d: %v", path, lineNum, err)
			fmt.Fprintln(os.Stderr, errMsg)
			Printf("%s", errMsg)
		}
	}
	return scanner.Err()
}

func (i *Interpreter) cmdExec(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: exec <file>")
	}
	return i.ExecFile(args[0])
}
