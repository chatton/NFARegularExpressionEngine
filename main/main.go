package main

import (
	"bufio"
	"strconv"

	"../nfa"
	"fmt"
	"os"
)

const (
	FEATURES = "features"
	OPTIONS  = "options"
)

// allow scanner to be seen and used as the global user input
var input *bufio.Scanner
var finished bool

// init variables before main runs
func init() {
	input = bufio.NewScanner(os.Stdin)
	finished = false
}

func main() {
	for !finished {
		printFile(OPTIONS, "txt")
		input.Scan()
		choice := input.Text()

		switch choice {
		case "1":
			matchStringOption()
		case "2":
			countOption()
		case "3":
			printFile(FEATURES, "txt")
			fmt.Println()
			fmt.Println()
		case "4":
			fmt.Println("Exiting application...")
			finished = true
			break
		default:
			fmt.Println("Enter a valid option.")
			continue
		}
	}
}

func matchStringOption() {
	fmt.Print("Enter regular expression: ")
	input.Scan()
	pattern := input.Text()
	fmt.Print("Enter a string to match: ")
	input.Scan()
	matchString := input.Text()

	result := nfa.MatchString(pattern, matchString)

	if matchString == "" {
		matchString = "the empty string"
	}

	if result {
		fmt.Println("The regular expression:", pattern, "matched", matchString+"!")
	} else {
		fmt.Println("The regular expression:", pattern, "did not match", matchString+".")
	}
}

func countOption() {
	fmt.Print("Enter regular expression: ")
	input.Scan()
	pattern := input.Text()
	fmt.Print("Enter a string to count number of occurrences: ")
	input.Scan()
	searchString := input.Text()
	num := nfa.Count(pattern, searchString)
	numStr := strconv.Itoa(num)
	fmt.Println("The pattern", pattern, "occurred", numStr, "times in the string", searchString)
}

// helper function to print all the lines in a file
func printFile(fileName, ext string) {
	lines := readLines("data/" + fileName + "." + ext)
	for _, line := range lines {
		fmt.Println(line)
	}
}

func readLines(path string) []string {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to read from " + path)
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
