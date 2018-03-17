package nfa

import (
	"github.com/golang-collections/collections/stack" // useful stack data structure
	"strings"
)

//type Nfa struct {
//	Initial, Accept *State
//}
//
//func (n Nfa) String() string {
//	return n.Initial.String()
//}
//
//type Token struct {
//	Value string
//}
//
//func (t Token) String() string {
//	return t.Value
//}
//
//type State struct {
//	Token        *Token
//	edge1, edge2 *State
//}
//
//func (s State) String() string {
//	result := ""
//	if s.Token == nil {
//		result += "Ɛ"
//	} else {
//		result += s.Token.Value
//	}
//	if s.edge1 != nil {
//		result += s.edge1.String()
//	}
//
//	if s.edge2 != nil {
//		result += s.edge2.String()
//	}
//
//	return result + " "
//}
//
//func contains(allStrings []string, substr string) bool {
//	for _, str := range allStrings {
//		if strings.Contains(str, substr) {
//			return true
//		}
//
//	}
//	return false
//}
//
//func splitBy(s string, operators []string) []string {
//	allStrings := []string{}
//
//	//for _, r := range s {
//	//
//	//}
//
//	return allStrings
//}
//
//func Tokenize(strings []string) []Token {
//	tokens := []Token{}
//
//	return tokens
//}
//
//func Pop2(nfaStack *NfaStack) (*Nfa, *Nfa) {
//	popped1 := nfaStack.Pop()
//	popped2 := nfaStack.Pop()
//	return popped2, popped1
//}
//
//func PostFixToNfa(postfix string) *Nfa {
//
//	nfaStack := NewNfaStack()
//	for _, r := range postfix {
//		switch r {
//		case '.':
//			// take 2 elements off the stack
//			frag1, frag2 := Pop2(nfaStack)
//			// join them together
//			frag1.Accept.edge1 = frag2.Initial
//			// place the single fragment on top of the stack
//			nfaStack.Push(&Nfa{Initial: frag1.Initial, Accept: frag2.Accept})
//		case '|':
//			// take 2 elements off the stack
//			frag1, frag2 := Pop2(nfaStack)
//
//			// create a new initial state which points at both fragments
//			initial := State{edge1: frag1.Initial, edge2: frag2.Initial}
//			accept := State{}
//			// point both fragments at the new accept state
//			frag1.Accept = &accept
//			frag2.Accept = &accept
//
//			// add the new fragment to the stack.
//			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
//		case '*':
//			// take a single element off of the stack
//			frag := nfaStack.Pop()
//			accept := State{}
//			// create and edge pointing back at itself and one towards the new accept state
//			initial := State{edge1: frag.Initial, edge2: &accept}
//			// the accept state points back at the initial state
//			frag.Accept.edge1 = frag.Initial
//			// the second edge loops around
//			frag.Accept.edge2 = frag.Accept
//			// add it to the stack
//			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
//		case '+':
//			fallthrough
//		default:
//			accept := State{}
//			initial := State{Token: &Token{string(r)}, edge1: &accept}
//			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
//		}
//	}
//
//	// TODO handle incorrect amount
//
//	return nfaStack.Pop() // the last element
//}

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
