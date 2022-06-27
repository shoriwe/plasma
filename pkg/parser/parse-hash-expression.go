package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseHashExpression() (*ast2.HashExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var values []*ast2.KeyValue
	var leftHandSide ast2.Node
	var rightHandSide ast2.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseBrace) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		leftHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := leftHandSide.(ast2.IExpression); !ok {
			return nil, parser.expectingExpressionError(HashExpression)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchDirectValue(lexer.Colon) {
			return nil, parser.newSyntaxError(HashExpression)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		rightHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast2.IExpression); !ok {
			return nil, parser.expectingExpressionError(HashExpression)
		}
		values = append(values, &ast2.KeyValue{
			Key:   leftHandSide.(ast2.IExpression),
			Value: rightHandSide.(ast2.IExpression),
		})
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
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.HashExpression{
		Values: values,
	}, nil
}
