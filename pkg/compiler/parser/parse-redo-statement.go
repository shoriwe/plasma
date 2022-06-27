package parser

import "github.com/shoriwe/gplasma/pkg/compiler/ast"

func (parser *Parser) parseRedoStatement() (*ast.RedoStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.RedoStatement{}, nil
}
