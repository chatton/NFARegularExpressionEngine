package nfa

import (
	"github.com/golang-collections/collections/stack" // useful stack data structure
	"strings"
	"unicode"
)

type nfa struct {
	Initial, Accept *State
	postfix         string
}

// provide mechanism for creating compiled Nfas so you don't need to construct the nfa each time you use it.
func Compile(infix string) *nfa {
	postfix := InfixToPostfix(infix)
	n := postFixToNfa(postfix)
	n.postfix = postfix
	return n
}

// use the saved nfa to match the provided string
//func (n *nfa) Matches(matchString string) bool {
//	return postfixMatch(n.postfix, matchString)
//}

type State struct {
	symbol       interface{}
	edge1, edge2 *State
}

type Token interface {
	Matches(r rune) bool
}

type CharacterClassToken struct {
	val string
}

//func (t *CharacterClassToken) Val() string {
//	return t.val
//}

// example every character in the character class
// if the rune in question matches any of them
// it is a match
func (t *CharacterClassToken) Matches(r rune) bool {
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

func (t *WordToken) Matches(r rune) bool {
	return unicode.IsLetter(r)
}

type DigitToken struct {
	val string
}

func (t *DigitToken) Matches(r rune) bool {
	return unicode.IsDigit(r)
}

func postFixToNfa(postfix string) *nfa {

	nfaStack := stack.New()
	for _, r := range postfix {
		switch r {
		case '.':
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*nfa)
			frag1 := nfaStack.Pop().(*nfa)
			// join them together
			frag1.Accept.edge1 = frag2.Initial
			// place the single fragment on top of the stack
			nfaStack.Push(&nfa{Initial: frag1.Initial, Accept: frag2.Accept})
		case '|':
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
		case '*':
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
		case '+':
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			accept := State{}
			initial := State{edge1: frag.Initial, edge2: &accept}

			frag.Accept.edge1 = &initial
			frag.Accept.edge2 = nil

			nfaStack.Push(&nfa{Initial: frag.Initial, Accept: &accept})
		case '?':
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			// create a new state that points to the existing item and also the accept state
			initial := State{edge1: frag.Initial, edge2: frag.Accept}
			// push the new Nfa onto the stack
			nfaStack.Push(&nfa{Initial: &initial, Accept: frag.Accept})
		default:
			accept := State{}
			initial := State{symbol: CharacterClassToken{string(r)}, edge1: &accept}
			nfaStack.Push(&nfa{Initial: &initial, Accept: &accept})
		}
	}

	if nfaStack.Len() != 1 {
		panic("Nfa stack didn't have just a single element in it.")
	}
	res := nfaStack.Pop().(*nfa)
	return res
}

func IsEmpty(s *stack.Stack) bool {
	return s.Len() == 0
}

func Tokenize(infix string) []interface{} {

	var tokens []interface{}

	var s string
	var escapeStr string
	appendToS := false
	wantsToEscape := false
	for _, r := range infix { // handle character classes

		if wantsToEscape { // second backslash in a row, escape the character

			escapeStr += string(r)
			wantsToEscape = false
			switch r {
			case 'd': // \d
				tokens = append(tokens, DigitToken{val: string(r)})
			case 'w': // \w
				tokens = append(tokens, WordToken{val: string(r)})
			case '\\': // a second backslash
				tokens = append(tokens, CharacterClassToken{val: string(r)})
			}
			continue
		}
		if r == '\\' { // potentially want to escape a character
			wantsToEscape = true
			escapeStr = ""
			continue
		}

		// don't want to append the last element in a token
		if appendToS && r != ']' {
			s += string(r) // add as a single character of a multi character token
		} else if r != '[' && r != ']' { // add the single character as a token
			tokens = append(tokens, CharacterClassToken{string(r)})
		}

		if r == '[' { // we're going to start a multi character token
			appendToS = true
		} else if r == ']' { // reached end of character class
			tokens = append(tokens, CharacterClassToken{s}) // add the full string as a single token
			s = ""
			appendToS = false // stop building up character class
		}
	}

	return tokens
}

func InfixToPostfix(infix string) string {

	specials := map[string]int{"*": 10, ".": 9, "+": 8, "|": 7, "?": 6}

	postfix := stack.New()
	tempStack := stack.New()

	for _, r := range infix {
		switch {
		case r == '(':
			tempStack.Push(string(r))
		case r == ')':
			for tempStack.Peek() != "(" {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Pop()
		case specials[string(r)] > 0:
			for !IsEmpty(tempStack) && specials[string(r)] <= specials[tempStack.Peek().(string)] {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Push(string(r))
		default:
			postfix.Push(string(r))
		}
	}

	for !IsEmpty(tempStack) {
		postfix.Push(tempStack.Pop())
	}

	// return the elements as a more generic slice instead of a stack.
	// stack is an implementation detail of the algorithm
	var result []string
	for !IsEmpty(postfix) {
		result = append([]string{postfix.Pop().(string)}, result...)
	}
	return strings.Join(result, "")
}

// returns all the possible states from the given state, including
// state transitions with e arrows
//func addState(possibilities []*State, from, to *State) []*State {
//	seen := make(map[*State]bool) // keep track of the states we've already visited
//
//	states := stack.New()
//	states.Push(from)
//
//	for !IsEmpty(states) { // keep looking until dead end
//		next := states.Pop().(*State)
//
//		if seen[next] { // we may be looking at a state we've already seen
//			continue
//		}
//
//		// mark the state as having been seen already so we don't examine it again.
//		seen[next] = true
//		possibilities = append(possibilities, next) // this is a valid destination
//
//		// e arrow if symbol == 0
//		if next != to && next.symbol == 0 {
//			states.Push(next.edge1)
//			if next.edge2 != nil {
//				states.Push(next.edge2)
//			}
//		}
//	}
//
//	return possibilities
//}

//func addState(states []*State, start, accept *State) []*State {
//	states = append(states, start)
//	if start != accept && start.symbol == 0 {
//		states = addState(states, start.edge1, accept)
//		if start.edge2 != nil {
//			states = addState(states, start.edge2, accept)
//		}
//	}
//	return states
//}
//
//func postfixMatch(postfix, matchString string) bool {
//	nfa := postFixToNfa(postfix)
//
//	var current []*State
//	current = addState(current, nfa.Initial, nfa.Accept)
//	var next []*State
//
//	for _, r := range matchString {
//		for _, curr := range current {
//			if curr.symbol == r {
//				next = addState(next[:], curr.edge1, nfa.Accept)
//			}
//		}
//		current, next = next, []*State{}
//	}
//
//	for _, curr := range current {
//		if curr == nfa.Accept {
//			return true
//		}
//	}
//
//	return false
//}
