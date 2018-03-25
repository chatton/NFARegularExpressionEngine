## Introduction

My name is Cian Hatton

This repository holds a regex expression engine created by modelling a non-deterministic finite automaton.

This project was part of my Graph Theory module in my 3rd year Software Development course at GMIT.

## Setup Instructions

## Design Decisions

## Problems Encountered

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

## Technologies used

## References
