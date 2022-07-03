package common

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/lexer"
	"math"
	"strconv"
	"strings"
)

/*
ResolveFloat returns a new token with simplified float notation of any float
will panic on invalid token
*/
func ResolveFloat(token *lexer.Token) *lexer.Token {
	s := strings.ReplaceAll(token.String(), "_", "")
	f, parseError := strconv.ParseFloat(s, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%f", f)),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        token.Line,
		Column:      token.Column,
		Index:       token.Index,
	}
}

func FloatGreaterToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue > rValue
}

func FloatGreatOrEqualToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue >= rValue
}

func FloatLessToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue < rValue
}

func FloatLessOrEqualToken(left, right *lexer.Token) bool {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return lValue <= rValue
}

func FloatAddToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%f", lValue+rValue)),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func FloatSubToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%f", lValue-rValue)),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func FloatMulToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%f", lValue*rValue)),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func FloatDivToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%f", lValue/rValue)),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func FloatFloorDivToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%d", int64(lValue/rValue))),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}

func FloatPowToken(left, right *lexer.Token) *lexer.Token {
	var (
		lValue, rValue float64
		parseError     error
	)
	lValue, parseError = strconv.ParseFloat(left.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	rValue, parseError = strconv.ParseFloat(right.String(), 64)
	if parseError != nil {
		panic(parseError)
	}
	return &lexer.Token{
		Contents:    []rune(fmt.Sprintf("%f", math.Pow(lValue, rValue))),
		DirectValue: lexer.Float,
		Kind:        lexer.Literal,
		Line:        left.Line,
		Column:      left.Column,
		Index:       left.Index,
	}
}
