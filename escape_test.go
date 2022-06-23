package sjson

import (
	"testing"
)

type testCase struct {
	in       string
	expected string
	err      bool
}

var testCases []testCase = []testCase{
	{"non-escaped", "non-escaped", false},
	{"\\n\\t\\b\\f\\r", "\n\t\b\f\r", false},
	{`\"\"\//\\`, `""//\`, false},
	{`\u263A`, "\u263A", false},
	{`\1\m`, ``, true},
	{`\u263Z`, ``, true},
}

func TestEscaping(t *testing.T) {
	for _, c := range testCases {
		out, err := EscapeUTF8([]byte(c.in))
		if err != nil && !c.err {
			t.Errorf("Error for input '%v' when no error expected", c.in)
		} else if c.err && err == nil {
			t.Errorf("Expected error for input '%v' but no error.", c.in)
		} else if c.expected != out {
			t.Errorf("Output '%v' does not match expected '%v' for input %v", out, c.expected, c.in)
		}
	}
}
