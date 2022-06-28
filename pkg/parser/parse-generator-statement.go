package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseGeneratorDefinitionStatement() (*ast2.GeneratorDefinitionStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer2.IdentifierKind) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	name := &ast2.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer2.OpenParentheses) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []*ast2.Identifier
	for parser.hasNext() {
		if parser.matchDirectValue(lexer2.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchKind(lexer2.IdentifierKind) {
			return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
		}
		argument := &ast2.Identifier{
			Token: parser.currentToken,
		}
		arguments = append(arguments, argument)
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer2.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.matchDirectValue(lexer2.CloseParentheses) {
			break
		} else {
			return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
		}
	}
	if !parser.matchDirectValue(lexer2.CloseParentheses) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(GeneratorDefinitionStatement)
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
		return nil, parser.statementNeverEndedError(GeneratorDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.GeneratorDefinitionStatement{
		Name:      name,
		Arguments: arguments,
		Body:      body,
	}, nil
}
