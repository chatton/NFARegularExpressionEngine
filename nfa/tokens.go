package nfa

import "unicode"

type Token interface {
	Val() string
	Matches(r rune) bool
}

type CharacterClassToken struct {
	val    string
	negate bool
}

func (t CharacterClassToken) Val() string {
	return t.val
}

// example every character in the character class
// if the rune in question matches any of them
// it is a match
func (t CharacterClassToken) Matches(r rune) bool {
	for _, char := range t.val {
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

type WordToken struct {
	val    string
	negate bool
}

func (t WordToken) Val() string {
	return t.val
}

func (t WordToken) Matches(r rune) bool {
	return unicode.IsLetter(r)
}

type DigitToken struct {
	val    string
	negate bool
}

func (t DigitToken) Matches(r rune) bool {
	if t.negate {
		return !unicode.IsDigit(r)
	}
	return unicode.IsDigit(r)
}

func (t DigitToken) Val() string {
	return t.val
}

type SpaceToken struct {
	val    string
	negate bool
}

func (t SpaceToken) Matches(r rune) bool {
	if t.negate {
		return r != ' '
	}
	return r == ' '
}

func (t SpaceToken) Val() string {
	return t.val
}
