package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseReturnStatement() (*ast2.ReturnStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var results []ast2.IExpression
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) || parser.matchKind(lexer2.EOF) {
			break
		}

		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast2.IExpression); !ok {
			return nil, parser.expectingExpressionError(ReturnStatement)
		}
		results = append(results, result.(ast2.IExpression))
		if parser.matchDirectValue(lexer2.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.matchKind(lexer2.Separator) || parser.matchKind(lexer2.EOF)) {
			return nil, parser.newSyntaxError(ReturnStatement)
		}
	}
	return &ast2.ReturnStatement{
		Results: results,
	}, nil
}
