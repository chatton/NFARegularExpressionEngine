package nfa

import (
	"github.com/golang-collections/collections/stack" // useful stack data structure
	"unicode"
)

type nfa struct {
	Initial, Accept *State
	tokens          []Token
}

// provide mechanism for creating compiled Nfas so you don't need to construct the nfa each time you use it.
func Compile(infix string) *nfa {
	tokens := InfixToPostfix(infix)
	n := tokensToNfa(tokens)
	n.tokens = tokens
	return n
}

type State struct {
	symbol       interface{}
	edge1, edge2 *State
}

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

func tokensToNfa(tokens []Token) *nfa {
	nfaStack := stack.New()
	for _, tok := range tokens {
		switch tok.(Token).Val() {
		case ".":
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*nfa)
			frag1 := nfaStack.Pop().(*nfa)
			// join them together
			frag1.Accept.edge1 = frag2.Initial
			// place the single fragment on top of the stack
			nfaStack.Push(&nfa{Initial: frag1.Initial, Accept: frag2.Accept})
		case "|":
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*nfa)
			frag1 := nfaStack.Pop().(*nfa)
			// create a new initial state which points at both fragments
			initial := State{edge1: frag1.Initial, edge2: frag2.Initial}
			accept := State{}
			// point both fragments at the new accept state
			frag1.Accept.edge1 = &accept
			frag2.Accept.edge1 = &accept

			// add the new fragment to the stack.
			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})
		case "*":
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			accept := State{}
			// create and edge pointing back at itself and one towards the new accept state
			initial := State{edge1: frag.Initial, edge2: &accept}
			// the accept state points back at the initial state
			frag.Accept.edge1 = frag.Initial
			// the second edge loops around
			frag.Accept.edge2 = &accept
			// add it to the stack
			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})
		case "+":
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			accept := State{}
			initial := State{edge1: frag.Initial, edge2: &accept}

			frag.Accept.edge1 = &initial
			frag.Accept.edge2 = nil

			nfaStack.Push(&nfa{Initial: frag.Initial, Accept: &accept})
		case "?":
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			// create a new state that points to the existing item and also the accept state
			initial := State{edge1: frag.Initial, edge2: frag.Accept}
			// push the new Nfa onto the stack
			nfaStack.Push(&nfa{Initial: &initial, Accept: frag.Accept})
		default:
			accept := State{}
			initial := State{edge1: &accept, symbol: tok}

			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})

		}
	}

	return nfaStack.Pop().(*nfa)
}

//func postFixToNfa(postfix string) *nfa {
//
//	nfaStack := stack.New()
//	for _, r := range postfix {
//		switch r {
//		case '.':
//			// take 2 elements off the stack
//			frag2 := nfaStack.Pop().(*nfa)
//			frag1 := nfaStack.Pop().(*nfa)
//			// join them together
//			frag1.Accept.edge1 = frag2.Initial
//			// place the single fragment on top of the stack
//			nfaStack.Push(&nfa{Initial: frag1.Initial, Accept: frag2.Accept})
//		case '|':
//			// take 2 elements off the stack
//			frag2 := nfaStack.Pop().(*nfa)
//			frag1 := nfaStack.Pop().(*nfa)
//			// create a new initial state which points at both fragments
//			initial := State{edge1: frag1.Initial, edge2: frag2.Initial}
//			accept := State{}
//			// point both fragments at the new accept state
//			frag1.Accept.edge1 = &accept
//			frag2.Accept.edge1 = &accept
//
//			// add the new fragment to the stack.
//			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})
//		case '*':
//			// take a single element off of the stack
//			frag := nfaStack.Pop().(*nfa)
//			accept := State{}
//			// create and edge pointing back at itself and one towards the new accept state
//			initial := State{edge1: frag.Initial, edge2: &accept}
//			// the accept state points back at the initial state
//			frag.Accept.edge1 = frag.Initial
//			// the second edge loops around
//			frag.Accept.edge2 = &accept
//			// add it to the stack
//			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})
//		case '+':
//			// take a single element off of the stack
//			frag := nfaStack.Pop().(*nfa)
//			accept := State{}
//			initial := State{edge1: frag.Initial, edge2: &accept}
//
//			frag.Accept.edge1 = &initial
//			frag.Accept.edge2 = nil
//
//			nfaStack.Push(&nfa{Initial: frag.Initial, Accept: &accept})
//		case '?':
//			// take a single element off of the stack
//			frag := nfaStack.Pop().(*nfa)
//			// create a new state that points to the existing item and also the accept state
//			initial := State{edge1: frag.Initial, edge2: frag.Accept}
//			// push the new Nfa onto the stack
//			nfaStack.Push(&nfa{Initial: &initial, Accept: frag.Accept})
//		default:
//			accept := State{}
//			initial := State{edge1: &accept}
//			switch r {
//			case 'd':
//				initial.symbol = DigitToken{string(r)}
//			case 'w':
//				initial.symbol = WordToken{string(r)}
//			default:
//				initial.symbol = CharacterClassToken{string(r)}
//			}
//
//			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})
//		}
//	}
//
//	if nfaStack.Len() != 1 {
//		panic("Nfa stack didn't have just a single element in it.")
//	}
//	res := nfaStack.Pop().(*nfa)
//	return res
//}

