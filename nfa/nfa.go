package nfa

import (
	"github.com/golang-collections/collections/stack" // useful stack data structure
	"strings"
)

type Nfa struct {
	Initial, Accept *State
}

type State struct {
	symbol       rune
	edge1, edge2 *State
}

func PostFixToNfa(postfix string) *Nfa {

	nfaStack := stack.New()
	for _, r := range postfix {
		switch r {
		case '.':
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*Nfa)
			frag1 := nfaStack.Pop().(*Nfa)
			// join them together
			frag1.Accept.edge1 = frag2.Initial
			// place the single fragment on top of the stack
			nfaStack.Push(&Nfa{Initial: frag1.Initial, Accept: frag2.Accept})
		case '|':
			// take 2 elements off the stack
			frag2 := nfaStack.Pop().(*Nfa)
			frag1 := nfaStack.Pop().(*Nfa)
			// create a new initial state which points at both fragments
			initial := State{edge1: frag1.Initial, edge2: frag2.Initial}
			accept := State{}
			// point both fragments at the new accept state
			frag1.Accept.edge1 = &accept
			frag2.Accept.edge1 = &accept

			// add the new fragment to the stack.
			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
		case '*':
			// take a single element off of the stack
			frag := nfaStack.Pop().(*Nfa)
			accept := State{}
			// create and edge pointing back at itself and one towards the new accept state
			initial := State{edge1: frag.Initial, edge2: &accept}
			// the accept state points back at the initial state
			frag.Accept.edge1 = frag.Initial
			// the second edge loops around
			frag.Accept.edge2 = &accept
			// add it to the stack
			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
		case '+':
			fallthrough
		default:
			accept := State{}
			initial := State{symbol: r, edge1: &accept}
			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
		}
	}

	if nfaStack.Len() != 1 {
		panic("Nfa stack didn't have just a single element in it.")
	}
	res := nfaStack.Pop().(*Nfa)
	return res
}

func IsEmpty(s *stack.Stack) bool {
	return s.Len() == 0
}

func InfixToPostfix(infix string) string {

	specials := map[string]int{"*": 10, ".": 9, "|": 8}

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
//func getAllPossibleStates(from, to *State) []*State {
//	seen := make(map[*State]bool) // keep track of the states we've already visited
//
//	states := stack.New()
//	states.Push(from)
//
//	var possibilities []*State
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
//
//	}
//	fmt.Println(len(possibilities))
//	return possibilities
//}

func addState(states []*State, start, accept *State) []*State {
	states = append(states, start)
	if start != accept && start.symbol == 0 {
		states = addState(states, start.edge1, accept)
		if start.edge2 != nil {
			states = addState(states, start.edge2, accept)
		}
	}
	return states
}

func PostfixMatch(postfix, matchString string) bool {
	nfa := PostFixToNfa(postfix)

	var current []*State
	current = addState(current, nfa.Initial, nfa.Accept)
	var next []*State

	for _, r := range matchString {
		for _, curr := range current {
			if curr.symbol == r {
				next = addState(next[:], curr.edge1, nfa.Accept)
			}
		}
		current, next = next, []*State{}

	}

	for _, curr := range current {
		if curr == nfa.Accept {
			return true
		}
	}

	return false
}