package parser

import (
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

func (parser *Parser) parseUnlessStatement() (*ast.UnlessStatement, error) {
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
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(UnlessStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(UnlessStatement)
	}
	// Parse Unless
	root := &ast.UnlessStatement{
		Condition: condition.(ast.IExpression),
		Body:      []ast.Node{},
		Else:      []ast.Node{},
	}
	var bodyNode ast.Node
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.Elif) ||
				parser.matchDirectValue(lexer.Else) ||
				parser.matchDirectValue(lexer.End) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		root.Body = append(root.Body, bodyNode)
	}
	// Parse Elifs
	lastCondition := root
	if parser.matchDirectValue(lexer.Elif) {
		var elifBody []ast.Node
	parsingElifLoop:
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.Else) ||
					parser.matchDirectValue(lexer.End) {
					break
				}
				continue
			}
			if !parser.matchDirectValue(lexer.Elif) {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			var elifCondition ast.Node
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := elifCondition.(ast.IExpression); !ok {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			if !parser.matchDirectValue(lexer.NewLine) {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			var elifBodyNode ast.Node
			for parser.hasNext() {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer.Else) ||
						parser.matchDirectValue(lexer.End) ||
						parser.matchDirectValue(lexer.Elif) {
						break
					}
					continue
				}
				elifBodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				elifBody = append(elifBody, elifBodyNode)
			}
			lastCondition.Else = append(
				lastCondition.Else,
				&ast.UnlessStatement{
					Condition: elifCondition.(ast.IExpression),
					Body:      elifBody,
				},
			)
			lastCondition = lastCondition.Else[0].(*ast.UnlessStatement)
			if parser.matchDirectValue(lexer.Else) ||
				parser.matchDirectValue(lexer.End) {
				break parsingElifLoop
			}
		}
	}
	// Parse Default
	if parser.matchDirectValue(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast.Node
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
		}
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.End) {
					break
				}
				continue
			}
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			lastCondition.Else = append(lastCondition.Else, elseBodyNode)
		}
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(UnlessStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}
