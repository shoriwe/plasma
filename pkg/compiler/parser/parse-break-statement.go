package parser

import "github.com/shoriwe/gplasma/pkg/compiler/ast"

func (parser *Parser) parseBreakStatement() (*ast.BreakStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BreakStatement{}, nil
}
