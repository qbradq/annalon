package command

import (
	"testing"
)

func TestInterpreter(t *testing.T) {
	i := NewInterpreter()
	i.RegisterBuiltins()

	// Test Set/Get String
	i.SetString("foo", "bar")
	if val := i.GetString("foo"); val != "bar" {
		t.Errorf("GetString(foo) = %s, expected bar", val)
	}

	// Test Set/Get Bool
	i.SetBool("flag", true)
	if val := i.GetBool("flag"); !val {
		t.Errorf("GetBool(flag) = false, expected true")
	}

	// Test Execute String Set
	if err := i.Execute("set baz qux"); err != nil {
		t.Errorf("Execute(set) failed: %v", err)
	}
	if val := i.GetString("baz"); val != "qux" {
		t.Errorf("Result of set baz qux = %s, expected qux", val)
	}

	// Test Boolean Shortcuts
	if err := i.Execute("+flag2"); err != nil {
		t.Errorf("Execute(+flag2) failed: %v", err)
	}
	if val := i.GetBool("flag2"); !val {
		t.Errorf("Result of +flag2 = false, expected true")
	}

	if err := i.Execute("-flag2"); err != nil {
		t.Errorf("Execute(-flag2) failed: %v", err)
	}
	if val := i.GetBool("flag2"); val {
		t.Errorf("Result of -flag2 = true, expected false")
	}
}

func TestAlias(t *testing.T) {
	i := NewInterpreter()
	i.RegisterBuiltins()

	// Register a mock command to verify execution
	executed := false
	var receivedArgs []string
	i.Register("echo", "Echoes args", func(args []string) error {
		executed = true
		receivedArgs = args
		return nil
	})

	// Create alias
	if err := i.Execute("alias e echo hello"); err != nil {
		t.Fatalf("Failed to create alias: %v", err)
	}

	// Execute alias
	if err := i.Execute("e world"); err != nil {
		t.Fatalf("Failed to execute alias: %v", err)
	}

	if !executed {
		t.Error("Alias did not trigger command")
	}
	if len(receivedArgs) != 2 || receivedArgs[0] != "hello" || receivedArgs[1] != "world" {
		t.Errorf("Alias expansion failed. Got args: %v, expected [hello world]", receivedArgs)
	}
}

func TestInverseBooleanCommands(t *testing.T) {
	i := NewInterpreter()

	tests := []struct {
		input    string
		expected string
	}{
		{"+forward", "-forward"},
		{"-jump", "+jump"},
		{"+forward; -jump", "-forward; +jump"},
		{"say hello", ""},
		{"+forward; say hello; -jump", "-forward; +jump"},
	}

	for _, tt := range tests {
		result := i.GetInverseBooleanCommands(tt.input)
		if result != tt.expected {
			t.Errorf("GetInverseBooleanCommands(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
