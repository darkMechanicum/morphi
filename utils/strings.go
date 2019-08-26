package utils

import (
	"errors"
	"text/scanner"
)

// Get first rune of passed string and return error if string is empty.
func GetFirstRune(self *string) (rune, error) {
	for _, c := range *self {
		return c, nil
	}
	return scanner.EOF, errors.New("Can't get first character from empty string.")
}

// Get first rune of passed string and return error if string is empty.
// Also, check it is the only rune in the string and return error is there are more.
func GetOnlyRune(self *string) (rune, error) {
	if len(*self) > 1 {
		return scanner.EOF, errors.New("String contains more than one character.")
	}
	return GetFirstRune(self)
}
