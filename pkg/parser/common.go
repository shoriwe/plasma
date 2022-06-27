package parser

import (
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) removeNewLines() error {
	for parser.matchDirectValue(lexer2.NewLine) {
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return tokenizingError
		}
	}
	return nil
}

func (parser *Parser) matchDirectValue(directValue lexer2.DirectValue) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) matchKind(kind lexer2.Kind) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.Kind == kind
}

func (parser *Parser) matchString(value string) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.String() == value
}
