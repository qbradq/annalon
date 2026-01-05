package command

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecFile_Mock(t *testing.T) {
	// Create a temporary directory for data.
	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "data")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a dummy config file
	configFile := filepath.Join(dataDir, "test.con")
	content := []byte("set foo bar\nalias a1 set foo baz\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Change working directory to tmpDir so "./data" works as expected
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	i := NewInterpreter()
	i.RegisterBuiltins()

	err := i.ExecFile("@test.con")
	if err != nil {
		t.Fatalf("ExecFile failed: %v", err)
	}

	if val := i.GetString("foo"); val != "bar" {
		t.Errorf("expected foo=bar, got %s", val)
	}
}
