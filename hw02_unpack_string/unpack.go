package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	str := []rune(s)
	if err := validate(str); err != nil {
		return "", fmt.Errorf("validating string error: %w", err)
	}

	var buf string
	var res strings.Builder
	for i, elem := range str {
		if unicode.IsDigit(elem) {
			res.WriteString(strings.Repeat(buf, int(elem-'0')))
			continue
		}
		if len(str) > i+1 && unicode.IsDigit(str[i+1]) {
			buf = string(elem)
			continue
		}
		res.WriteRune(elem)
	}
	return res.String(), nil
}

func validate(str []rune) error {
	if unicode.IsDigit(str[0]) {
		return ErrInvalidString
	}

	var previousRuneIsDigit bool
	for _, char := range str {
		currentRuneIsDigit := unicode.IsDigit(char)
		if currentRuneIsDigit && previousRuneIsDigit {
			return ErrInvalidString
		}
		previousRuneIsDigit = currentRuneIsDigit
	}
	return nil
}
