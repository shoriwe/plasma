package parser

import (
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

func (parser *Parser) parseGeneratorExpression(operation ast.IExpression) (*ast.GeneratorExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var variables []*ast.Identifier
	numberOfVariables := 0
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.In) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchKind(lexer.IdentifierKind) {
			return nil, parser.newSyntaxError(GeneratorExpression)
		}
		variables = append(variables, &ast.Identifier{
			Token: parser.currentToken,
		})
		numberOfVariables++
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	if numberOfVariables == 0 {
		return nil, parser.newSyntaxError(GeneratorExpression)
	}
	if !parser.matchDirectValue(lexer.In) {
		return nil, parser.newSyntaxError(GeneratorExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(GeneratorExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	// Finally detect the closing parentheses
	if !parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.newSyntaxError(GeneratorExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.GeneratorExpression{
		Operation: operation,
		Receivers: variables,
		Source:    source.(ast.IExpression),
	}, nil
}
