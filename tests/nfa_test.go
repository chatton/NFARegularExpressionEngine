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
		{"(a?).(b?).(c?)", "abc", true},
		{"(a?).(b?).(c?)", "ac", true},
		{"(a?).(b?).(c?)", "c", true},
		{"a.b.c.(d|e)", "abcd", true},
		{"a.b.c.(d|e)", "abce", true},
		{"(\\d|\\w)*", "123", true},
		{"(\\d|\\w)*", "abc", true},
		{"(\\d|\\w)*", "abc123", true},
		{"a.b.c.d", "abc", false},
		{"a.b|c.d", "ab", true},
		{"a.b|c.d", "ad", false},
		{`\d*.\w*`, "123hello", true},
		{`\d*.\w*`, "hello", true},
		{`\d*.\w*`, "123", true},
		{`\d+`, "123", true},
		{`\d+`, "", false},
		{`\d+.\w+`, "123", false},
		{`\d+.\w+`, "abc", false},
		{`\d+.\w+`, "123abc", true},
		{"[123].[abc]", "1a", true},
		{"[123].[abc]", "12", false},
		{"[123].[abc]", "123abc", false},
		{"[123].[abc]", "3a", true},
		{"[123]+.[abc]*", "12321312", true},
		{"[123]+.[abc]*", "12321312abc", true},
		{"h.e.l.l.o.w.(o?).r.l.d", "helloworld", true},
		{"h.e.l.l.o.w.(o?).r.l.d", "hellowrld", true},
		{"h.e.l.l.o.w.(o?).r.l.d", "hellowrl", false},
	}

	for _, data := range testData {
		n := nfa.Compile(data.NfaCreationRegex)
		result := n.Matches(data.InputRegex)
		testPassed := result == data.Expected
		if !testPassed {
			t.Error("Test failed. Pattern was: '", data.InputRegex, "' Input was : '", data.InputRegex, "' Match: ", data.Expected)
		}
	}
}
