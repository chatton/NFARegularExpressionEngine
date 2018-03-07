package nfa

import (
	"strings"
)

type Nfa struct {
	Initial, Accept *State
}

func (n Nfa) String() string {
	return n.Initial.String()
}

type Token struct {
	Value string
}

func (t Token) String() string {
	return t.Value
}

type State struct {
	Token        *Token
	edge1, edge2 *State
}

func (s State) String() string {
	result := ""
	if s.Token == nil {
		result += "Ɛ"
	} else {
		result += s.Token.Value
	}
	if s.edge1 != nil {
		result += s.edge1.String()
	}

	if s.edge2 != nil {
		result += s.edge2.String()
	}

	return result + " "
}

func contains(allStrings []string, substr string) bool {
	for _, str := range allStrings {
		if strings.Contains(str, substr) {
			return true
		}

	}
	return false
}

func splitBy(s string, operators []string) []string {
	allStrings := []string{}

	//for _, r := range s {
	//
	//}

	return allStrings
}

func Tokenize(strings []string) []Token {
	tokens := []Token{}

	return tokens
}

func Pop2(nfaStack *NfaStack) (*Nfa, *Nfa) {
	popped1 := nfaStack.Pop()
	popped2 := nfaStack.Pop()
	return popped2, popped1
}

func PostFixToNfa(postfix string) *Nfa {

	nfaStack := NewNfaStack()
	for _, r := range postfix {
		switch r {
		case '.':
			// take 2 elements off the stack
			frag1, frag2 := Pop2(nfaStack)
			// join them together
			frag1.Accept.edge1 = frag2.Initial
			// place the single fragment on top of the stack
			nfaStack.Push(&Nfa{Initial: frag1.Initial, Accept: frag2.Accept})
		case '|':
			// take 2 elements off the stack
			frag1, frag2 := Pop2(nfaStack)

			// create a new initial state which points at both fragments
			initial := State{edge1: frag1.Initial, edge2: frag2.Initial}
			accept := State{}
			// point both fragments at the new accept state
			frag1.Accept = &accept
			frag2.Accept = &accept

			// add the new fragment to the stack.
			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
		case '*':
			// take a single element off of the stack
			frag := nfaStack.Pop()
			accept := State{}
			// create and edge pointing back at itself and one towards the new accept state
			initial := State{edge1: frag.Initial, edge2: &accept}
			// the accept state points back at the initial state
			frag.Accept.edge1 = frag.Initial
			// the second edge loops around
			frag.Accept.edge2 = frag.Accept
			// add it to the stack
			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
		case '+':
			fallthrough
		default:
			accept := State{}
			initial := State{Token: &Token{string(r)}, edge1: &accept}
			nfaStack.Push(&Nfa{Initial: &initial, Accept: &accept})
		}
	}

	// TODO handle incorrect amount

	return nfaStack.Pop() // the last element
}

func InfixToPostfix(infix string) string {

	specials := map[string]int{"*": 10, ".": 9, "|": 8}

	postfix := NewStringStack()
	tempStack := NewStringStack()

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
			for !tempStack.IsEmpty() && specials[string(r)] <= specials[tempStack.Peek()] {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Push(string(r))
		default:
			postfix.Push(string(r))
		}
	}

	for !tempStack.IsEmpty() {
		postfix.Push(tempStack.Pop())
	}

	result := postfix.Data()
	return strings.Join(result, "")
}

func NewStringStack() *StringStack {
	return &StringStack{}
}

type StringStack struct {
	data []string
}

func (s *StringStack) Data() []string {
	tmp := make([]string, len(s.data))
	copy(tmp, s.data)
	return tmp
}

func (s *StringStack) Peek() string {
	return s.data[len(s.data)-1]
}

func (s *StringStack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *StringStack) Size() int {
	return len(s.data)
}

func (s *StringStack) Push(elem string) {
	s.data = append(s.data, elem)
}

func (s *StringStack) Pop() string {
	lastIndex := len(s.data) - 1
	lastItem := s.data[lastIndex]
	s.data = removeString(s.data, lastIndex)
	return lastItem
}

// this SO post had an answer for how to remove an element from a slice while maintaining order.
// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang/37335777
func removeString(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// NFA stack.

func NewNfaStack() *NfaStack {
	return &NfaStack{}
}

type NfaStack struct {
	data []*Nfa
}

func (s *NfaStack) Data() []*Nfa {
	tmp := make([]*Nfa, len(s.data))
	copy(tmp, s.data)
	return tmp
}

func (s *NfaStack) Peek() *Nfa {
	return s.data[len(s.data)-1]
}

func (s *NfaStack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *NfaStack) Size() int {
	return len(s.data)
}

func (s *NfaStack) Push(elem *Nfa) {
	s.data = append(s.data, elem)
}

func (s *NfaStack) Pop() *Nfa {
	lastIndex := len(s.data) - 1
	lastItem := s.data[lastIndex]
	s.data = removeNfa(s.data, lastIndex)
	return lastItem
}

func removeNfa(slice []*Nfa, s int) []*Nfa {
	return append(slice[:s], slice[s+1:]...)
}
