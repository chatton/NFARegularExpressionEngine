package stack

import "../nfa"

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
	data []*nfa.Nfa
}

func (s *NfaStack) Data() []*nfa.Nfa {
	tmp := make([]*nfa.Nfa, len(s.data))
	copy(tmp, s.data)
	return tmp
}

func (s *NfaStack) Peek() *nfa.Nfa {
	return s.data[len(s.data)-1]
}

func (s *NfaStack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *NfaStack) Size() int {
	return len(s.data)
}

func (s *NfaStack) Push(elem *nfa.Nfa) {
	s.data = append(s.data, elem)
}

func (s *NfaStack) Pop() *nfa.Nfa {
	lastIndex := len(s.data) - 1
	lastItem := s.data[lastIndex]
	s.data = removeNfa(s.data, lastIndex)
	return lastItem
}

func removeNfa(slice []*nfa.Nfa, s int) []*nfa.Nfa {
	return append(slice[:s], slice[s+1:]...)
}
