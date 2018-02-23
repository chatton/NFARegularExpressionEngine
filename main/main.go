package main

import (
	"../nfa"
	"fmt"
)

func main() {

	n := nfa.New("a.b.c*")
	fmt.Printf(n.Regex)
	//fmt.Println("Infix:    ", "a.b.c*")
	//fmt.Println("Postfix:  ", infixToPostfix("a.b.c*"))
	//
	//fmt.Println("Infix:    ", "(a.(b|d))*")
	//fmt.Println("Postfix:  ", infixToPostfix("(a.(b|d))*"))
	//
	//fmt.Println("Infix:    ", "a.(b|d).c*")
	//fmt.Println("Postfix:  ", infixToPostfix("a.(b|d).c*"))
	//
	//fmt.Println("Infix:    ", "a.(b.b)+.c")
	//fmt.Println("Postfix:  ", infixToPostfix("a.(b.b)+.c"))
}
