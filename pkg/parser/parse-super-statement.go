package parser

import ast2 "github.com/shoriwe/gplasma/pkg/ast"

func (parser *Parser) parseSuperExpression() (*ast2.SuperExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	x, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := x.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(SuperExpression)
	}
	return &ast2.SuperExpression{
		X: x.(ast2.IExpression),
	}, nil
}
