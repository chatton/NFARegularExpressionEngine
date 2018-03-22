package main

import (
	"../nfa"
	"fmt"
)

func main() {
	fmt.Println(nfa.Count("\\d\\w", "1a2b3cc"))
}
