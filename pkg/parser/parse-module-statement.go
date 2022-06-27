package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseModuleStatement() (*ast2.ModuleStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer2.IdentifierKind) {
		return nil, parser.newSyntaxError(ModuleStatement)
	}
	name := &ast2.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(ModuleStatement)
	}
	var body []ast2.Node
	var bodyNode ast2.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer2.End) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer2.End) {
		return nil, parser.statementNeverEndedError(ModuleStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.ModuleStatement{
		Name: name,
		Body: body,
	}, nil
}
