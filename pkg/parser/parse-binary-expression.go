package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseBinaryExpression(precedence lexer2.DirectValue) (ast2.Node, error) {
	var leftHandSide ast2.Node
	var rightHandSide ast2.Node
	var parsingError error
	leftHandSide, parsingError = parser.parseUnaryExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := leftHandSide.(ast2.IStatement); ok {
		return leftHandSide, nil
	}
	for parser.hasNext() {
		if !parser.matchKind(lexer2.Operator) &&
			!parser.matchKind(lexer2.Comparator) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		operator := parser.currentToken
		operatorPrecedence := parser.currentToken.DirectValue
		if operatorPrecedence < precedence {
			return leftHandSide, nil
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}

		rightHandSide, parsingError = parser.parseBinaryExpression(operatorPrecedence + 1)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast2.IExpression); !ok {
			return nil, parser.expectingExpressionError(BinaryExpression)
		}

		leftHandSide = &ast2.BinaryExpression{
			LeftHandSide:  leftHandSide.(ast2.IExpression),
			Operator:      operator,
			RightHandSide: rightHandSide.(ast2.IExpression),
		}
	}
	return leftHandSide, nil
}
