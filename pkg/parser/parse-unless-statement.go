package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseUnlessStatement() (*ast2.UnlessStatement, error) {
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
		return nil, parser.expectingExpressionError(UnlessStatement)
	}
	if !parser.matchDirectValue(lexer2.NewLine) {
		return nil, parser.newSyntaxError(UnlessStatement)
	}
	// Parse Unless
	root := &ast2.UnlessStatement{
		Condition: condition.(ast2.IExpression),
		Body:      []ast2.Node{},
		Else:      []ast2.Node{},
	}
	var bodyNode ast2.Node
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer2.Elif) ||
				parser.matchDirectValue(lexer2.Else) ||
				parser.matchDirectValue(lexer2.End) {
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
	if parser.matchDirectValue(lexer2.Elif) {
		var elifBody []ast2.Node
	parsingElifLoop:
		for parser.hasNext() {
			if parser.matchKind(lexer2.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer2.Else) ||
					parser.matchDirectValue(lexer2.End) {
					break
				}
				continue
			}
			if !parser.matchDirectValue(lexer2.Elif) {
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
			var elifCondition ast2.Node
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := elifCondition.(ast2.IExpression); !ok {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			if !parser.matchDirectValue(lexer2.NewLine) {
				return nil, parser.newSyntaxError(ElifBlock)
			}
			var elifBodyNode ast2.Node
			for parser.hasNext() {
				if parser.matchKind(lexer2.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer2.Else) ||
						parser.matchDirectValue(lexer2.End) ||
						parser.matchDirectValue(lexer2.Elif) {
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
				&ast2.UnlessStatement{
					Condition: elifCondition.(ast2.IExpression),
					Body:      elifBody,
				},
			)
			lastCondition = lastCondition.Else[0].(*ast2.UnlessStatement)
			if parser.matchDirectValue(lexer2.Else) ||
				parser.matchDirectValue(lexer2.End) {
				break parsingElifLoop
			}
		}
	}
	// Parse Default
	if parser.matchDirectValue(lexer2.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast2.Node
		if !parser.matchDirectValue(lexer2.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
		}
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
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			lastCondition.Else = append(lastCondition.Else, elseBodyNode)
		}
	}
	if !parser.matchDirectValue(lexer2.End) {
		return nil, parser.statementNeverEndedError(UnlessStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}
