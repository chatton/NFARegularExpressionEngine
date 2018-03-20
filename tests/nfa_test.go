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
		{"abc", "abc", true},
		{"a?b?c?", "abc", true},
		{"a?b?c?", "ac", true},
		{"a?b?c?", "c", true},
		{"abc(d|e)", "abcd", true},
		{"abc(d|e)", "abce", true},
		{"(a|b)cde", "acde", true},
		{"(a|b)(c|d)", "ad", true},
		{"(\\d|\\w)*", "123", true},
		{"(\\d|\\w)*", "abc", true},
		{"(\\d|\\w)*", "abc123", true},
		{"abcd", "abc", false},
		{"ab|cd", "ab", true},
		{"ab|cd", "ad", false},
		{`\d*\w*`, "123hello", true},
		{`\d*\w*`, "hello", true},
		{`\d*\w*`, "123", true},
		{`\d+`, "123", true},
		{`\d+`, "", false},
		{`\d+\w+`, "123", false},
		{`\d+\w+`, "abc", false},
		{`\d+\w+`, "123abc", true},
		{"[123][abc]", "1a", true},
		{"[123][abc]", "12", false},
		{"[123][abc]", "123abc", false},
		{"[123][abc]", "3a", true},
		{"[123]+[abc]*", "12321312", true},
		{"[123]+[abc]*", "12321312abcabcaaabbbccc", true},
		{"hellow(o?)rld", "helloworld", true},
		{"hellow(o?)rld", "hellowrld", true},
		{"hellow(o?)rld", "hellowrl", false},
	}

	for _, data := range testData {
		n := nfa.Compile(data.NfaCreationRegex)
		result := n.Matches(data.InputRegex)
		testPassed := result == data.Expected
		if !testPassed {
			t.Error("Test failed. Pattern was: '", data.InputRegex, "' Input was : '", data.InputRegex, "' Match: ", testPassed)
		}
	}
}
