package nfa

import "unicode"

type token interface {
	Val() string
	Matches(r rune) bool
}

type characterClassToken struct {
	val             string
	negate          bool
	caseInsensitive bool
}

func (t characterClassToken) Val() string {
	return t.val
}

// example every character in the character class
// if the rune in question matches any of them
// it is a match
func (t characterClassToken) Matches(r rune) bool {
	for _, char := range t.val {
		if t.caseInsensitive {
			// ignore the case of both of the characters in question
			r = unicode.ToLower(r)
			char = unicode.ToLower(char)
		}

		if r == char {

			if t.negate {
				return false
			}
			return true
		}
	}
	if t.negate {
		return true
	}
	return false
}

type wordToken struct {
	val    string
	negate bool
}

func (t wordToken) Val() string {
	return t.val
}

func (t wordToken) Matches(r rune) bool {
	return unicode.IsLetter(r)
}

type digitToken struct {
	val    string
	negate bool
}

func (t digitToken) Matches(r rune) bool {
	if t.negate {
		return !unicode.IsDigit(r)
	}
	return unicode.IsDigit(r)
}

func (t digitToken) Val() string {
	return t.val
}

type spaceToken struct {
	val    string
	negate bool
}

func (t spaceToken) Matches(r rune) bool {
	if t.negate {
		return r != ' '
	}
	return r == ' '
}

func (t spaceToken) Val() string {
	return t.val
}

type anyToken struct {
	val    string
	negate bool
}

func (t anyToken) Val() string {
	return t.val
}

// Any token will match any character
func (t anyToken) Matches(r rune) bool {
	if t.negate {
		return false
	}
	return true
}
