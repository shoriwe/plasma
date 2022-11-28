package parser

import (
	"github.com/shoriwe/plasma/pkg/ast"
)

func (parser *Parser) parseContinueStatement() (*ast.ContinueStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ContinueStatement{}, nil
}
