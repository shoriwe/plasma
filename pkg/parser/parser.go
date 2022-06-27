package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/common"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

type Parser struct {
	lineStack    common.Stack[int]
	lexer        *lexer2.Lexer
	complete     bool
	currentToken *lexer2.Token
}

func (parser *Parser) Parse() (*ast2.Program, error) {
	result := &ast2.Program{
		Begin: nil,
		End:   nil,
		Body:  nil,
	}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var beginStatement *ast2.BeginStatement
	var endStatement *ast2.EndStatement
	var parsedExpression ast2.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer2.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		switch {
		case parser.matchDirectValue(lexer2.BEGIN):
			if result.Begin != nil {
				return nil, BeginRepeated
			}
			beginStatement, parsingError = parser.parseBeginStatement()
			if parsingError != nil {
				return nil, parsingError
			}
			result.Begin = beginStatement
		case parser.matchDirectValue(lexer2.END):
			if result.End != nil {
				return nil, EndRepeated
			}
			endStatement, parsingError = parser.parseEndStatement()
			if parsingError != nil {
				return nil, parsingError
			}
			result.End = endStatement
		default:
			parsedExpression, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			result.Body = append(result.Body, parsedExpression)
		}
	}
	parser.complete = true
	return result, nil
}

func NewParser(lexer_ *lexer2.Lexer) *Parser {
	return &Parser{
		lineStack:    common.Stack[int]{},
		lexer:        lexer_,
		complete:     false,
		currentToken: nil,
	}
}
