package stack

import "errors"

func New() *Stack {
	return &Stack{}
}

type Stack struct {
	data []string
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Size() int {
	return len(s.data)
}

func (s *Stack) Push(elem string) {
	s.data = append(s.data, elem)
}

func (s *Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("stack was empty")
	}
	lastIndex := len(s.data) - 1
	lastItem := s.data[lastIndex]
	s.data = remove(s.data, lastIndex)
	return lastItem, nil
}

// this SO post had an answer for how to remove an element from a slice while maintaining order.
// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang/37335777
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
