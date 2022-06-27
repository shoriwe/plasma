package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parsePrimaryExpression() (ast2.Node, error) {
	var parsedNode ast2.Node
	var parsingError error
	parsedNode, parsingError = parser.parseOperand()
	if parsingError != nil {
		return nil, parsingError
	}
expressionPendingLoop:
	for {
		switch parser.currentToken.DirectValue {
		case lexer2.Dot: // Is selector
			parsedNode, parsingError = parser.parseSelectorExpression(parsedNode.(ast2.IExpression))
		case lexer2.OpenParentheses: // Is function Call
			parsedNode, parsingError = parser.parseMethodInvocationExpression(parsedNode.(ast2.IExpression))
		case lexer2.OpenSquareBracket: // Is indexing
			parsedNode, parsingError = parser.parseIndexExpression(parsedNode.(ast2.IExpression))
		case lexer2.If: // One line If
			parsedNode, parsingError = parser.parseIfOneLinerExpression(parsedNode.(ast2.IExpression))
		case lexer2.Unless: // One line Unless
			parsedNode, parsingError = parser.parseUnlessOneLinerExpression(parsedNode.(ast2.IExpression))
		default:
			if parser.matchKind(lexer2.Assignment) {
				parsedNode, parsingError = parser.parseAssignmentStatement(parsedNode.(ast2.IExpression))
			}
			break expressionPendingLoop
		}
	}
	if parsingError != nil {
		return nil, parsingError
	}
	return parsedNode, nil
}
