package common

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/lexer"
	"strconv"
	"strings"
)

var directCharEscapeValue = map[rune][]rune{
	'a': {7},
	'b': {8},
	'e': {'\\', 'e'},
	'f': {12},
	'n': {10},
	'r': {13},
	't': {9},
	'?': {'\\', '?'},
}

func Repeat(s string, times int64) string {
	result := ""
	for i := int64(0); i < times; i++ {
		result += s
	}
	return result
}

func ReplaceEscaped(s []rune) []rune {
	sLength := len(s)
	escaped := false
	var result []rune
	for index := 0; index < sLength; index++ {
		char := s[index]
		if escaped {
			switch char {
			case 'a', 'b', 'e', 'f', 'n', 'r', 't', '?':
				// Replace char based
				result = append(result, directCharEscapeValue[char]...)
			case '\\', '\'', '"', '`':
				// Replace escaped literals
				result = append(result, char)
			case 'x':
				// Replace hex with numbers
				index++
				a := s[index]
				index++
				b := s[index]
				number, parsingError := strconv.ParseUint(string([]rune{a, b}), 16, 32)
				if parsingError != nil {
					panic(parsingError)
				}
				result = append(result, rune(number))
			case 'u':
				// Replace unicode with numbers
				index++
				a := s[index]
				index++
				b := s[index]
				index++
				c := s[index]
				index++
				d := s[index]
				number, parsingError := strconv.ParseUint(string([]rune{a, b, c, d}), 16, 32)
				if parsingError != nil {
					panic(parsingError)
				}
				result = append(result, rune(number))
			}
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else {
			result = append(result, char)
		}
	}
	return result
}

func ResolveString(token *lexer.Token) *lexer.Token {
	return &lexer.Token{
		Contents: []rune(
			fmt.Sprintf("'%s'", string(token.Contents[1:len(token.Contents)-1])),
		),
		DirectValue: lexer.SingleQuoteString,
		Kind:        lexer.Literal,
		Line:        token.Line,
		Column:      token.Column,
		Index:       token.Index,
	}
}

func ResolveBytesString(token *lexer.Token) *lexer.Token {
	return &lexer.Token{
		Contents: []rune(
			fmt.Sprintf("b'%s'", string(token.Contents[2:len(token.Contents)-1])),
		),
		DirectValue: lexer.SingleQuoteString,
		Kind:        lexer.Literal,
		Line:        token.Line,
		Column:      token.Column,
		Index:       token.Index,
	}
}

func StringAddToken(left, right *lexer.Token) *lexer.Token {
	return &lexer.Token{
		Contents:    append(left.Contents[1:len(left.Contents)-1], right.Contents[1:len(right.Contents)-1]...),
		DirectValue: lexer.SingleQuoteString,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func StringMulToken(left, right *lexer.Token) *lexer.Token {
	var (
		rValue     int64
		parseError error
	)
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(strings.Repeat(string(left.Contents[1:len(left.Contents)-1]), int(rValue))),
		DirectValue: lexer.SingleQuoteString,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func BytesStringAddToken(left, right *lexer.Token) *lexer.Token {
	return &lexer.Token{
		Contents:    append(left.Contents[2:len(left.Contents)-1], right.Contents[2:len(right.Contents)-1]...),
		DirectValue: lexer.ByteString,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func BytesStringMulToken(left, right *lexer.Token) *lexer.Token {
	var (
		rValue     int64
		parseError error
	)
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(strings.Repeat(string(left.Contents[2:len(left.Contents)-1]), int(rValue))),
		DirectValue: lexer.SingleQuoteString,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}
