package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
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

func (simp *simplify) simplifyInteger(s string) *ast2.Integer {
	s = strings.ReplaceAll(strings.ToLower(s), "_", "")
	value, parseError := strconv.ParseInt(s, 0, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &ast2.Integer{
		Value: value,
	}
}

func (simp *simplify) simplifyFloat(s string) *ast2.Float {
	s = strings.ReplaceAll(strings.ToLower(s), "_", "")
	value, parseError := strconv.ParseFloat(s, 64)
	if parseError != nil {
		panic(parseError)
	}
	return &ast2.Float{
		Value: value,
	}
}

func (simp *simplify) simplifyString(rawS string) *ast2.String {
	s := []rune(rawS)
	s = s[1 : len(s)-1]
	sLength := len(s)
	escaped := false
	resolved := make([]rune, 0, len(s))
	for index := 0; index < sLength; index++ {
		char := s[index]
		if escaped {
			switch char {
			case 'a', 'b', 'e', 'f', 'n', 'r', 't', '?':
				// Replace char based
				resolved = append(resolved, directCharEscapeValue[char]...)
			case '\\', '\'', '"', '`':
				// Replace escaped literals
				resolved = append(resolved, char)
			case 'x':
				// Replace hex with numbers
				index++
				a := s[index]
				index++
				b := s[index]
				resolved = append(resolved, a*16+b)
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
				resolved = append(resolved, a*4096+b*256+c*16+d)
			}
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else {
			resolved = append(resolved, char)
		}
	}
	return &ast2.String{
		Contents: []byte(string(resolved)),
	}
}

func (simp *simplify) simplifyBytes(rawS string) *ast2.Bytes {
	s := []rune(rawS)
	s = s[2 : len(s)-1]
	sLength := len(s)
	escaped := false
	resolved := make([]rune, 0, len(s))
	for index := 0; index < sLength; index++ {
		char := s[index]
		if escaped {
			switch char {
			case 'a', 'b', 'e', 'f', 'n', 'r', 't', '?':
				// Replace char based
				resolved = append(resolved, directCharEscapeValue[char]...)
			case '\\', '\'', '"', '`':
				// Replace escaped literals
				resolved = append(resolved, char)
			case 'x':
				// Replace hex with numbers
				index++
				a := s[index]
				index++
				b := s[index]
				resolved = append(resolved, a*16+b)
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
				resolved = append(resolved, a*4096+b*256+c*16+d)
			}
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else {
			resolved = append(resolved, char)
		}
	}
	return &ast2.Bytes{
		Contents: []byte(string(resolved)),
	}
}

func (simp *simplify) simplifyLiteral(literal *ast.BasicLiteralExpression) ast2.Expression {
	switch literal.DirectValue {
	case lexer.Integer, lexer.BinaryInteger, lexer.OctalInteger, lexer.HexadecimalInteger:
		return simp.simplifyInteger(literal.Token.String())
	case lexer.Float, lexer.ScientificFloat:
		return simp.simplifyFloat(literal.Token.String())
	case lexer.SingleQuoteString, lexer.DoubleQuoteString, lexer.CommandOutput:
		return simp.simplifyString(literal.Token.String())
	case lexer.True:
		return &ast2.True{}
	case lexer.False:
		return &ast2.False{}
	case lexer.None:
		return &ast2.None{}
	case lexer.ByteString:
		return simp.simplifyBytes(literal.Token.String())
	default:
		panic(fmt.Sprintf("unknown literal %d", literal.DirectValue))
	}
}
