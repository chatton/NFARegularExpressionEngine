package nfa

import (
	"../utils"
	"github.com/golang-collections/collections/set"   // used to keep track of seen states
	"github.com/golang-collections/collections/stack" // useful stack data structure used throughout the code
	"strings"
)

type nfa struct {
	Initial, Accept *State
	tokens          []Token
}

// provide mechanism for creating compiled Nfas so you don't need to construct the nfa each time you use it.
func Compile(infix string) *nfa {
	tokens := InfixToPostfix(infix)
	return tokensToNfa(tokens)
}

type State struct {
	symbol       interface{} // Token
	edge1, edge2 *State
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
	result := nfaStack.Pop().(*nfa)
	result.tokens = tokens
	return result
}

// Helper method as stacks don't have an IsEmpty method
func IsEmpty(s *stack.Stack) bool {
	return s.Len() == 0
}

func Tokenize(infix string) []interface{} {

	var tokens []interface{}

	var s string
	var escapeStr string
	appendToS := false
	wantsToEscape := false
	negate := false

	// all an input of "(?i)" + regex which will make it case insensitive.
	// put the infix string to lower case to eliminate differences between cases.
	if ignoreCase := len(infix) >= 4 && strings.HasPrefix(infix, "(?i)"); ignoreCase {
		infix = strings.ToLower(infix[4:]) // ignore the (?i) on the string
	}

	for i, r := range infix {

		if wantsToEscape { // there was a backslash escape the character
			escapeStr += string(r)
			wantsToEscape = false
			switch r {
			case 'd': // \d
				tokens = append(tokens, DigitToken{val: `\d`, negate: negate})
			case 'w': // \w
				tokens = append(tokens, WordToken{val: `\w`, negate: negate})
			case 's':
				tokens = append(tokens, SpaceToken{val: `\s`, negate: negate})
			default: // it's an escaped character
				tokens = append(tokens, CharacterClassToken{val: string(r), negate: negate})
			}
			negate = false

			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) && !isClosingBracket(infix[i+1]) {
				if !isOpeningBracket(r) && r != '|' {
					tokens = append(tokens, CharacterClassToken{val: ".", negate: false})
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

		// if it's an underscore, add a token that will match any character
		if r == '_' {
			tokens = append(tokens, AnyToken{val: "_", negate: negate})
			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) && !isClosingBracket(infix[i+1]) {
				if !isOpeningBracket(r) && r != '|' {
					tokens = append(tokens, CharacterClassToken{val: ".", negate: false})
				}
			}
			continue
		}

		startingClass, endingClass := r == '[', r == ']'

		// don't want to append the last element in a token
		if appendToS && !endingClass {
			s += string(r) // add as a single character of a multi character token
		} else if !startingClass && !endingClass { // add the single character as a token
			tokens = append(tokens, CharacterClassToken{string(r), negate})
			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) && !isClosingBracket(infix[i+1]) {
				if !isOpeningBracket(r) && r != '|' {
					tokens = append(tokens, CharacterClassToken{val: ".", negate: false})
				}
			}
			negate = false
		}

		if startingClass { // we're going to start a multi character token
			appendToS = true
		} else if endingClass { // reached end of character class
			tokens = append(tokens, CharacterClassToken{val: s, negate: negate}) // add the full string as a single token
			negate = false

			atEnd := i == len(infix)-1
			if !atEnd && !isExplicitOperator(infix[i+1]) {
				// only don't add implicit concat if the next character is an explicit operator
				tokens = append(tokens, CharacterClassToken{".", false})
			}

			s = ""
			appendToS = false // stop building up character class

		}
	}

	return tokens
}

func isExplicitOperator(r uint8) bool {
	return r == '+' || r == '|' || r == '*' || r == '?'
}

func isClosingBracket(r uint8) bool {
	return r == ']' || r == ')'
}

func isOpeningBracket(r int32) bool {
	return r == '(' || r == '['
}

func InfixToPostfix(infix string) []Token {

	specials := map[string]int{"*": 10, "+": 9, "|": 5, "?": 8, ".": 6}

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

// returns all the possible states from the given state, including
// state transitions with e arrows
func addState(possibilities []*State, from, to *State) []*State {
	seen := set.New() // keep track of the states we've already visited
	states := stack.New()
	states.Push(from)

	for !IsEmpty(states) { // keep looking until dead end
		next := states.Pop().(*State)

		if seen.Has(next) { // we may be looking at a state we've already seen
			continue
		}

		// mark the state as having been seen already so we don't examine it again.
		seen.Insert(next)

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
			// each token type implements the logic to match the specific character
			if curr.symbol != nil && curr.symbol.(Token).Matches(r) {
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

// returns the number of unique occurrences of subString in fullString
func Count(subString, fullString string) int {
	count := 0
	for _, s := range utils.AllSubstrings(fullString) {
		if MatchString(subString, s) {
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
