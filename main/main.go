package main

import (
	"bufio"

	"../nfa"
	"fmt"
	"github.com/fatih/color"
	"os"
)

const (
	FEATURES = "features"
	OPTIONS  = "options"
)

var input *bufio.Scanner

func main() {
	input = bufio.NewScanner(os.Stdin)
	finished := false
	for !finished {
		printFile(OPTIONS)
		input.Scan()
		choice := input.Text()

		switch choice {
		case "1":
			matchStringOption()
		case "2":
			countOption()
		case "3":
			printFile(FEATURES)
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
		color.Green("The regular expression: " + pattern + " matched " + matchString + "!")
	} else {
		color.Yellow("The regular expression: " + pattern + " did not match " + matchString + ".")
	}
}

func countOption() {

}

func printFile(fileName string) {
	lines := readLines("data/" + fileName + ".txt")
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
