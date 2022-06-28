package parser

import ast2 "github.com/shoriwe/gplasma/pkg/ast"

func (parser *Parser) parseRequireStatement() (*ast2.RequireStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	x, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := x.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(RequireStatement)
	}
	return &ast2.RequireStatement{
		X: x.(ast2.IExpression),
	}, nil
}
