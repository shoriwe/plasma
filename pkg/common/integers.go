package common

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/lexer"
	"strconv"
	"strings"
)

/*
ResolveInteger returns a new token with simplified decimal notation of any integer
will panic on invalid token
*/
func ResolveInteger(token *lexer.Token) *lexer.Token {
	s := strings.ReplaceAll(token.String(), "_", "")
	i, parseError := strconv.ParseInt(s, 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", i)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        token.Line,
		Column:      token.Column,
		Index:       token.Index,
	}
}

func IntegerGreaterToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue > rValue
}

func IntegerGreatOrEqualToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue >= rValue
}

func IntegerLessToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue < rValue
}

func IntegerLessOrEqualToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue <= rValue
}

func IntegerBitwiseOrToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue|rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerBitwiseXorToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue^rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerBitwiseAndToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue&rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerBitwiseLeftToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue<<rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerBitwiseRightToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue>>rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerAddToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue+rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerSubToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue-rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerMulToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue*rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerDivToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue/rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerModToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue int64
		parseError     error
	)
	lValue, parseError = strconv.ParseInt(left.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseInt(right.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", lValue%rValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func IntegerNegateBitsToken(x *lexer.Token) *lexer.Token {
	var (
		xValue     int64
		parseError error
	)
	xValue, parseError = strconv.ParseInt(x.String(), 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", ^xValue)),
		DirectValue: lexer.Integer,
		Kind:        lexer.Literal,
		Line:        x.Line,
		Column:      x.Column,
		Index:       x.Index,
	}
}
