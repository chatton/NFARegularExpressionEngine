package tests

import (
	"../nfa"
	"testing"
)

func TestNfa(t *testing.T) {
	testData := []struct {
		NfaCreationRegex string
		InputRegex       string
		Expected         bool
	}{
		{"a.b.c", "abc", true},
		{"a.b.c.d", "abc", false},
		{"(a.b)|(c.d)", "ab", true},
		{`\d*.\w*`, "123hello", true},
		{`\d*.\w*`, "hello", true},
		{`\d*.\w*`, "123", true},
		{`\d+`, "123", true},
		{`\d+`, "", false},
	}

	for _, data := range testData {
		n := nfa.Compile(data.NfaCreationRegex)
		result := n.Matches(data.InputRegex)
		testPassed := result == data.Expected
		if !testPassed {
			t.Error("Test failed. ", data)
		}
	}
}
