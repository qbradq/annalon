package command

import (
	"fmt"
	"strings"
)

// Command is a function that executes a command.
// It receives the arguments passed to the command.
type Command func(args []string) error

var printer func(string)

// SetPrinter sets the function used for printing output.
func SetPrinter(fn func(string)) {
	printer = fn
}

// Printf formats according to a format specifier and writes to the printer.
func Printf(format string, a ...any) {
	if printer != nil {
		printer(fmt.Sprintf(format, a...))
	}
}

// CommandInfo holds the command function and its help text.
type CommandInfo struct {
	Name string
	Help string
	Fn   Command
}

// Interpreter manages the command environment.
type Interpreter struct {
	commands    map[string]CommandInfo
	aliases     map[string]string
	stringVars  map[string]string
	boolVars    map[string]bool
	quitHandler func()
	bindings    map[string]string
}

// NewInterpreter creates a new command interpreter.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		commands:   make(map[string]CommandInfo),
		aliases:    make(map[string]string),
		stringVars: make(map[string]string),
		boolVars:   make(map[string]bool),
		bindings:   make(map[string]string),
	}
}

// Register adds a new command to the interpreter.
func (i *Interpreter) Register(name, help string, cmd Command) {
	i.commands[strings.ToLower(name)] = CommandInfo{
		Name: name,
		Help: help,
		Fn:   cmd,
	}
}

// SetQuitHandler sets the function to be called when the quit command is executed.
func (i *Interpreter) SetQuitHandler(handler func()) {
	i.quitHandler = handler
}

// RunQuit executes the quit handler if set.
func (i *Interpreter) RunQuit() {
	if i.quitHandler != nil {
		i.quitHandler()
	}
}

// SetString sets a string variable.
func (i *Interpreter) SetString(name, value string) {
	i.stringVars[strings.ToLower(name)] = value
}

// GetString gets a string variable.
func (i *Interpreter) GetString(name string) string {
	return i.stringVars[strings.ToLower(name)]
}

// SetBool sets a boolean variable.
func (i *Interpreter) SetBool(name string, value bool) {
	i.boolVars[strings.ToLower(name)] = value
}

// GetBool gets a boolean variable.
func (i *Interpreter) GetBool(name string) bool {
	return i.boolVars[strings.ToLower(name)]
}

// AddAlias adds a command alias.
func (i *Interpreter) AddAlias(name, cmdLine string) {
	i.aliases[strings.ToLower(name)] = cmdLine
}

// Execute parses and runs a command line.
// A command line can contain multiple commands separated by semicolons.
func (i *Interpreter) Execute(line string) error {
	cmds, err := ParseLine(line)
	if err != nil {
		return err
	}

	for _, args := range cmds {
		if len(args) == 0 {
			continue
		}

		cmdName := strings.ToLower(args[0])

		// Handle boolean variable shortcuts
		if strings.HasPrefix(cmdName, "+") {
			varName := strings.TrimPrefix(cmdName, "+")
			if len(varName) > 0 {
				i.SetBool(varName, true)
				continue
			}
		}
		if strings.HasPrefix(cmdName, "-") {
			varName := strings.TrimPrefix(cmdName, "-")
			if len(varName) > 0 {
				i.SetBool(varName, false)
				continue
			}
		}

		// Handle aliases
		// Note: The prompt example `alias exit quit` implies direct substitution.
		// `alias a1 cmd2 foo bar` -> `a1 baz` = `cmd2 foo bar baz`.
		// So we need to check if it's an alias, parse the alias expansion, and prepend it.
		// Be careful of recursion.
		if expansion, ok := i.aliases[cmdName]; ok {
			// Basic recursion check or depth limit could be added here if needed.
			// For now, let's just do a single level expansion or re-parse?
			// The example `cmd1; a1 baz; cmd3` -> `cmd1; cmd2 foo bar baz; cmd3`
			// This suggests alias expansion happens at the command word level.

			// If we re-parse the expansion, we might handle complexity better.
			// But the prompt says `alias name command [argument_list]`.
			// `alias creates a new command word name that acts as command call with the argument list.`

			// So `a1` -> `cmd2 foo bar`.
			// `a1 baz` -> `cmd2 foo bar` + `baz`.

			aliasParts, err := ParseLine(expansion)
			if err == nil && len(aliasParts) > 0 {
				// We only take the first command from the alias expansion if it parses to multiple?
				// "Binds a command alias, such that when the alias is used in place of a command word, the alias expands."
				// A simple implementation is to replace args[0] with the alias parts and continue.
				// However, the alias string might contain multiple tokens.
				// e.g. expansion "cmd2 foo bar".
				// args is ["a1", "baz"]
				// newArgs should be ["cmd2", "foo", "bar", "baz"]

				// There is a edge case: what if the alias itself contains multiple commands?
				// "alias name command [argument_list]" implies a single command.
				// But the expansion logic should probably just prepend the tokens.

				// Let's assume alias is a single command chain.
				// We'll take the first command from the parsed alias.
				expandedArgs := append(aliasParts[0], args[1:]...)

				// Recursive lookup? The prompt doesn't strictly require it but it's good practice.
				// For now, let's just set args to expandedArgs and let the loop continue with the new cmdName.
				args = expandedArgs
				cmdName = strings.ToLower(args[0])
			}
		}

		if cmdInfo, ok := i.commands[cmdName]; ok {
			if err := cmdInfo.Fn(args[1:]); err != nil {
				return fmt.Errorf("command '%s' failed: %w", cmdName, err)
			}
		} else {
			return fmt.Errorf("unknown command: %s", cmdName)
		}
	}
	return nil
}
