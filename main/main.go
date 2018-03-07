package main

import (
	"../nfa"
	"fmt"
)

func main() {
	regex := "a.b.(c|d).u*"
	postfix := nfa.InfixToPostfix(regex)
	fmt.Println(postfix)
	n := nfa.PostFixToNfa(postfix)
	fmt.Println(n)
}
