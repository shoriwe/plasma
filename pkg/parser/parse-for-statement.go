package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseForStatement() (*ast2.ForLoopStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var receivers []*ast2.Identifier
	for parser.hasNext() {
		if parser.matchDirectValue(lexer2.In) {
			break
		} else if !parser.matchKind(lexer2.IdentifierKind) {
			return nil, parser.newSyntaxError(ForStatement)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		receivers = append(receivers, &ast2.Identifier{
			Token: parser.currentToken,
		})
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
		} else if parser.matchDirectValue(lexer2.In) {
			break
		} else {
			return nil, parser.newSyntaxError(ForStatement)
		}
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer2.In) {
		return nil, parser.newSyntaxError(ForStatement)
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
	if _, ok := source.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(ForStatement)
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
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
		return nil, parser.statementNeverEndedError(ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.ForLoopStatement{
		Receivers: receivers,
		Source:    source.(ast2.IExpression),
		Body:      body,
	}, nil
}
