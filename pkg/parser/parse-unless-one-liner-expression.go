package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseUnlessOneLinerExpression(result ast2.IExpression) (*ast2.UnlessOneLinerExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(UnlessOneLinerExpression)
	}
	if !parser.matchDirectValue(lexer.Else) {
		return &ast2.UnlessOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast2.IExpression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var elseResult ast2.Node

	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(OneLineElseBlock)
	}
	return &ast2.UnlessOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast2.IExpression),
		ElseResult: elseResult.(ast2.IExpression),
	}, nil
}
