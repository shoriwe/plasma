package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseInterfaceStatement() (*ast2.InterfaceStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer2.IdentifierKind) {
		return nil, parser.newSyntaxError(InterfaceStatement)
	}
	name := &ast2.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast2.IExpression
	var base ast2.Node
	var parsingError error
	if parser.matchDirectValue(lexer2.OpenParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		for parser.hasNext() {
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			base, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := base.(ast2.IExpression); !ok {
				return nil, parser.newSyntaxError(InterfaceStatement)
			}
			bases = append(bases, base.(ast2.IExpression))
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if parser.matchDirectValue(lexer2.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.matchDirectValue(lexer2.CloseParentheses) {
				break
			} else {
				return nil, parser.newSyntaxError(InterfaceStatement)
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchDirectValue(lexer2.CloseParentheses) {
			return nil, parser.newSyntaxError(InterfaceStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(InterfaceStatement)
	}
	var methods []*ast2.FunctionDefinitionStatement
	var node ast2.Node
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
		node, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := node.(*ast2.FunctionDefinitionStatement); !ok {
			return nil, parser.expectingFunctionDefinition(InterfaceStatement)
		}
		methods = append(methods, node.(*ast2.FunctionDefinitionStatement))
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer2.End) {
		return nil, parser.statementNeverEndedError(InterfaceStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast2.InterfaceStatement{
		Name:              name,
		Bases:             bases,
		MethodDefinitions: methods,
	}, nil
}
