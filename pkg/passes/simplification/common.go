package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func literalIsInteger(literal *ast.BasicLiteralExpression) bool {
	switch literal.DirectValue {
	case lexer.Integer,
		lexer.BinaryInteger,
		lexer.OctalInteger,
		lexer.HexadecimalInteger:
		return true
	default:
		return false
	}
}

func literalIsFloat(literal *ast.BasicLiteralExpression) bool {
	switch literal.DirectValue {
	case lexer.Float,
		lexer.ScientificFloat:
		return true
	default:
		return false
	}
}

func literalIsString(literal *ast.BasicLiteralExpression) bool {
	switch literal.DirectValue {
	case lexer.SingleQuoteString,
		lexer.DoubleQuoteString,
		lexer.CommandOutput:
		return true
	default:
		return false
	}
}

func literalIsBytesString(literal *ast.BasicLiteralExpression) bool {
	switch literal.DirectValue {
	case lexer.ByteString:
		return true
	default:
		return false
	}
}

func literalIsBytes(literal *ast.BasicLiteralExpression) bool {
	switch literal.DirectValue {
	case lexer.ByteString:
		return true
	default:
		return false
	}
}

func booleanValue(value bool) *ast.Identifier {
	var (
		s      string
		sValue lexer.DirectValue
	)
	if value {
		s = lexer.TrueString
		sValue = lexer.True
	} else {
		s = lexer.FalseString
		sValue = lexer.False
	}
	return &ast.Identifier{
		Token: &lexer.Token{
			Contents:    []rune(s),
			DirectValue: sValue,
			Kind:        lexer.Boolean,
		},
	}
}
