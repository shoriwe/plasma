package parser

import (
	"github.com/shoriwe/gplasma/pkg/ast"
)

func (parser *Parser) parseRaiseStatement() (*ast.RaiseStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	x, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := x.(ast.Expression); !ok {
		return nil, parser.expectingExpressionError(RaiseStatement)
	}
	return &ast.RaiseStatement{
		X: x.(ast.Expression),
	}, nil
}
