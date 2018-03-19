package main

import (
	"../nfa"
	"fmt"
)

func main() {
	n := nfa.Compile("b.(c+).(e?)")
	fmt.Println(n.Matches("bccccccc"))
}
