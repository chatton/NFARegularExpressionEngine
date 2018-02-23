package nfa

import "../stack"

type Nfa struct {
	Regex string
}

func New(regex string) *Nfa {
	postFix := infixToPostfix(regex)
	return &Nfa{postFix}
}

func infixToPostfix(infix string) string {

	specials := map[rune]int{'*': 10, '.': 9, '|': 8}
	postfix := stack.New()
	tempStack := stack.New()

	for _, r := range infix {
		switch {
		case r == '(':
			tempStack.Push(r)
		case r == ')':
			for tempStack.Peek() != '(' {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Pop()

		case specials[r] > 0:
			for !tempStack.IsEmpty() && specials[r] <= specials[tempStack.Peek()] {
				postfix.Push(tempStack.Pop())
			}
			tempStack.Push(r)
		default:
			postfix.Push(r)
		}
	}

	for !tempStack.IsEmpty() {
		postfix.Push(tempStack.Pop())
	}

	result := postfix.Data()
	return string(result)
}
