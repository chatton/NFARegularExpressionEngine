package stack

func New() *Stack {
	return &Stack{}
}

type Stack struct {
	data []rune
}

func (s *Stack) Data() []rune {
	tmp := make([]rune, len(s.data))
	copy(tmp, s.data)
	return tmp
}

func (s *Stack) Peek() rune {
	return s.data[len(s.data)-1]
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Size() int {
	return len(s.data)
}

func (s *Stack) Push(elem rune) {
	s.data = append(s.data, elem)
}

func (s *Stack) Pop() rune {

	lastIndex := len(s.data) - 1
	lastItem := s.data[lastIndex]
	s.data = remove(s.data, lastIndex)
	return lastItem
}

// this SO post had an answer for how to remove an element from a slice while maintaining order.
// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-array-in-golang/37335777
func remove(slice []rune, s int) []rune {
	return append(slice[:s], slice[s+1:]...)
}
