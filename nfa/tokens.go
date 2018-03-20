package nfa

import "unicode"

type Token interface {
	Val() string
	Matches(r rune) bool
}

type CharacterClassToken struct {
	val string
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
			return true
		}
	}
	return false
}

type WordToken struct {
	val string
}

func (t WordToken) Val() string {
	return t.val
}

func (t WordToken) Matches(r rune) bool {
	return unicode.IsLetter(r)
}

type DigitToken struct {
	val string
}

func (t DigitToken) Matches(r rune) bool {
	return unicode.IsDigit(r)
}

func (t DigitToken) Val() string {
	return t.val
}
