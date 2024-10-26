package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// This package contains an implementation of a simple regular expression engine
// It is not intended to be used in a production environment and was created as an exercise to become more familiar with go
//
// It supports the operators "^", "$", "?", "*" and "+", as well as the wildcard "."
// These symbols can be escaped by prepending a "\" in the given regex
// This package does not support sub-patterns or many other common features

// computeGreedyOperator returns a boolean value
// '*' and '+' are known as greedy operators. They try to match as much of stringToMatch as they can.
// This functions parses them in a way that gives chars after them a chance to be matched as well.
// It recursively parses the operator and preceding char, until there is no match.
// Once there is no match, e.g. stringToMatch has been consumed but regEx has chars after the greedy operator,
// the call immediately before that stops matching the greedy operator and parses the rest of regEx and stringToMatch
// recursively.
func computeGreedyOperator(regEx, stringToMatch string) bool {
	// If match, checks if stringToMatch has more chars.
	if len(stringToMatch) > 1 {
		// If it has more chars, recursively match repeated char until there is no match.
		// This also triggers in case of ".*" or ".+" trying to consume the whole string, while there are
		// chars after the greedy operator, giving the rest of regEx a chance to be matched.
		if !matchEqualLength(regEx, stringToMatch[1:]) {
			// If there is a mismatch, passes regEx after the greedy operator and rest of stringToMatch to matchEqualLength().
			// Returns its result.
			return matchEqualLength(regEx[2:], stringToMatch[1:])
		} else {
			// If there is no mismatch, returns true.
			return true
		}
	} else {
		// If it has no more chars, passes regEx after the greedy operator and empty string to matchEqualLength().
		// Returns its results.
		return matchEqualLength(regEx[2:], "")
	}
}

// matchSingleChar returns a boolean value.
// The returned value is true, if the char represented by regExChar is the '.' wildcard
// or the chars represented by regExChar and charFromString match.
func matchSingleChar(regExChar, charFromString rune) bool {
	// Checks if the regExChar and charFromString are the same OR regExChar is wildcard.
	// Returns result of check.
	return regExChar == charFromString || regExChar == '.'
}

// computeQuestionMarkOperator returns a boolean value.
// The '?' operator in a regex means that the preceding char needs to occur zero or one times in the string.
// Recursively parses rest of regEx and stringToMatch.
func computeQuestionMarkOperator(regEx, stringToMatch string) bool {
	// Checks if char preceding '?' does NOT match first char in stringToMatch.
	if !matchSingleChar([]rune(regEx)[0], []rune(stringToMatch)[0]) {
		// If no match, passes regEx after '?' and intact stringToMatch to matchEqualLength().
		// Returns its result.
		return matchEqualLength(regEx[2:], stringToMatch)
	} else {
		// If match, checks if stringToMatch has more chars
		if len(stringToMatch) > 1 {
			// If it has more chars, passes regEx after '?' and rest of stringToMatch to matchEqualLength().
			// Returns its result.
			return matchEqualLength(regEx[2:], stringToMatch[1:])
		} else {
			// If it has no more chars, passes regEx after '?' and empty string to matchEqualLength().
			// Returns its results.
			return matchEqualLength(regEx[2:], "")
		}
	}
}

// computeAsteriskOperator returns a boolean value.
// The '*' operator in regex means that the preceding char needs to occur zero or more times in the string.
// Recursively parses regEx and rest of stringToMatch until there is no match.
// Then recursively parses rest regEx and rest of stringToMatch.
func computeAsteriskOperator(regEx, stringToMatch string) bool {
	// Checks if char preceding '*' does NOT match first char in stringToMatch.
	if !matchSingleChar([]rune(regEx)[0], []rune(stringToMatch)[0]) {
		// If no match, passes regEx after '*' and intact stringToMatch to matchEqualLength().
		// Returns its result.
		return matchEqualLength(regEx[2:], stringToMatch)
	} else {
		// Recursively computes greedy operator '*'
		return computeGreedyOperator(regEx, stringToMatch)
	}
}

// computePlusOperator returns a boolean value.
// The '+' operator in regex means that the preceding char needs to occur one or more times in the string.
// Returns false if there is no match at all.
// Recursively parses regEx and rest of stringToMatch until there is no match.
// Then recursively parses rest regEx and rest of stringToMatch.
func computePlusOperator(regEx, stringToMatch string) bool {
	// Checks if char preceding '+' does match first char in stringToMatch.
	if matchSingleChar([]rune(regEx)[0], []rune(stringToMatch)[0]) {
		// Recursively computes greedy operator '+'
		return computeGreedyOperator(regEx, stringToMatch)
	} else {
		// If there is no match, returns false.
		return false
	}
}

