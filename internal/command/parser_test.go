package command

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		input    string
		expected [][]string
	}{
		{
			input:    "cmd1 arg1",
			expected: [][]string{{"cmd1", "arg1"}},
		},
		{
			input:    "cmd1; cmd2",
			expected: [][]string{{"cmd1"}, {"cmd2"}},
		},
		{
			input:    "cmd1 arg1; cmd2 arg2 arg3",
			expected: [][]string{{"cmd1", "arg1"}, {"cmd2", "arg2", "arg3"}},
		},
		{
			input:    `cmd "quoted string"`,
			expected: [][]string{{"cmd", "quoted string"}},
		},
		{
			input:    `cmd "quote with \" escaped"`,
			expected: [][]string{{"cmd", `quote with " escaped`}},
		},
		{
			input:    `cmd "escaped backslash \\"`,
			expected: [][]string{{"cmd", `escaped backslash \`}},
		},
		{
			input:    "  cmd1   arg1  ;  cmd2  ",
			expected: [][]string{{"cmd1", "arg1"}, {"cmd2"}},
		},
		{
			input:    "",
			expected: nil,
		},
		{
			input:    ";;",
			expected: nil,
		},
	}

	for _, test := range tests {
		result, err := ParseLine(test.input)
		if err != nil {
			t.Errorf("ParseLine(%q) returned error: %v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ParseLine(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseLineError(t *testing.T) {
	_, err := ParseLine(`cmd "unclosed quote`)
	if err == nil {
		t.Error("Expected error for unclosed quote, got nil")
	}
}
