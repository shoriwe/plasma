package parser

import (
	"errors"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
)

type Parser struct {
	lexer    *lexer.Lexer
	complete bool
}

func (parser *Parser) hasNext() bool {
	return !parser.complete
}

func (parser *Parser) next() (ast.Node, error) {
	if !parser.lexer.HasNext() {
		parser.complete = true
		return nil, nil
	}
	token, tokenizingError := parser.lexer.Next()
	if tokenizingError != nil {
		parser.complete = true
		return nil, tokenizingError
	}
	var node ast.Node
	var parsingError error
	switch token.Kind {
	case lexer.CommandOutput,
		lexer.SingleQuoteString,
		lexer.DoubleQuoteString,
		lexer.ByteString,
		lexer.Integer,
		lexer.HexadecimalInteger,
		lexer.BinaryInteger,
		lexer.OctalInteger,
		lexer.Float,
		lexer.ScientificFloat:
	}
	return node, parsingError
}

func (parser *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	for ; parser.hasNext(); {
		node, parsingError := parser.next()
		if parsingError != nil {
			return nil, parsingError
		}
		switch node.(type) {
		case *ast.BeginStatement:
			if program.Begin != nil {
				return nil, errors.New("found multiples BEGIN statements in the code")
			}
			program.Begin = node.(*ast.BeginStatement)
		case *ast.EndStatement:
			if program.End != nil {
				return nil, errors.New("found multiples END statements in the code")
			}
			program.End = node.(*ast.EndStatement)
		default:
			program.Body = append(program.Body, node)
		}
	}
	return program, nil
}

func NewParser(lexer_ *lexer.Lexer) *Parser {
	return &Parser{lexer_, false}
}
