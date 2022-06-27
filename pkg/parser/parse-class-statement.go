package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseClassStatement() (*ast2.ClassStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer2.IdentifierKind) {
		return nil, parser.newSyntaxError(ClassStatement)
	}
	name := &ast2.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast2.IExpression
	var base ast2.Node
	var parsingError error
	if parser.matchDirectValue(lexer2.OpenParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		for parser.hasNext() {
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			base, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := base.(ast2.IExpression); !ok {
				return nil, parser.expectingExpressionError(ClassStatement)
			}
			bases = append(bases, base.(ast2.IExpression))
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
				return nil, parser.newSyntaxError(ClassStatement)
			}
		}
		if !parser.matchDirectValue(lexer2.CloseParentheses) {
			return nil, parser.newSyntaxError(ClassStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(ClassStatement)
	}
	var body []ast2.Node
	var bodyNode ast2.Node
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
		return nil, parser.statementNeverEndedError(ClassStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.ClassStatement{
		Name:  name,
		Bases: bases,
		Body:  body,
	}, nil
}
