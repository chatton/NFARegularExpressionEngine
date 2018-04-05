## Introduction

My name is Cian Hatton

This repository holds a regex expression engine created by modelling a non-deterministic finite automaton (NFA).

This project was part of my Graph Theory module in my 3rd year Software Development course at GMIT.

## Setup Instructions

Before running this application, you must make sure to have [Go](https://golang.org/dl/) installed.

This application makes use of several additional modules which can be aquired using the `go get` command. Before using this command, you must ensure that the `GOPATH` is set up. 

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
otherwise, see [this page](https://github.com/golang/go/wiki/SettingGOPATH) for instructions on how to set everything up in your operating system.

if you're on a Ubuntu based system, you can run the commands


in order to set up your go path temporarily to test this application.


if you have [git](https://git-scm.com/downloads) installed, you can run the command.

```bash
git clone https://github.com/chatton/NFARegularExpressionEngine.git
```

to clone the repository. If you don't want to use git, you can click on the Clone or Download button and download the source code as a zip file.

Navigate to the new directory

```bash
cd NFARegularExpressionEngine
```

download the dependencies for the project.

```bash
go get github.com/golang-collections/collections/set github.com/golang-collections/collections/stack github.com/fatih/color
```

and run it.

```bash
go run main/main.go
```

If you want to also run 

```bash
go test .\tests\nfa_test.go
```

you should get an output similar to

```bash
ok      command-line-arguments  0.076s
```

indicating that all tests pass.

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
number2 := nfa.Count(`\d`, "123456abc78") // 8
```

## The Supported Regular Expression Lanuguage

- `\d`  Digit: Matches a digit between 0 and 9

```go
nfa.MatchString(`\d234`, "1234") // true
```

- `\w`  Word: Matches a letter in the alphabet - case sensitively

```go
nfa.MatchString(`hello \w\w\w\w\w`, "hello world") // true
```
- `\s`  Space: Matches a single space
```go
nfa.MatchString(`\w+\s\w+\s`, "hello friend ") // true
```
- `\` Escape: Use \ to escape any special characters
```go
nfa.MatchString("\\d", `\d`) // true
```
- **[abc1]** Character Class
```go
nfa.MatchString("[123][456]","15") // true
nfa.MatchString("[123][456]","12") // false
nfa.MatchString("[123]+[456]+","123123456654") // true
```
- `^` Negate: Negates the result.
```go
nfa.MatchString("^[abc]","d") // true
nfa.MatchString("^[1]","1") // false
```
- `_` Any Character: Matches exactly 1 character.

```go
nfa.MatchString("he_lo","hello") // true
nfa.MatchString("_1_2_3","a1b2v3") // true
```

- `+` One or More: Causes the regular expression to match one or more occurrences.

```go
nfa.MatchString(`\d+`,"1234") // true
nfa.MatchString(`\d+`,"") // false
nfa.MatchString(`1+2+1+`,"11112221111") // true
```

- `*` Zero or More: Causes the regular expression to match zero or more occurrences.
```go
nfa.MatchString(`\d*`,"1234") // true
nfa.MatchString(`\d*`,"") // true
nfa.MatchString(`1*2*1*`,"2221111") // true
```

- `?` One or Zero: Causes the regular expression to match exactly 0 or 1 occurrence.

```go
nfa.MatchString(`hello?`,"hell") // true
nfa.MatchString(`hello?`,"hello") // true
```

- `|` Or: Matches if either the LHS or RHS matches.
```go
nfa.MatchString(`world ([123]|[abc])`,"world 1") // true
nfa.MatchString(`world ([123]|[abc])`,"world b") // true
nfa.MatchString(`world ([123]|[abc])`,"world p") // false

```
- `(?i)` Ignore Case: Prefix any expression with (?i)

```go
nfa.MatchString(`(?i)HeLlo`,"hElLo") // true
```

## Limitations

This Regular Expression Engine cannot do capture groups, or perform any other tasks that would require backtracking. This is a [limitation of NFAs](http://useless-factor.blogspot.ie/2008/05/regexp-research.html).

There is not any validation of the regular expressions passed in. As a result, the program will likely crash given a syntactically incorrect expression (within the given language).

Validation needs to be added in order to ensure that an error is passed back instead of allowing the program to crash.

## Technologies used

In this project, I used the [Go Programming Language](https://golang.org/)

## References

This project was initially based on a series of video lectures created by the lecturer for this course [Ian Mcloughlin](https://github.com/ianmcloughlin), I then built upon them and added additional features.

During development I frequently consulted documentation for both the [Python RE module](https://docs.python.org/3/library/re.html) and the [Go Regexp package](https://golang.org/pkg/regexp/)

I consulted [this paper](https://swtch.com/~rsc/regexp/regexp1.html) written by Russ Cox in order to implement Thompson's construction algorithm.


