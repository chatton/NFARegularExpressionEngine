package nfa

import (
	"../utils"
	"github.com/golang-collections/collections/set"   // used to keep track of seen states
	"github.com/golang-collections/collections/stack" // useful stack data structure used throughout the code
	"strings"
)

type nfa struct {
	initial, accept *state
	tokens          []token
}

// provide mechanism for creating compiled Nfas so you don't need to construct the nfa each time you use it.
func Compile(infix string) *nfa {
	tokens := infixToPostfix(infix)
	return tokensToNfa(tokens)
}

type state struct {
	symbol       interface{} // the symbol of each state should be an implementation of the token interface
	edge1, edge2 *state
}

func tokensToNfa(tokens []token) *nfa {
	nfaStack := stack.New()
	for _, tok := range tokens {
		switch tok.(token).Val() {
		case ".":
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*nfa)
			frag1 := nfaStack.Pop().(*nfa)
			// join them together
			frag1.accept.edge1 = frag2.initial
			// place the single fragment on top of the stack
			nfaStack.Push(&nfa{initial: frag1.initial, accept: frag2.accept})
		case "|":
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*nfa)
			frag1 := nfaStack.Pop().(*nfa)
			// create a new initial state which points at both fragments
			initial := state{edge1: frag1.initial, edge2: frag2.initial}
			accept := state{}
			// point both fragments at the new accept state
			frag1.accept.edge1 = &accept
			frag2.accept.edge1 = &accept

			// add the new fragment to the stack.
			nfaStack.Push(&nfa{initial: &initial, accept: &accept})
		case "*":
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			accept := state{}
			// create and edge pointing back at itself and one towards the new accept state
			initial := state{edge1: frag.initial, edge2: &accept}
			// the accept state points back at the initial state
			frag.accept.edge1 = frag.initial
			// the second edge loops around
			frag.accept.edge2 = &accept
			// add it to the stack
			nfaStack.Push(&nfa{initial: &initial, accept: &accept})
		case "+":
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			accept := state{}
			initial := state{edge1: frag.initial, edge2: &accept}

			frag.accept.edge1 = &initial

			nfaStack.Push(&nfa{initial: frag.initial, accept: &accept})
		case "?":
			// take a single element off of the stack
			frag := nfaStack.Pop().(*nfa)
			// create a new state that points to the existing item and also the accept state
			initial := state{edge1: frag.initial, edge2: frag.accept}
			// push the new Nfa onto the stack
			nfaStack.Push(&nfa{initial: &initial, accept: frag.accept})
		default:
			accept := state{}
			initial := state{edge1: &accept, symbol: tok}
			nfaStack.Push(&nfa{initial: &initial, accept: &accept})
		}
	}
	result := nfaStack.Pop().(*nfa)
	result.tokens = tokens
	return result
}

// Helper method as stacks don't have an IsEmpty method
func IsEmpty(s *stack.Stack) bool {
	return s.Len() == 0
}

