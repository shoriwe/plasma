package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseDoWhileStatement() (*ast2.DoWhileStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast2.Node
	var bodyNode ast2.Node
	var parsingError error
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(DoWhileStatement)
	}
	// Parse Body
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer2.While) {
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
	// Parse Condition
	if !parser.matchDirectValue(lexer2.While) {
		return nil, parser.newSyntaxError(DoWhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	var condition ast2.Node
	condition, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(WhileStatement)
	}
	return &ast2.DoWhileStatement{
		Condition: condition.(ast2.IExpression),
		Body:      body,
	}, nil
}
