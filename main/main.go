package main

import (
	"../nfa"
	"fmt"
)

func main() {
	fmt.Println(nfa.MatchString("abc+", "abccccccc"))
}
