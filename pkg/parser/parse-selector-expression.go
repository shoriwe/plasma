package parser

import (
	"github.com/shoriwe/plasma/pkg/ast"
	"github.com/shoriwe/plasma/pkg/lexer"
)

func (parser *Parser) parseSelectorExpression(expression ast.Expression) (*ast.SelectorExpression, error) {
	selector := expression
	for parser.hasNext() {
		if !parser.matchDirectValue(lexer.Dot) {
			break
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		identifier := parser.currentToken
		if identifier.Kind != lexer.IdentifierKind {
			return nil, parser.newSyntaxError(SelectorExpression)
		}
		selector = &ast.SelectorExpression{
			X: selector,
			Identifier: &ast.Identifier{
				Token: identifier,
			},
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	return selector.(*ast.SelectorExpression), nil
}
