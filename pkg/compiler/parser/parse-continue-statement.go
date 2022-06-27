package parser

import "github.com/shoriwe/gplasma/pkg/compiler/ast"

func (parser *Parser) parseContinueStatement() (*ast.ContinueStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ContinueStatement{}, nil
}
