package command

import (
	"fmt"
	"sort"
	"strings"
)

// RegisterBuiltins registers standard commands with the interpreter.
func (i *Interpreter) RegisterBuiltins() {
	i.Register("help", "Prints help for commands.", i.cmdHelp)
	i.Register("quit", "Exits the application.", i.cmdQuit)
	i.Register("set", "Sets a string variable.", i.cmdSet)
	i.Register("alias", "Creates a command alias.", i.cmdAlias)
	i.Register("bind", "Binds a command to an input.", i.cmdBind)
}

func (i *Interpreter) cmdHelp(args []string) error {
	if len(args) == 0 {
		var names []string
		for name := range i.commands {
			names = append(names, name)
		}
		sort.Strings(names)
		Printf("Available commands:")
		for _, name := range names {
			Printf("  %s", name)
		}
		return nil
	}

	cmdName := args[0]
	if cmd, ok := i.commands[strings.ToLower(cmdName)]; ok {
		Printf("%s\n  %s", cmd.Name, cmd.Help)
	} else {
		return fmt.Errorf("unknown command: %s", cmdName)
	}
	return nil
}

func (i *Interpreter) cmdQuit(args []string) error {
	i.RunQuit()
	return nil
}

func (i *Interpreter) cmdSet(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: set <variable> <value>")
	}
	i.SetString(args[0], args[1])
	return nil
}

func (i *Interpreter) cmdAlias(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: alias <name> <command...>")
	}
	name := args[0]
	// Join the rest of the arguments back into a command line string?
	// This is a bit tricky because the parser already stripped quotes.
	// However, simple reconstruction might be enough for basic aliases.
	// Ideally, we'd want the raw argument string, but our interface is []string.
	// For now, let's join with spaces.
	cmdLine := strings.Join(args[1:], " ")
	i.AddAlias(name, cmdLine)
	return nil
}

func (i *Interpreter) cmdBind(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: bind <input> <command...>")
	}
	input := args[0]
	cmdLine := strings.Join(args[1:], " ")
	i.Bind(input, cmdLine)
	return nil
}

// Bind associates an input event with a command line.
func (i *Interpreter) Bind(input, cmdLine string) {
	if i.bindings == nil {
		i.bindings = make(map[string]string)
	}
	i.bindings[strings.ToLower(input)] = cmdLine
}

// GetBinding returns the command line for an input event.
func (i *Interpreter) GetBinding(input string) string {
	if i.bindings == nil {
		return ""
	}
	return i.bindings[strings.ToLower(input)]
}

// InvertBooleanOps creates a new command line with boolean operations inverted.
// e.g., "+jump" becomes "-jump", "-fire" becomes "+fire".
func InvertBooleanOps(line string) string {
	// We need to parse this robustly to avoid inverting things inside strings.
	// However, since we need to return a string, using the Parser that returns []string might destroy formatting.
	// But for the purpose of "key up", the user probably wants the commands to run.

	// Simple approach: split by semicolon and spaces, look for boolean pattern?
	// "The entire command line should be scanned..."

	// A better way might be to use ParseLine, invert the ops in the args, and reconstruct.
	// But reconstruction is lossy (quotes).

	// Let's iterate manually or regex? Regex is risky with quotes.
	// Let's use a custom iterator similar to Parser but just for replacement?
	// Or just accept that we target token starts.

	// "A command of the from +boolean_variable_name with no arguments"

	// Let's try to parse, identifying boolean commands, and replace them in the original string?
	// Or simpler: Just Tokenize, Invert, Re-assemble with spaces.
	// Re-assembly with spaces is fine for execution usually, unless there were specific quotes.
	// But boolean commands don't use quotes.

	cmds, err := ParseLine(line)
	if err != nil {
		return ""
	}

	var newCmds []string
	for _, args := range cmds {
		if len(args) == 0 {
			continue
		}
		cmd := args[0]
		if strings.HasPrefix(cmd, "+") {
			args[0] = "-" + strings.TrimPrefix(cmd, "+")
		} else if strings.HasPrefix(cmd, "-") {
			args[0] = "+" + strings.TrimPrefix(cmd, "-")
		}

		// Quote arguments if they contain spaces?
		// This is a minimal reconstruction.
		for j, arg := range args {
			if strings.Contains(arg, " ") || strings.Contains(arg, ";") {
				args[j] = fmt.Sprintf("%q", arg)
			}
		}
		newCmds = append(newCmds, strings.Join(args, " "))
	}
	return strings.Join(newCmds, "; ")
}
