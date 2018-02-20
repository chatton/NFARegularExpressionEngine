package main

import "fmt"
import "../stack"

func infixToPostfix(regex string) string {
	return ""
}

func main() {

	s := stack.New()
	s.Push("1")
	s.Push("2")
	s.Push("3")
	s.Push("4")
	s.Push("5")

	for i := 0; i < 5; i++ {
		fmt.Println(s.Pop())
	}
}
