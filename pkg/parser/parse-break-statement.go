package parser

import (
	"github.com/shoriwe/plasma/pkg/ast"
)

func (parser *Parser) parseBreakStatement() (*ast.BreakStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BreakStatement{}, nil
}