func IsEmpty(s *stack.Stack) bool {
	return s.Len() == 0
}

func Tokenize(infix string) []interface{} {

	var tokens []interface{}

	var s string
	var escapeStr string
	appendToS := false
	wantsToEscape := false
	for _, r := range infix {

		if wantsToEscape { // there was a backslash escape the character
			escapeStr += string(r)
			wantsToEscape = false
			switch r {
			case 'd': // \d
				tokens = append(tokens, DigitToken{val: `\d`})
			case 'w': // \w
				tokens = append(tokens, WordToken{val: `\w`})
			case '\\': // a second backslash, treat it as a backslash literal
				tokens = append(tokens, CharacterClassToken{val: string(r)})
			}
			continue
		}

		if r == '\\' { // want to escape a character
			wantsToEscape = true
			escapeStr = ""
			continue
		}
		startingClass, endingClass := r == '[', r == ']'

		// don't want to append the last element in a token
		if appendToS && startingClass {
			s += string(r) // add as a single character of a multi character token
		} else if !startingClass && !endingClass { // add the single character as a token
			tokens = append(tokens, CharacterClassToken{string(r)})
		}

		if startingClass { // we're going to start a multi character token
			appendToS = true
		} else if endingClass { // reached end of character class
			tokens = append(tokens, CharacterClassToken{s}) // add the full string as a single token
			s = ""
			appendToS = false // stop building up character class
		}
	}

	return tokens
}

func InfixToPostfix(infix string) []Token {

	specials := map[string]int{"*": 10, ".": 9, "+": 8, "|": 7, "?": 6}

	postfix := stack.New()
	tempStack := stack.New()

	tokens := Tokenize(infix)
	for _, tok := range tokens {
		val := tok.(Token).Val()

		switch {
		case val == "(":
			tempStack.Push(tok)
		case val == ")":
			for tempStack.Peek().(Token).Val() != "(" {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Pop()
		case specials[val] > 0:
			for !IsEmpty(tempStack) && specials[val] <= specials[tempStack.Peek().(Token).Val()] {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Push(tok)
		default:
			postfix.Push(tok)
		}
	}

	for !IsEmpty(tempStack) {
		postfix.Push(tempStack.Pop())
	}

	var result []Token

	for !IsEmpty(postfix) {
		// insert every element at the start of the list.
		result = append([]Token{postfix.Pop().(Token)}, result...)
	}

	return result
}

//returns all the possible states from the given state, including
//state transitions with e arrows
func addState(possibilities []*State, from, to *State) []*State {
	seen := make(map[*State]bool) // keep track of the states we've already visited

	states := stack.New()
	states.Push(from)

	for !IsEmpty(states) { // keep looking until dead end
		next := states.Pop().(*State)

		if seen[next] { // we may be looking at a state we've already seen
			continue
		}

		// mark the state as having been seen already so we don't examine it again.
		seen[next] = true
		possibilities = append(possibilities, next) // this is a valid destination

		// e arrow if symbol == nil
		// zero value for Token is nil
		if next != to && next.symbol == nil {
			states.Push(next.edge1)
			if next.edge2 != nil {
				states.Push(next.edge2)
			}
		}
	}

	return possibilities
}

func (n *nfa) Matches(matchString string) bool {

	var current []*State
	current = addState(current, n.Initial, n.Accept)
	var next []*State

	for _, r := range matchString {
		for _, curr := range current {
			if curr.symbol != nil && curr.symbol.(Token).Val() == string(r) {
				next = addState(next[:], curr.edge1, n.Accept)
			}
		}
		current, next = next, []*State{}
	}

	for _, curr := range current {
		if curr == n.Accept {
			return true
		}
	}

	return false
}
