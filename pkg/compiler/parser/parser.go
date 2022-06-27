package parser

import (
	"github.com/shoriwe/gplasma/pkg/common"
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
)

type Parser struct {
	lineStack    common.Stack[int]
	lexer        *lexer.Lexer
	complete     bool
	currentToken *lexer.Token
}

func (parser *Parser) parseForStatement() (*ast.ForLoopStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var receivers []*ast.Identifier
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.In) {
			break
		} else if !parser.matchKind(lexer.IdentifierKind) {
			return nil, parser.newSyntaxError(ForStatement)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		receivers = append(receivers, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.matchDirectValue(lexer.In) {
			break
		} else {
			return nil, parser.newSyntaxError(ForStatement)
		}
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.In) {
		return nil, parser.newSyntaxError(ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(ForStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast.Node
	var bodyNode ast.Node
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(ForStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ForLoopStatement{
		Receivers: receivers,
		Source:    source.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseUntilStatement() (*ast.UntilLoopStatement, error) {
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
		return nil, parser.newSyntaxError(UntilStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(UntilStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bodyNode ast.Node
	var body []ast.Node
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(UntilStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.UntilLoopStatement{
		Condition: condition.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseModuleStatement() (*ast.ModuleStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(ModuleStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(ModuleStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(ModuleStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ModuleStatement{
		Name: name,
		Body: body,
	}, nil
}

func (parser *Parser) parseFunctionDefinitionStatement() (*ast.FunctionDefinitionStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(FunctionDefinitionStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.OpenParentheses) {
		return nil, parser.newSyntaxError(FunctionDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []*ast.Identifier
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchKind(lexer.IdentifierKind) {
			return nil, parser.newSyntaxError(FunctionDefinitionStatement)
		}
		argument := &ast.Identifier{
			Token: parser.currentToken,
		}
		arguments = append(arguments, argument)
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.matchDirectValue(lexer.CloseParentheses) {
			break
		} else {
			return nil, parser.newSyntaxError(FunctionDefinitionStatement)
		}
	}
	if !parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.newSyntaxError(FunctionDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(FunctionDefinitionStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(FunctionDefinitionStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.FunctionDefinitionStatement{
		Name:      name,
		Arguments: arguments,
		Body:      body,
	}, nil
}

func (parser *Parser) parseClassStatement() (*ast.ClassStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(ClassStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.IExpression
	var base ast.Node
	var parsingError error
	if parser.matchDirectValue(lexer.OpenParentheses) {
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
			if _, ok := base.(ast.IExpression); !ok {
				return nil, parser.expectingExpressionError(ClassStatement)
			}
			bases = append(bases, base.(ast.IExpression))
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if parser.matchDirectValue(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.matchDirectValue(lexer.CloseParentheses) {
				break
			} else {
				return nil, parser.newSyntaxError(ClassStatement)
			}
		}
		if !parser.matchDirectValue(lexer.CloseParentheses) {
			return nil, parser.newSyntaxError(ClassStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(ClassStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(ClassStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ClassStatement{
		Name:  name,
		Bases: bases,
		Body:  body,
	}, nil
}

func (parser *Parser) parseRaiseStatement() (*ast.RaiseStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	x, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := x.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(RaiseStatement)
	}
	return &ast.RaiseStatement{
		X: x.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseTryStatement() (*ast.TryStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(TryStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.End) ||
				parser.matchDirectValue(lexer.Except) ||
				parser.matchDirectValue(lexer.Else) ||
				parser.matchDirectValue(lexer.Finally) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	var exceptBlocks []*ast.ExceptBlock
	for parser.hasNext() {
		if !parser.matchDirectValue(lexer.Except) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var targets []ast.IExpression
		var target ast.Node
		for parser.hasNext() {
			if parser.matchDirectValue(lexer.NewLine) ||
				parser.matchDirectValue(lexer.As) {
				break
			}
			target, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := target.(ast.IExpression); !ok {
				return nil, parser.newSyntaxError(ExceptBlock)
			}
			targets = append(targets, target.(ast.IExpression))
			if parser.matchDirectValue(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			}
		}
		var captureName *ast.Identifier
		if parser.matchDirectValue(lexer.As) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError := parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if !parser.matchKind(lexer.IdentifierKind) {
				return nil, parser.newSyntaxError(ExceptBlock)
			}
			captureName = &ast.Identifier{
				Token: parser.currentToken,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(TryStatement)
		}
		var exceptBody []ast.Node
		var exceptBodyNode ast.Node
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.End) ||
					parser.matchDirectValue(lexer.Except) ||
					parser.matchDirectValue(lexer.Else) ||
					parser.matchDirectValue(lexer.Finally) {
					break
				}
				continue
			}
			exceptBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			exceptBody = append(exceptBody, exceptBodyNode)
		}
		exceptBlocks = append(exceptBlocks, &ast.ExceptBlock{
			Targets:     targets,
			CaptureName: captureName,
			Body:        exceptBody,
		})
	}
	var elseBody []ast.Node
	var elseBodyNode ast.Node
	if parser.matchDirectValue(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(ElseBlock)
		}
		for parser.hasNext() {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirectValue(lexer.End) ||
					parser.matchDirectValue(lexer.Finally) {
					break
				}
				continue
			}
			elseBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			elseBody = append(elseBody, elseBodyNode)
		}
	}
	var finallyBody []ast.Node
	var finallyBodyNode ast.Node
	if parser.matchDirectValue(lexer.Finally) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(FinallyBlock)
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
			finallyBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			finallyBody = append(finallyBody, finallyBodyNode)
		}
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.newSyntaxError(TryStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.TryStatement{
		Body:         body,
		ExceptBlocks: exceptBlocks,
		Else:         elseBody,
		Finally:      finallyBody,
	}, nil
}

func (parser *Parser) parseBeginStatement() (*ast.BeginStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(BeginStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(BeginStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BeginStatement{
		Body: body,
	}, nil
}

func (parser *Parser) parseEndStatement() (*ast.EndStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(EndStatement)
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
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
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(EndStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.EndStatement{
		Body: body,
	}, nil
}

func (parser *Parser) parseInterfaceStatement() (*ast.InterfaceStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, parser.newSyntaxError(InterfaceStatement)
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.IExpression
	var base ast.Node
	var parsingError error
	if parser.matchDirectValue(lexer.OpenParentheses) {
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
			if _, ok := base.(ast.IExpression); !ok {
				return nil, parser.newSyntaxError(InterfaceStatement)
			}
			bases = append(bases, base.(ast.IExpression))
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			if parser.matchDirectValue(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.matchDirectValue(lexer.CloseParentheses) {
				break
			} else {
				return nil, parser.newSyntaxError(InterfaceStatement)
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchDirectValue(lexer.CloseParentheses) {
			return nil, parser.newSyntaxError(InterfaceStatement)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(InterfaceStatement)
	}
	var methods []*ast.FunctionDefinitionStatement
	var node ast.Node
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
		node, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := node.(*ast.FunctionDefinitionStatement); !ok {
			return nil, parser.expectingFunctionDefinition(InterfaceStatement)
		}
		methods = append(methods, node.(*ast.FunctionDefinitionStatement))
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(InterfaceStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.InterfaceStatement{
		Name:              name,
		Bases:             bases,
		MethodDefinitions: methods,
	}, nil
}

func (parser *Parser) parseLiteral() (ast.IExpression, error) {
	if !parser.matchKind(lexer.Literal) &&
		!parser.matchKind(lexer.Boolean) &&
		!parser.matchKind(lexer.NoneType) {
		return nil, parser.invalidTokenKind()
	}

	switch parser.currentToken.DirectValue {
	case lexer.SingleQuoteString, lexer.DoubleQuoteString, lexer.ByteString,
		lexer.Integer, lexer.HexadecimalInteger, lexer.BinaryInteger, lexer.OctalInteger,
		lexer.Float, lexer.ScientificFloat,
		lexer.True, lexer.False, lexer.None:
		currentToken := parser.currentToken
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.BasicLiteralExpression{
			Token:       currentToken,
			Kind:        currentToken.Kind,
			DirectValue: currentToken.DirectValue,
		}, nil
	}
	return nil, parser.invalidTokenKind()
}

func (parser *Parser) parseBinaryExpression(precedence lexer.DirectValue) (ast.Node, error) {
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	leftHandSide, parsingError = parser.parseUnaryExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := leftHandSide.(ast.Statement); ok {
		return leftHandSide, nil
	}
	for parser.hasNext() {
		if !parser.matchKind(lexer.Operator) &&
			!parser.matchKind(lexer.Comparator) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		operator := parser.currentToken
		operatorPrecedence := parser.currentToken.DirectValue
		if operatorPrecedence < precedence {
			return leftHandSide, nil
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}

		rightHandSide, parsingError = parser.parseBinaryExpression(operatorPrecedence + 1)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(BinaryExpression)
		}

		leftHandSide = &ast.BinaryExpression{
			LeftHandSide:  leftHandSide.(ast.IExpression),
			Operator:      operator,
			RightHandSide: rightHandSide.(ast.IExpression),
		}
	}
	return leftHandSide, nil
}

func (parser *Parser) parseUnaryExpression() (ast.Node, error) {
	// Do something to parse Unary
	if parser.matchKind(lexer.Operator) {
		switch parser.currentToken.DirectValue {
		case lexer.Sub, lexer.Add, lexer.NegateBits, lexer.SignNot, lexer.Not:
			operator := parser.currentToken
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}

			x, parsingError := parser.parseUnaryExpression()
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast.IExpression); !ok {
				return nil, parser.expectingExpressionError(PointerExpression)
			}
			return &ast.UnaryExpression{
				Operator: operator,
				X:        x.(ast.IExpression),
			}, nil
		}
	}
	return parser.parsePrimaryExpression()
}

func (parser *Parser) removeNewLines() error {
	for parser.matchDirectValue(lexer.NewLine) {
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return tokenizingError
		}
	}
	return nil
}

func (parser *Parser) parseLambdaExpression() (*ast.LambdaExpression, error) {
	var arguments []*ast.Identifier
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.Colon) {
			break
		}
		newLinesRemoveError := parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		identifier, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := identifier.(*ast.Identifier); !ok {
			return nil, parser.expectingIdentifier(LambdaExpression)
		}
		arguments = append(arguments, identifier.(*ast.Identifier))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirectValue(lexer.Colon) {
			return nil, parser.newSyntaxError(LambdaExpression)
		}
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.Colon) {
		return nil, parser.newSyntaxError(LambdaExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	code, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := code.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(LambdaExpression)
	}
	return &ast.LambdaExpression{
		Arguments: arguments,
		Code:      code.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseParentheses() (ast.IExpression, error) {
	/*
		This should also parse generators
	*/
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.newSyntaxError(ParenthesesExpression)
	}

	firstExpression, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := firstExpression.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(ParenthesesExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if parser.matchDirectValue(lexer.For) {
		return parser.parseGeneratorExpression(firstExpression.(ast.IExpression))
	}
	if parser.matchDirectValue(lexer.CloseParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.ParenthesesExpression{
			X: firstExpression.(ast.IExpression),
		}, nil
	}
	if !parser.matchDirectValue(lexer.Comma) {
		return nil, parser.newSyntaxError(ParenthesesExpression)
	}
	var values []ast.IExpression
	values = append(values, firstExpression.(ast.IExpression))
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var nextValue ast.Node
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		nextValue, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := nextValue.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(ParenthesesExpression)
		}
		values = append(values, nextValue.(ast.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirectValue(lexer.CloseParentheses) {
			return nil, parser.newSyntaxError(TupleExpression)
		}
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.expressionNeverClosedError(TupleExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.TupleExpression{
		Values: values,
	}, nil
}
func (parser *Parser) parseArrayExpression() (*ast.ArrayExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var values []ast.IExpression
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseSquareBracket) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		value, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := value.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(ArrayExpression)
		}
		values = append(values, value.(ast.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirectValue(lexer.CloseSquareBracket) {
			return nil, parser.newSyntaxError(ArrayExpression)
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ArrayExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseHashExpression() (*ast.HashExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var values []*ast.KeyValue
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseBrace) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		leftHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := leftHandSide.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(HashExpression)
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchDirectValue(lexer.Colon) {
			return nil, parser.newSyntaxError(HashExpression)
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		rightHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(HashExpression)
		}
		values = append(values, &ast.KeyValue{
			Key:   leftHandSide.(ast.IExpression),
			Value: rightHandSide.(ast.IExpression),
		})
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.HashExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseWhileStatement() (*ast.WhileLoopStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, parser.newSyntaxError(WhileStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(WhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var whileChild ast.Node
	var body []ast.Node
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
		whileChild, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, whileChild)
	}
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(WhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.WhileLoopStatement{
		Condition: condition.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseDoWhileStatement() (*ast.DoWhileStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(DoWhileStatement)
	}
	// Parse Body
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirectValue(lexer.While) {
				break
			}
			continue
		}
		bodyNode, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, bodyNode)
	}
	// Parse Condition
	if !parser.matchDirectValue(lexer.While) {
		return nil, parser.newSyntaxError(DoWhileStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	var condition ast.Node
	condition, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(WhileStatement)
	}
	return &ast.DoWhileStatement{
		Condition: condition.(ast.IExpression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseIfStatement() (*ast.IfStatement, error) {
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
		return nil, parser.expectingExpressionError(IfStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(IfStatement)
	}
	// Parse If
	root := &ast.IfStatement{
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
				&ast.IfStatement{
					Condition: elifCondition.(ast.IExpression),
					Body:      elifBody,
				},
			)
			lastCondition = lastCondition.Else[0].(*ast.IfStatement)
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
		return nil, parser.statementNeverEndedError(IfStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return root, nil
}

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

func (parser *Parser) parseSwitchStatement() (*ast.SwitchStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	target, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := target.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(SwitchStatement)
	}
	if !parser.matchDirectValue(lexer.NewLine) {
		return nil, parser.newSyntaxError(SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	// parse Cases
	var caseBlocks []*ast.CaseBlock
	if parser.matchDirectValue(lexer.Case) {
		for parser.hasNext() {
			if parser.matchDirectValue(lexer.Default) ||
				parser.matchDirectValue(lexer.End) {
				break
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			newLinesRemoveError = parser.removeNewLines()
			if newLinesRemoveError != nil {
				return nil, newLinesRemoveError
			}
			var cases []ast.IExpression
			var caseTarget ast.Node
			for parser.hasNext() {
				caseTarget, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				if _, ok := caseTarget.(ast.IExpression); !ok {
					return nil, parser.expectingExpressionError(CaseBlock)
				}
				cases = append(cases, caseTarget.(ast.IExpression))
				if parser.matchDirectValue(lexer.NewLine) {
					break
				} else if parser.matchDirectValue(lexer.Comma) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
				} else {
					return nil, parser.newSyntaxError(CaseBlock)
				}
			}
			if !parser.matchDirectValue(lexer.NewLine) {
				return nil, parser.newSyntaxError(CaseBlock)
			}
			// Targets Body
			var caseBody []ast.Node
			var caseBodyNode ast.Node
			for parser.hasNext() {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirectValue(lexer.Case) ||
						parser.matchDirectValue(lexer.Default) ||
						parser.matchDirectValue(lexer.End) {
						break
					}
					continue
				}
				caseBodyNode, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				caseBody = append(caseBody, caseBodyNode)
			}
			// Targets block
			caseBlocks = append(caseBlocks, &ast.CaseBlock{
				Cases: cases,
				Body:  caseBody,
			})
		}
	}
	// Parse Default
	var defaultBody []ast.Node
	if parser.matchDirectValue(lexer.Default) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirectValue(lexer.NewLine) {
			return nil, parser.newSyntaxError(DefaultBlock)
		}
		var defaultBodyNode ast.Node
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
			defaultBodyNode, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			defaultBody = append(defaultBody, defaultBodyNode)
		}
	}
	// Finally detect valid end
	if !parser.matchDirectValue(lexer.End) {
		return nil, parser.statementNeverEndedError(SwitchStatement)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SwitchStatement{
		Target:     target.(ast.IExpression),
		CaseBlocks: caseBlocks,
		Default:    defaultBody,
	}, nil
}

func (parser *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var results []ast.IExpression
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF) {
			break
		}

		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(ReturnStatement)
		}
		results = append(results, result.(ast.IExpression))
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF)) {
			return nil, parser.newSyntaxError(ReturnStatement)
		}
	}
	return &ast.ReturnStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseYieldStatement() (*ast.YieldStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var results []ast.IExpression
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF) {
			break
		}

		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(YieldStatement)
		}
		results = append(results, result.(ast.IExpression))
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF)) {
			return nil, parser.newSyntaxError(YieldStatement)
		}
	}
	return &ast.YieldStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseContinueStatement() (*ast.ContinueStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ContinueStatement{}, nil
}

func (parser *Parser) parseBreakStatement() (*ast.BreakStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.BreakStatement{}, nil
}

func (parser *Parser) parseRedoStatement() (*ast.RedoStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.RedoStatement{}, nil
}

func (parser *Parser) parsePassStatement() (*ast.PassStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.PassStatement{}, nil
}

func (parser *Parser) parseOperand() (ast.Node, error) {
	switch parser.currentToken.Kind {
	case lexer.Literal, lexer.Boolean, lexer.NoneType:
		return parser.parseLiteral()
	case lexer.IdentifierKind:
		identifier := parser.currentToken

		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.Identifier{
			Token: identifier,
		}, nil
	case lexer.Keyboard:
		switch parser.currentToken.DirectValue {
		case lexer.Lambda:
			return parser.parseLambdaExpression()
		case lexer.While:
			return parser.parseWhileStatement()
		case lexer.For:
			return parser.parseForStatement()
		case lexer.Until:
			return parser.parseUntilStatement()
		case lexer.If:
			return parser.parseIfStatement()
		case lexer.Unless:
			return parser.parseUnlessStatement()
		case lexer.Switch:
			return parser.parseSwitchStatement()
		case lexer.Module:
			return parser.parseModuleStatement()
		case lexer.Def:
			return parser.parseFunctionDefinitionStatement()
		case lexer.Interface:
			return parser.parseInterfaceStatement()
		case lexer.Class:
			return parser.parseClassStatement()
		case lexer.Raise:
			return parser.parseRaiseStatement()
		case lexer.Try:
			return parser.parseTryStatement()
		case lexer.Return:
			return parser.parseReturnStatement()
		case lexer.Yield:
			return parser.parseYieldStatement()
		case lexer.Continue:
			return parser.parseContinueStatement()
		case lexer.Break:
			return parser.parseBreakStatement()
		case lexer.Redo:
			return parser.parseRedoStatement()
		case lexer.Pass:
			return parser.parsePassStatement()
		case lexer.Do:
			return parser.parseDoWhileStatement()
		}
	case lexer.Punctuation:
		switch parser.currentToken.DirectValue {
		case lexer.OpenParentheses:
			return parser.parseParentheses()
		case lexer.OpenSquareBracket: // Parse Arrays
			return parser.parseArrayExpression()
		case lexer.OpenBrace: // Parse Dictionaries
			return parser.parseHashExpression()
		}
	}
	return nil, UnknownToken
}

func (parser *Parser) parseSelectorExpression(expression ast.IExpression) (*ast.SelectorExpression, error) {
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

func (parser *Parser) parseMethodInvocationExpression(expression ast.IExpression) (*ast.MethodInvocationExpression, error) {
	var arguments []ast.IExpression
	// The first token is open parentheses
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.CloseParentheses) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}

		argument, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := argument.(ast.IExpression); !ok {
			return nil, parser.expectingExpressionError(MethodInvocationExpression)
		}
		arguments = append(arguments, argument.(ast.IExpression))
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.MethodInvocationExpression{
		Function:  expression,
		Arguments: arguments,
	}, nil
}

func (parser *Parser) parseIndexExpression(expression ast.IExpression) (*ast.IndexExpression, error) {
	tokenizationError := parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	// var rightIndex ast.Node

	index, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := index.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(IndexExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	if !parser.matchDirectValue(lexer.CloseSquareBracket) {
		return nil, parser.newSyntaxError(IndexExpression)
	}
	tokenizationError = parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	return &ast.IndexExpression{
		Source: expression,
		Index:  index.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseIfOneLinerExpression(result ast.IExpression) (*ast.IfOneLinerExpression, error) {
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
		return nil, parser.expectingExpressionError(IfOneLinerExpression)
	}
	if !parser.matchDirectValue(lexer.Else) {
		return &ast.IfOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast.IExpression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var elseResult ast.Node
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := elseResult.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(OneLineElseBlock)
	}
	return &ast.IfOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast.IExpression),
		ElseResult: elseResult.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseUnlessOneLinerExpression(result ast.IExpression) (*ast.UnlessOneLinerExpression, error) {
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
		return nil, parser.expectingExpressionError(UnlessOneLinerExpression)
	}
	if !parser.matchDirectValue(lexer.Else) {
		return &ast.UnlessOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast.IExpression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var elseResult ast.Node

	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(OneLineElseBlock)
	}
	return &ast.UnlessOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast.IExpression),
		ElseResult: elseResult.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseGeneratorExpression(operation ast.IExpression) (*ast.GeneratorExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	var variables []*ast.Identifier
	numberOfVariables := 0
	for parser.hasNext() {
		if parser.matchDirectValue(lexer.In) {
			break
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if !parser.matchKind(lexer.IdentifierKind) {
			return nil, parser.newSyntaxError(GeneratorExpression)
		}
		variables = append(variables, &ast.Identifier{
			Token: parser.currentToken,
		})
		numberOfVariables++
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		newLinesRemoveError = parser.removeNewLines()
		if newLinesRemoveError != nil {
			return nil, newLinesRemoveError
		}
		if parser.matchDirectValue(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
	}
	if numberOfVariables == 0 {
		return nil, parser.newSyntaxError(GeneratorExpression)
	}
	if !parser.matchDirectValue(lexer.In) {
		return nil, parser.newSyntaxError(GeneratorExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(GeneratorExpression)
	}
	newLinesRemoveError = parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}
	// Finally detect the closing parentheses
	if !parser.matchDirectValue(lexer.CloseParentheses) {
		return nil, parser.newSyntaxError(GeneratorExpression)
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.GeneratorExpression{
		Operation: operation,
		Receivers: variables,
		Source:    source.(ast.IExpression),
	}, nil
}

func (parser *Parser) parseAssignmentStatement(leftHandSide ast.IExpression) (*ast.AssignStatement, error) {
	assignmentToken := parser.currentToken
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	newLinesRemoveError := parser.removeNewLines()
	if newLinesRemoveError != nil {
		return nil, newLinesRemoveError
	}

	rightHandSide, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := rightHandSide.(ast.IExpression); !ok {
		return nil, parser.expectingExpressionError(AssignStatement)
	}
	return &ast.AssignStatement{
		LeftHandSide:   leftHandSide,
		AssignOperator: assignmentToken,
		RightHandSide:  rightHandSide.(ast.IExpression),
	}, nil
}

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

func (parser *Parser) Parse() (*ast.Program, error) {
	result := &ast.Program{
		Begin: nil,
		End:   nil,
		Body:  nil,
	}
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var beginStatement *ast.BeginStatement
	var endStatement *ast.EndStatement
	var parsedExpression ast.Node
	var parsingError error
	for parser.hasNext() {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		switch {
		case parser.matchDirectValue(lexer.BEGIN):
			if result.Begin != nil {
				return nil, BeginRepeated
			}
			beginStatement, parsingError = parser.parseBeginStatement()
			if parsingError != nil {
				return nil, parsingError
			}
			result.Begin = beginStatement
		case parser.matchDirectValue(lexer.END):
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

func NewParser(lexer_ *lexer.Lexer) *Parser {
	return &Parser{
		lineStack:    common.Stack[int]{},
		lexer:        lexer_,
		complete:     false,
		currentToken: nil,
	}
}
