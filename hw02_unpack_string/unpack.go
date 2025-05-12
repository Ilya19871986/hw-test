package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var ErrInternal = errors.New("internal error")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if unicode.IsDigit(rune(s[0])) {
		return "", ErrInvalidString
	}

	var builder strings.Builder
	var last rune

	for i, current := range s {
		if unicode.IsLetter(current) || unicode.IsSpace(current) {
			last = current
			builder.WriteRune(current)
		} else if unicode.IsDigit(current) {
			if unicode.IsDigit(rune(s[i-1])) {
				return "", ErrInvalidString
			}
			count, err := strconv.Atoi(string(current))
			if err != nil || count < 0 {
				return "", ErrInternal
			}
			if count == 0 {
				result := builder.String()
				builder.Reset()
				builder.WriteString(result[:len(result)-1])
			} else {
				builder.WriteString(strings.Repeat(string(last), count-1))
			}
		}
	}
	return builder.String(), nil
}
