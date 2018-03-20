package main

import (
	"../nfa"
	"fmt"
)

func main() {
	n := nfa.Compile("\\w+.\\d+")
	fmt.Println(n.Matches("hello"))
}
