package parser

import (
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

func (parser *Parser) parseOperand() (ast2.Node, error) {
	switch parser.currentToken.Kind {
	case lexer2.Literal, lexer2.Boolean, lexer2.NoneType:
		return parser.parseLiteral()
	case lexer2.IdentifierKind:
		identifier := parser.currentToken

		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast2.Identifier{
			Token: identifier,
		}, nil
	case lexer2.Keyword:
		switch parser.currentToken.DirectValue {
		case lexer2.Lambda:
			return parser.parseLambdaExpression()
		case lexer2.Super:
			return parser.parseSuperExpression()
		case lexer2.Delete:
			return parser.parseDeleteStatement()
		case lexer2.Require:
			return parser.parseRequireStatement()
		case lexer2.While:
			return parser.parseWhileStatement()
		case lexer2.For:
			return parser.parseForStatement()
		case lexer2.Until:
			return parser.parseUntilStatement()
		case lexer2.If:
			return parser.parseIfStatement()
		case lexer2.Unless:
			return parser.parseUnlessStatement()
		case lexer2.Switch:
			return parser.parseSwitchStatement()
		case lexer2.Module:
			return parser.parseModuleStatement()
		case lexer2.Def:
			return parser.parseFunctionDefinitionStatement()
		case lexer2.Generator:
			return parser.parseGeneratorDefinitionStatement()
		case lexer2.Interface:
			return parser.parseInterfaceStatement()
		case lexer2.Class:
			return parser.parseClassStatement()
		case lexer2.Raise:
			return parser.parseRaiseStatement()
		case lexer2.Try:
			return parser.parseTryStatement()
		case lexer2.Return:
			return parser.parseReturnStatement()
		case lexer2.Yield:
			return parser.parseYieldStatement()
		case lexer2.Continue:
			return parser.parseContinueStatement()
		case lexer2.Break:
			return parser.parseBreakStatement()
		case lexer2.Redo:
			return parser.parseRedoStatement()
		case lexer2.Pass:
			return parser.parsePassStatement()
		case lexer2.Do:
			return parser.parseDoWhileStatement()
		}
	case lexer2.Punctuation:
		switch parser.currentToken.DirectValue {
		case lexer2.OpenParentheses:
			return parser.parseParentheses()
		case lexer2.OpenSquareBracket: // Parse Arrays
			return parser.parseArrayExpression()
		case lexer2.OpenBrace: // Parse Dictionaries
			return parser.parseHashExpression()
		}
	}
	return nil, UnknownToken
}
