package parser

import (
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

func (parser *Parser) parseUnaryExpression() (ast.Node, error) {
	// Do something to parse Unary
	if parser.matchKind(lexer.Operator) {
		switch parser.currentToken.DirectValue {
		case lexer.Sub, lexer.Add, lexer.NegateBits, lexer.SignNot, lexer.Not:
			operator := parser.currentToken
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}

			x, parsingError := parser.parseUnaryExpression()
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast.IExpression); !ok {
				return nil, parser.expectingExpressionError(PointerExpression)
			}
			return &ast.UnaryExpression{
				Operator: operator,
				X:        x.(ast.IExpression),
			}, nil
		}
	}
	return parser.parsePrimaryExpression()
}