// function to covert an infix string into a list of tokens (still in infix)
func tokenize(infix string) []interface{} {

	var tokens []interface{}

	var s string
	var escapeStr string
	appendToS := false
	wantsToEscape := false
	negate := false
	ignoreCase := len(infix) >= 4 && strings.HasPrefix(infix, "(?i)")
	// all an input of "(?i)" + regex which will make it case insensitive.
	// put the infix string to lower case to eliminate differences between cases.
	if ignoreCase {
		infix = infix[4:] // remove the (?i) on the string
	}

	for i, r := range infix {

		if wantsToEscape { // there was a backslash escape the character
			escapeStr += string(r)
			wantsToEscape = false
			switch r {
			case 'd': // \d
				tokens = append(tokens, digitToken{val: `\d`, negate: negate})
			case 'w': // \w
				tokens = append(tokens, wordToken{val: `\w`, negate: negate})
			case 's':
				tokens = append(tokens, spaceToken{val: `\s`, negate: negate})
			default: // it's an escaped character
				tokens = append(tokens, characterClassToken{val: string(r), negate: negate, caseInsensitive: ignoreCase})
			}
			negate = false

			// handle implicit concatenation
			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) && !isClosingBracket(infix[i+1]) {
				if !isOpeningBracket(r) && r != '|' {
					tokens = append(tokens, characterClassToken{val: ".", negate: false})
				}
			}

			continue
		}

		if r == '\\' { // want to escape a character
			wantsToEscape = true
			escapeStr = ""
			continue
		}

		// if we should invert the match
		if r == '^' {
			negate = true
			continue
		}

		// if it's an underscore (and isn't being escaped), add a token that will match any character
		if r == '_' {
			tokens = append(tokens, anyToken{val: "_", negate: negate})

			// handle implicit concatenation
			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) && !isClosingBracket(infix[i+1]) {
				if !isOpeningBracket(r) && r != '|' {
					tokens = append(tokens, characterClassToken{val: ".", negate: false})
				}
			}
			continue
		}

		startingClass, endingClass := r == '[', r == ']'

		// don't want to append the last element in a token
		if appendToS && !endingClass {
			s += string(r) // add as a single character of a multi character token
		} else if !startingClass && !endingClass { // add the single character as a token
			tokens = append(tokens, characterClassToken{string(r), negate, ignoreCase})

			// handle implicit concatenation
			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) && !isClosingBracket(infix[i+1]) {
				if !isOpeningBracket(r) && r != '|' {
					tokens = append(tokens, characterClassToken{val: ".", negate: false})
				}
			}
			negate = false
		}

		if startingClass { // we're going to start a multi character token
			appendToS = true
		} else if endingClass { // reached end of character class
			tokens = append(tokens, characterClassToken{val: s, negate: negate}) // add the full string as a single token
			negate = false

			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) {
				// only don't add implicit concat if the next character is an explicit operator
				tokens = append(tokens, characterClassToken{".", false, ignoreCase})
			}

			s = ""
			appendToS = false // stop building up character class

		}
	}

	return tokens
}

// helper functions to provide meaning to certain characters
func isExplicitOperator(r uint8) bool {
	return r == '+' || r == '|' || r == '*' || r == '?'
}

func isClosingBracket(r uint8) bool {
	return r == ']' || r == ')'
}

func isOpeningBracket(r int32) bool {
	return r == '(' || r == '['
}

func infixToPostfix(infix string) []token {

	specials := map[string]int{"*": 10, "+": 9, "|": 5, "?": 8, ".": 6}

	postfix := stack.New()
	tempStack := stack.New()

	tokens := tokenize(infix) // get the tokens from the infix string

	// apply the shunting yard algorithm to the list of tokens
	for _, tok := range tokens {
		val := tok.(token).Val()

		switch {
		case val == "(":
			tempStack.Push(tok)
		case val == ")":
			for tempStack.Peek().(token).Val() != "(" {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Pop()
		case specials[val] > 0:
			for !IsEmpty(tempStack) && specials[val] <= specials[tempStack.Peek().(token).Val()] {
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

	// the final list of tokens in postfix
	var result []token

	for !IsEmpty(postfix) {
		// insert every element at the start of the list.
		result = append([]token{postfix.Pop().(token)}, result...)
	}

	return result
}

// returns all the possible states from the given state, including
// state transitions with e arrows
func addState(possibilities []*state, from, to *state) []*state {
	seen := set.New() // keep track of the states we've already visited
	states := stack.New()
	states.Push(from)

	for !IsEmpty(states) { // keep looking until dead end
		next := states.Pop().(*state)

		if seen.Has(next) { // we may be looking at a state we've already seen
			continue
		}

		// mark the state as having been seen already so we don't examine it again.
		seen.Insert(next)

		possibilities = append(possibilities, next) // this is a valid destination

		// e arrow if symbol == nil
		// zero value for token is nil
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

	var current []*state
	current = addState(current, n.initial, n.accept)
	var next []*state

	for _, r := range matchString {
		for _, curr := range current {
			// each token type implements the logic to match the specific character
			if curr.symbol != nil && curr.symbol.(token).Matches(r) {
				next = addState(next[:], curr.edge1, n.accept)
			}
		}
		current, next = next, []*state{}
	}

	for _, curr := range current {
		if curr == n.accept { // the accept state is found
			return true
		}
	}

	return false
}

// returns the number of unique occurrences of subString in fullString
func Count(subString, fullString string) int {
	count := 0
	for _, s := range utils.AllSubstrings(fullString) {
		if MatchString(subString, s) { // need to reconstruct the nfa for each substring
			count++
		}
	}
	return count
}

// function for use when you don't need to keep the NFA for multiple uses.
// Constructs nfa and calls the Matches method.
func MatchString(infix, matchString string) bool {
	n := Compile(infix)
	return n.Matches(matchString)
}
