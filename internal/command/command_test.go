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

func TestBindAndInvert(t *testing.T) {
	i := NewInterpreter()
	i.RegisterBuiltins()

	i.Execute("bind SPACE +jump")

	binding := i.GetBinding("SPACE")
	if binding != "+jump" {
		t.Errorf("GetBinding(SPACE) = %s, expected +jump", binding)
	}

	inverted := InvertBooleanOps(binding)
	if inverted != "-jump" {
		t.Errorf("InvertBooleanOps(+jump) = %s, expected -jump", inverted)
	}

	complexCmd := "+jump; -crouch; say hello"
	invertedComplex := InvertBooleanOps(complexCmd)
	// Output format depends on InvertBooleanOps implementation (spaces, etc.)
	// We expect "-jump ; +crouch ; say hello" or similar structure
	// Let's rely on basic check
	if msg, _ := ParseLine(invertedComplex); msg[0][0] != "-jump" || msg[1][0] != "+crouch" || msg[2][0] != "say" {
		t.Errorf("InvertBooleanOps failed for complex: %s", invertedComplex)
	}
}
