package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseUnaryExpression() (ast2.Node, error) {
	// Do something to parse Unary
	if parser.matchKind(lexer2.Operator) {
		switch parser.currentToken.DirectValue {
		case lexer2.Sub, lexer2.Add, lexer2.NegateBits, lexer2.SignNot, lexer2.Not:
			operator := parser.currentToken
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}

			x, parsingError := parser.parseUnaryExpression()
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast2.IExpression); !ok {
				return nil, parser.expectingExpressionError(PointerExpression)
			}
			return &ast2.UnaryExpression{
				Operator: operator,
				X:        x.(ast2.IExpression),
			}, nil
		}
	}
	return parser.parsePrimaryExpression()
}
