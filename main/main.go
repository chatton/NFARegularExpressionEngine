package main

import (
	"../nfa"
	"fmt"
)

func main() {
	n := nfa.Compile("[def]")
	fmt.Println(n.Matches("d"))
}
