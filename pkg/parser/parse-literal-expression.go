package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseLiteral() (ast.IExpression, error) {
	if !parser.matchKind(lexer2.Literal) &&
		!parser.matchKind(lexer2.Boolean) &&
		!parser.matchKind(lexer2.NoneType) {
		return nil, parser.invalidTokenKind()
	}

	switch parser.currentToken.DirectValue {
	case lexer2.SingleQuoteString, lexer2.DoubleQuoteString, lexer2.ByteString,
		lexer2.Integer, lexer2.HexadecimalInteger, lexer2.BinaryInteger, lexer2.OctalInteger,
		lexer2.Float, lexer2.ScientificFloat,
		lexer2.True, lexer2.False, lexer2.None:
		currentToken := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.BasicLiteralExpression{
			Token:       currentToken,
			Kind:        currentToken.Kind,
			DirectValue: currentToken.DirectValue,
		}, nil
	}
	return nil, parser.invalidTokenKind()
}
