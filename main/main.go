package main

import "fmt"
import "../stack"

func infixToPostfix(infix string) string {

	specials := map[rune]int{'*': 10, '.': 9, '|': 8}
	postfix := stack.New()
	stack := stack.New()

	for _, r := range infix {
		switch {
		case r == '(':
			stack.Push(r)
		case r == ')':
			for stack.Peek() != '(' {
				postfix.Push(stack.Pop())
			}
			stack.Pop()

		case specials[r] > 0:
			for !stack.IsEmpty() && specials[r] <= specials[stack.Peek()] {
				postfix.Push(stack.Pop())
			}
			stack.Push(r)
		default:
			postfix.Push(r)
		}
	}

	for !stack.IsEmpty() {
		postfix.Push(stack.Pop())
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