// matchEqualLength returns a boolean value
// Matches a regex against a string recursively.
// Returns false if stringToMatch has been consumed, but regEx has not.
func matchEqualLength(regEx, stringToMatch string) bool {
	// Checks if regEx has been consumed.
	if len(regEx) == 0 {
		// If regEx has been completely consumed, (sub) string matches.
		// Returns true.
		return true
		// Checks if the regEx is at its end, its last char is '$' and stringToMatch has been consumed.
	} else if len(regEx) == 1 && regEx[:1] == "$" && len(stringToMatch) == 0 {
		// This means the regEx matched so far and is at its end.
		// Returns true.
		return true
		// Checks if stringToMatch has been consumed
	} else if len(stringToMatch) == 0 {
		// If stringToMatch has been consumed, but not regEx, they don't match.
		// Returns false.
		return false
		// Checks if first char in regEx at the moment is escape char '\'.
	} else if []rune(regEx)[0] == '\\' {
		// If first char in regEx is '\', check if escaped char in regEx matches first char in stringToMatch.
		if len(regEx) > 1 && matchSingleChar([]rune(regEx)[1], []rune(stringToMatch)[0]) {
			// If escaped char matches, continue parsing rest of regEx and stringToMatch
			return matchEqualLength(regEx[2:], stringToMatch[1:])
		} else {
			// If escaped char does not match first char in stringToMatch, returns false.
			return false
		}
		// Checks if remaining regEx is longer than 1 char.
	} else if len(regEx) > 1 {
		// If remaining regEx is longer than 1 char, checks if 2nd char is an operator.
		switch regEx[1:2] {
		// If 2nd char is "?", calls computeQuestionMarkOperator() to recursively parse regEx and stringToMatch.
		case "?":
			return computeQuestionMarkOperator(regEx, stringToMatch)
		// If 2nd char is "*", calls computeAsteriskOperator() to recursively parse regEx and stringToMatch.
		case "*":
			return computeAsteriskOperator(regEx, stringToMatch)
		// If 2nd char is "+", calls computePlusOperator() to recursively parse regEx and stringToMatch.
		case "+":
			return computePlusOperator(regEx, stringToMatch)
		// If 2nd char is not one of the above operators.
		default:
			// Checks if 1st char of regEx and stringToMatch match.
			if matchSingleChar([]rune(regEx)[0], []rune(stringToMatch)[0]) {
				// If they match, recursively pass the rest of regEx and stringToMatch.
				return matchEqualLength(regEx[1:], stringToMatch[1:])
			} else {
				// If 1st char of regEx and stringToMatch don't match, return false.
				return false
			}
		}
		// If this case is reached, there remains only one char in regEx.
	} else {
		// Checks if only char left in regEx matches first char in stringToMatch and returns result.
		return matchSingleChar([]rune(regEx)[0], []rune(stringToMatch)[0])
	}
}

// match returns a boolean value.
// Recursively matches a given the given regEx to the given stringToMatch.
func match(regEx, stringToMatch string) bool {
	// Checks if length if regEx is not 0 AND if the first char in regEx is "^".
	if len(regEx) != 0 && regEx[:1] == "^" {
		// If regEx starts with the "^" operator, recursively match regEx exactly against beginning of stringToMatch.
		return matchEqualLength(regEx[1:], stringToMatch)
		//	Recursively tries to match regEx and stringToMatch.
	} else if matchEqualLength(regEx, stringToMatch) {
		// Returns true if regEx could be completely matched to stringToMatch.
		return true
		//	If regEx and stringToMatch couldn't be matched, checks if stringToMatch hasn't been consumed.
	} else if len(stringToMatch) != 0 {
		// If stringToMatch hasn't been consumed, try to match the regEx against the rest of stringToMatch.
		// This ensures that regEx is still matched against later parts of stringToMatch, in case regEx
		// and stringToMatch differ in length.
		return match(regEx, stringToMatch[1:])
	} else {
		// Return false in the case that stringToMatch has been consumed.
		return false
	}
}

// main is entry point of program.
func main() {
	// Creates new buffered input scanner
	scanner := bufio.NewScanner(os.Stdin)

	// Continuously scans for user input
	for scanner.Scan() {
		// Gets user input
		input := scanner.Text()
		// Separates input into regEx and stringToMatch.
		regEx, stringToMatch, _ := strings.Cut(input, "|")
		// Tries to match regEx and stringToMatch and prints the result.
		fmt.Println(match(regEx, stringToMatch))
	}
}
