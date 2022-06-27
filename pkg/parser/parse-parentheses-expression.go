package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseParentheses() (ast2.IExpression, error) {
	/*
		This should also parse generators
	*/
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.newSyntaxError(ParenthesesExpression)
	}

	firstExpression, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := firstExpression.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(ParenthesesExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.matchDirectValue(lexer.For) {
		return parser.parseGeneratorExpression(firstExpression.(ast2.IExpression))
	}
	if parser.matchDirectValue(lexer.CloseParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast2.ParenthesesExpression{
			X: firstExpression.(ast2.IExpression),
		}, nil
	}
	if !parser.matchDirectValue(lexer.Comma) {
		return nil, parser.newSyntaxError(ParenthesesExpression)
	}
	var values []ast2.IExpression
	values = append(values, firstExpression.(ast2.IExpression))
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var nextValue ast2.Node
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		nextValue, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := nextValue.(ast2.IExpression); !ok {
			return nil, parser.expectingExpressionError(ParenthesesExpression)
		}
		values = append(values, nextValue.(ast2.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirectValue(lexer.CloseParentheses) {
			return nil, parser.newSyntaxError(TupleExpression)
		}
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.expressionNeverClosedError(TupleExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.TupleExpression{
		Values: values,
	}, nil
}
