package parser

import (
	"github.com/shoriwe/plasma/pkg/ast"
)

func (parser *Parser) parsePassStatement() (*ast.PassStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.PassStatement{}, nil
}
