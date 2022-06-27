package parser

import (
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

func (parser *Parser) parseParentheses() (ast.IExpression, error) {
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
	if _, ok := firstExpression.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(ParenthesesExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.matchDirectValue(lexer.For) {
		return parser.parseGeneratorExpression(firstExpression.(ast.IExpression))
	}
	if parser.matchDirectValue(lexer.CloseParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.ParenthesesExpression{
			X: firstExpression.(ast.IExpression),
		}, nil
	}
	if !parser.matchDirectValue(lexer.Comma) {
		return nil, parser.newSyntaxError(ParenthesesExpression)
	}
	var values []ast.IExpression
	values = append(values, firstExpression.(ast.IExpression))
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var nextValue ast.Node
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
		if _, ok := nextValue.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(ParenthesesExpression)
		}
		values = append(values, nextValue.(ast.IExpression))
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
	return &ast.TupleExpression{
		Values: values,
	}, nil
}
