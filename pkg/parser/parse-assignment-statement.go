package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
)

func (parser *Parser) parseAssignmentStatement(leftHandSide ast2.IExpression) (*ast2.AssignStatement, error) {
	assignmentToken := parser.currentToken
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	rightHandSide, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := rightHandSide.(ast2.IExpression); !ok {
		return nil, parser.expectingExpressionError(AssignStatement)
	}
	return &ast2.AssignStatement{
		LeftHandSide:   leftHandSide,
		AssignOperator: assignmentToken,
		RightHandSide:  rightHandSide.(ast2.IExpression),
	}, nil
}
