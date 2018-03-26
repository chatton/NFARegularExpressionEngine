## Introduction

My name is Cian Hatton

This repository holds a regex expression engine created by modelling a non-deterministic finite automaton (NFA).

This project was part of my Graph Theory module in my 3rd year Software Development course at GMIT.

## Setup Instructions

## Design Decisions

Initially each token in the regular expression was a single character. This restricted the flexibility of the regular expressions that would be possible to create.

In particular, character classes. In order to get around this limitation, I created a `Token` interface. And the `Nfa` `States` would hold a reference to a `Token` (or interface as a Go specific implementation detail).

I then created several `Token` implementations.

`WordToken`, `CharacterClassToken`, `DigitToken`, `SpaceToken` etc.

Each of these `Tokens` provides a method `Match(r rune) bool {...}`. This method returns true or false for if a given character should match that `Token`. For example, a `CharaceterClassToken` created with the string [123] should match true if the provided rune is 1, 2 or 3. And false for everything else.

Each 'Token' also has a negate property. This makes it trivial to implement the negate (^) character by simply toggling this value.

This added an additional layer of complexity when parsing the regular expression which is provided as a string to the `nfa.Compile` method.

I also chose to implement implicit concatenation. i.e. in order you write a regular expression that is the 1 followed by 2 followed by 3. You can write 123 instead of 1.2.3 


## Features and Functionality

There are two main features. That is

1. `MatchString` This has the method signature 'MatchString(infix, matchString string) bool' And is used when you want to do a one off check for a string against a regular expression. It is also possible to create an `nfa` with the `nfa.Compile` method if you want to re-use the same nfa on multiple strings.

Single use
```go
isMatch := nfa.MatchString(`\s+(\d|[abc])`, "    b")
fmt.Println(isMatch) // true
```

Re-using an nfa
```go
nfa := nfa.Compile(`\d+\w+\d+`)
fmt.Println(nfa.Matches("123abc456")) // true
fmt.Println(nfa.Matches("123456")) // false
```

2. `Count` This method lets you find the number of occurrences of a regular expression in a given string. For example

```go
number := nfa.Count("ab", "ababab") // 3
number2 := nfa.Count(`\d`, "123456abc78) // 8
```

## The Supported Regular Expression Lanuguage

- **\d**  Digit: Matches a digit between 0 and 9
- **\w**  Word: Matches a letter in the alphabet - case sensitively
- **\s**  Space: Matches a single space
- **\** Escape: Use \ to escape any special characters. E.g. \\d\\w will match the string "\d\w"
- **[abc1]** Character Class: Matches one of a, b, c or 1
- **^** Negate: Negates the result .e.g '^\d' will match a single character that isn't a digit
- **_** Any Character: Matches exactly 1 character. E.g \d_\d will match 1k8
- **\+** One or More: Causes the regular expression to match one or more occurrences. E.g 1+ will match one or more 1s
- \* Zero or More: Causes the regular expression to match zero or more occurrences. E.g. H* will match zero or more H characters
- **?** One or Zero: Causes the regular expression to match exactly 0 or 1 occurrence. E.g. hel?o will match "heo" and "helo"
- **|** Or: Matches if either the LHS or RHS matches. E.g. ([\d]|j) will match any single digit or the letter j
- **(?i)** Ignore Case: Prefix any expression with (?i) and the case will be ignored. E.g. (?i)HeLlo will match hELlO

## Limitations

This Regular Expression Engine cannot do capture groups, or perform any other tasks that would require backtracking. This is a [limitation of NFAs](http://useless-factor.blogspot.ie/2008/05/regexp-research.html).

## Technologies used

In this project, I used the [Go Programming Language](https://golang.org/)

## References

During development I frequently consulted documentation for both the [Python RE module](https://docs.python.org/3/library/re.html) and the [Go Regexp package](https://golang.org/pkg/regexp/)

I consulted [this paper](https://swtch.com/~rsc/regexp/regexp1.html) written by Russ Cox in order to implement Thompson's construction algorithm.


