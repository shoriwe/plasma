package parser

import "github.com/shoriwe/gplasma/pkg/compiler/lexer"

func (parser *Parser) matchDirectValue(directValue lexer.DirectValue) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) matchKind(kind lexer.Kind) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.Kind == kind
}

func (parser *Parser) currentLine() int {
	if parser.currentToken == nil {
		return 0
	}
	return parser.currentToken.Line
}

func (parser *Parser) matchString(value string) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.String() == value
}
