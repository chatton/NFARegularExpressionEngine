package main

import (
	"../nfa"
	"fmt"
	//"fmt"
)

func main() {
	//n := nfa.Compile("b.(c+).(e?)")
	//fmt.Println(n.Matches("bccccccc"))
	//nfa.Tokenize(`\d+e`)
	//nfa.InfixToPostfix("a.(b|c)")
	n := nfa.Compile(`\w+.\d+`)
	fmt.Println(n.Matches("abvc123"))
}
