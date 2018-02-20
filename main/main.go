package main

import "fmt"
import "../stack"

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

func main() {

	fmt.Println("Infix:    ", "a.b.c*")
	fmt.Println("Postfix:  ", infixToPostfix("a.b.c*"))

	fmt.Println("Infix:    ", "(a.(b|d))*")
	fmt.Println("Postfix:  ", infixToPostfix("(a.(b|d))*"))

	fmt.Println("Infix:    ", "a.(b|d).c*")
	fmt.Println("Postfix:  ", infixToPostfix("a.(b|d).c*"))

	fmt.Println("Infix:    ", "a.(b.b)+.c")
	fmt.Println("Postfix:  ", infixToPostfix("a.(b.b)+.c"))
}
