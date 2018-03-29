package tests

import (
	"../nfa"
	"strconv"
	"testing"
)

func TestCount(t *testing.T) {
	testData := []struct {
		InputRegex, MatchString string
		ExpectedResult          int
	}{
		{"ab", "abab", 2},
		{"a", "aabva", 3},
		{"1*", "11", 4}, // matches the empty string, "1", "1" and "11"
		{"[12][34]", "13", 1},
		{"abcd", "abcdabcdabcf", 2},
	}

	for _, data := range testData {
		result := nfa.Count(data.InputRegex, data.MatchString)
		if result != data.ExpectedResult {
			t.Error("The pattern " + data.MatchString + " had " + strconv.Itoa(result) + " occurrences the string when it should have had: " + strconv.Itoa(data.ExpectedResult))
		}
	}
}

func TestNfa(t *testing.T) {
	testData := []struct {
		NfaCreationRegex string
		InputRegex       string
		Expected         bool
	}{
		{"abc", "abc", true},
		{"1____6", "123456", true},
		{"1____6", "654321", false},
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
		{"ab|c_", "cb", true},
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
		{"hello\\s+world", "hello           world", true},
		{"hello\\s+world", "helloworld", false},
		{"\\s+helloworld\\s+", "         helloworld                       ", true},
		{`^a`, "a", false},
		{`^[123]`, "2", false},
		{`^[123]`, "4", true},
		{`abc^\def`, "abc5ef", false},
		{`abc\def`, "abc5ef", true},
		{`abc\d+ef`, "abc5678765446ef", true},
		{`\d\w+\d`, "1sdfsdfds2", true},
		{`\d\w+\d`, "12", false},
		{`\d\w?\d`, "12", true},
		{`\\hello`, `\hello`, true},
		{`hel_o`, `hello`, true},
		{`h_llo`, `hello`, true},
		{`_ello`, `hello`, true},
		{"(?i)HELlo", `helLO`, true},
		{"HELlo", `hello`, false},
		{"^[1]", `h`, true},
		{"^\\d", `h`, true},
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
