package parser

import (
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

func (parser *Parser) parsePrimaryExpression() (ast.Node, error) {
	var parsedNode ast.Node
	var parsingError error
	parsedNode, parsingError = parser.parseOperand()
	if parsingError != nil {
		return nil, parsingError
	}
expressionPendingLoop:
	for {
		switch parser.currentToken.DirectValue {
		case lexer.Dot: // Is selector
			parsedNode, parsingError = parser.parseSelectorExpression(parsedNode.(ast.IExpression))
		case lexer.OpenParentheses: // Is function Call
			parsedNode, parsingError = parser.parseMethodInvocationExpression(parsedNode.(ast.IExpression))
		case lexer.OpenSquareBracket: // Is indexing
			parsedNode, parsingError = parser.parseIndexExpression(parsedNode.(ast.IExpression))
		case lexer.If: // One line If
			parsedNode, parsingError = parser.parseIfOneLinerExpression(parsedNode.(ast.IExpression))
		case lexer.Unless: // One line Unless
			parsedNode, parsingError = parser.parseUnlessOneLinerExpression(parsedNode.(ast.IExpression))
		default:
			if parser.matchKind(lexer.Assignment) {
				parsedNode, parsingError = parser.parseAssignmentStatement(parsedNode.(ast.IExpression))
			}
			break expressionPendingLoop
		}
	}
	if parsingError != nil {
		return nil, parsingError
	}
	return parsedNode, nil
}
