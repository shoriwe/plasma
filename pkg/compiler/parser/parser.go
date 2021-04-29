package parser

import (
	"errors"
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
)

type Parser struct {
	lexer        *lexer.Lexer
	complete     bool
	currentToken *lexer.Token
}

func (parser *Parser) hasNext() bool {
	return !parser.complete
}

func (parser *Parser) next() error {
	token, tokenizingError := parser.lexer.Next()
	if tokenizingError != nil {
		return tokenizingError
	}
	if token.Kind == lexer.EOF {
		parser.complete = true
	}
	parser.currentToken = token
	return nil
}

func (parser *Parser) matchDirect(directValue int) bool {
	if parser.currentToken == nil {
		return false
	}
	return parser.currentToken.DirectValue == directValue
}

func (parser *Parser) matchKind(kind int) bool {
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
	return parser.currentToken.String == value
}

func (parser *Parser) updateState() {
	parser.complete = true
}

func (parser *Parser) parseForLoop() (*ast.ForLoopStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var receivers []*ast.Identifier
	for ; !parser.complete; {
		if parser.matchDirect(lexer.In) {
			break
		} else if !parser.matchKind(lexer.IdentifierKind) {
			return nil, errors.New(fmt.Sprintf("invalid definition of for loop at line %d", parser.currentLine()))
		}
		receivers = append(receivers, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.matchDirect(lexer.In) {
			break
		} else {
			return nil, errors.New(fmt.Sprintf("invalid declaration of for loop at line %d", parser.currentLine()))
		}
	}
	if !parser.matchDirect(lexer.In) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of for loop at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("recevied a non expression as for loop source at line %d", parser.currentLine()))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of for loop at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var body []ast.Node
	var bodyNode ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("for loop never ended at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.ForLoopStatement{
		Receivers: receivers,
		Source:    source.(ast.Expression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseUntilLoop() (*ast.UntilLoopStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("invalid until loop declaration at line %d", parser.currentLine()))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid until loop declaration at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var untilChild ast.Node
	var body []ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
				break
			}
			continue
		}
		untilChild, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		body = append(body, untilChild)
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("until statement never closed at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.UntilLoopStatement{
		Condition: condition.(ast.Expression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseModuleStatement() (*ast.ModuleStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of module at line %d", parser.currentLine()))
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of module at line %d", parser.currentLine()))
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("Module declaration never closed at line %d", parser.currentLine()))
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
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid Function definition at line %d", parser.currentLine()))
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.OpenParentheses) {
		return nil, errors.New(fmt.Sprintf("invalid Function definition at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []*ast.Identifier
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		if !parser.matchKind(lexer.IdentifierKind) {
			return nil, errors.New(fmt.Sprintf("invalid Function definition at line %d", parser.currentLine()))
		}
		argument := &ast.Identifier{
			Token: parser.currentToken,
		}
		arguments = append(arguments, argument)
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if parser.matchDirect(lexer.CloseParentheses) {
			break
		} else {
			return nil, errors.New(fmt.Sprintf("invalid Function definition at line %d", parser.currentLine()))
		}
	}
	if !parser.matchDirect(lexer.CloseParentheses) {
		return nil, errors.New(fmt.Sprintf("invalid Function definition at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid Function definition at line %d", parser.currentLine()))
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("function declaration never closed at line %d", parser.currentLine()))
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

func (parser *Parser) parseAsyncFunctionDefinitionStatement() (*ast.AsyncFunctionDefinitionStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	functionDefinition, parsingError := parser.parseFunctionDefinitionStatement()
	if parsingError != nil {
		return nil, parsingError
	}
	return &ast.AsyncFunctionDefinitionStatement{
		Name:      functionDefinition.Name,
		Arguments: functionDefinition.Arguments,
		Body:      functionDefinition.Body,
	}, nil
}

func (parser *Parser) parseClassStatement() (*ast.ClassStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of class at line %d", parser.currentLine()))
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.Expression
	var base ast.Node
	var parsingError error
	if parser.matchDirect(lexer.OpenParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		for ; !parser.complete; {
			base, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := base.(*ast.Identifier); !ok {
				if _, ok2 := base.(*ast.SelectorExpression); !ok2 {
					return nil, errors.New(fmt.Sprintf("received an invalid expression as class base at line %d", parser.currentLine()))
				}
			}
			bases = append(bases, base.(ast.Expression))
			if parser.matchDirect(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.matchDirect(lexer.CloseParentheses) {
				break
			} else {
				return nil, errors.New(fmt.Sprintf("invalid definition of a class at line %d", parser.currentLine()))
			}
		}
		if !parser.matchDirect(lexer.CloseParentheses) {
			return nil, errors.New(fmt.Sprintf("invalid declaration of class bases at line %d", parser.currentLine()))
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of class at line %d", parser.currentLine()))
	}
	var body []ast.Node
	var bodyNode ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("class declaration never closed at line %d", parser.currentLine()))
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
	if _, ok := x.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("raise statement must receive an expression not an statement at line %d", parser.currentLine()))
	}
	return &ast.RaiseStatement{
		X: x.(ast.Expression),
	}, nil
}
func (parser *Parser) parseTryStatement() (*ast.TryStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid try statement definition at line %d", parser.currentLine()))
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) ||
				parser.matchDirect(lexer.Except) ||
				parser.matchDirect(lexer.Else) ||
				parser.matchDirect(lexer.Finally) {
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
	for ; !parser.complete; {
		if !parser.matchDirect(lexer.Except) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var targets []ast.Expression
		var target ast.Node
		for ; !parser.complete; {
			if parser.matchDirect(lexer.NewLine) ||
				parser.matchDirect(lexer.As) {
				break
			}
			target, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := target.(ast.Expression); !ok {
				return nil, errors.New(fmt.Sprintf("invalid definition of try statement exception at line %d", parser.currentLine()))
			}
			targets = append(targets, target.(ast.Expression))
			if parser.matchDirect(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			}
		}
		var captureName *ast.Identifier
		if parser.matchDirect(lexer.As) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if !parser.matchKind(lexer.IdentifierKind) {
				return nil, errors.New(fmt.Sprintf("invalid declaratione of except block at line %d", parser.currentLine()))
			}
			captureName = &ast.Identifier{
				Token: parser.currentToken,
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		}
		if !parser.matchDirect(lexer.NewLine) {
			return nil, errors.New(fmt.Sprintf("invalid declaratione of except block at line %d", parser.currentLine()))
		}
		var exceptBody []ast.Node
		var exceptBodyNode ast.Node
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.End) ||
					parser.matchDirect(lexer.Except) ||
					parser.matchDirect(lexer.Else) ||
					parser.matchDirect(lexer.Finally) {
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
	if parser.matchDirect(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirect(lexer.NewLine) {
			return nil, errors.New(fmt.Sprintf("invalid definition of else block in try statement at line %d", parser.currentLine()))
		}
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.End) ||
					parser.matchDirect(lexer.Finally) {
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
	if parser.matchDirect(lexer.Finally) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirect(lexer.NewLine) {
			return nil, errors.New(fmt.Sprintf("invalid definition of finally block in try statement at line %d", parser.currentLine()))
		}
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("invalid declaratione of except block at line %d", parser.currentLine()))
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
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid definition of begin statement at line %d", parser.currentLine()))
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("invalid definition of begin statement at line %d", parser.currentLine()))
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
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid definition of begin statement at line %d", parser.currentLine()))
	}
	var body []ast.Node
	var bodyNode ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("invalid definition of begin statement at line %d", parser.currentLine()))
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
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of class at line %d", parser.currentLine()))
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var bases []ast.Expression
	var base ast.Node
	var parsingError error
	if parser.matchDirect(lexer.OpenParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		for ; !parser.complete; {
			base, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := base.(*ast.Identifier); !ok {
				if _, ok2 := base.(*ast.SelectorExpression); !ok2 {
					return nil, errors.New(fmt.Sprintf("received an invalid expression as class base at line %d", parser.currentLine()))
				}
			}
			bases = append(bases, base.(ast.Expression))
			if parser.matchDirect(lexer.Comma) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
			} else if parser.matchDirect(lexer.CloseParentheses) {
				break
			} else {
				return nil, errors.New(fmt.Sprintf("invalid definition of a class at line %d", parser.currentLine()))
			}
		}
		if !parser.matchDirect(lexer.CloseParentheses) {
			return nil, errors.New(fmt.Sprintf("invalid declaration of class bases at line %d", parser.currentLine()))
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of class at line %d", parser.currentLine()))
	}
	var methods []*ast.FunctionDefinitionStatement
	var asyncMethods []*ast.AsyncFunctionDefinitionStatement
	var node ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
				break
			}
			continue
		}
		node, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := node.(*ast.FunctionDefinitionStatement); !ok {
			if _, ok2 := node.(*ast.AsyncFunctionDefinitionStatement); !ok2 {
				return nil, errors.New(fmt.Sprintf("received a non function definition in interface at line %d", parser.currentLine()))
			}
			asyncMethods = append(asyncMethods, node.(*ast.AsyncFunctionDefinitionStatement))
		} else {
			methods = append(methods, node.(*ast.FunctionDefinitionStatement))
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("class declaration never closed at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.InterfaceStatement{
		Name:                   name,
		Bases:                  bases,
		MethodDefinitions:      methods,
		AsyncMethodDefinitions: asyncMethods,
	}, nil
}

func (parser *Parser) parseGoToStatement() (*ast.GoToStatement, error) {
	line := parser.currentLine()
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	identifier, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := identifier.(*ast.Identifier); !ok {
		return nil, errors.New(fmt.Sprintf("goto statement must receive an identifier at line %d", line))
	}
	return &ast.GoToStatement{
		Name: identifier.(*ast.Identifier),
	}, nil
}

func (parser *Parser) parseLiteral() (ast.Expression, error) {
	if !parser.matchKind(lexer.Literal) &&
		!parser.matchKind(lexer.Boolean) &&
		!parser.matchKind(lexer.NoneType) {
		return nil, errors.New(fmt.Sprintf("invalid kind of token %s at line %d", parser.currentToken.String, parser.currentLine()))
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
	return nil, errors.New(fmt.Sprintf("could not determine the directValue of token %s at line %d", parser.currentToken.String, parser.currentLine()))
}

func (parser *Parser) parseBinaryExpression(precedence int) (ast.Node, error) {
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
	for ; !parser.complete; {
		if !parser.matchKind(lexer.Operator) &&
			!parser.matchKind(lexer.Comparator) {
			break
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
		line := parser.currentLine()
		rightHandSide, parsingError = parser.parseBinaryExpression(operatorPrecedence + 1)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression child for binary expression at line %d", line))
		}
		leftHandSide = &ast.BinaryExpression{
			LeftHandSide:  leftHandSide.(ast.Expression),
			Operator:      operator,
			RightHandSide: rightHandSide.(ast.Expression),
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
			line := parser.currentLine()
			x, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast.Expression); !ok {
				return nil, errors.New(fmt.Sprintf("received a non expression child for unary expression at line %d", line))
			}
			return &ast.UnaryExpression{
				Operator: operator,
				X:        x.(ast.Expression),
			}, nil
		case lexer.BitWiseAnd:
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			line := parser.currentLine()
			x, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast.Expression); !ok {
				return nil, errors.New(fmt.Sprintf("received a non expression child for pointer expression at line %d", line))
			}
			return &ast.PointerExpression{
				X: x.(ast.Expression),
			}, nil
		case lexer.Star:
			tokenizingError := parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			line := parser.currentLine()
			x, parsingError := parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := x.(ast.Expression); !ok {
				return nil, errors.New(fmt.Sprintf("received a non expression child for star expression at line %d", line))
			}
			return &ast.StarExpression{
				X: x.(ast.Expression),
			}, nil
		}
	} else if parser.matchKind(lexer.AwaitKeyboard) {
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		line := parser.currentLine()
		x, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := x.(*ast.MethodInvocationExpression); !ok {
			return nil, errors.New(fmt.Sprintf("await must receive a method invocation at line %d", line))
		}
		return &ast.AwaitExpression{
			X: x.(*ast.MethodInvocationExpression),
		}, nil
	}
	return parser.parsePrimaryExpression()
}

func (parser *Parser) parseLambdaExpression() (*ast.LambdaExpression, error) {
	var arguments []*ast.Identifier
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for ; !parser.complete; {
		if parser.matchDirect(lexer.Colon) {
			break
		}
		line := parser.currentLine()
		identifier, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := identifier.(*ast.Identifier); !ok {
			return nil, errors.New(fmt.Sprintf("recevied a non identifier value in lambda arguments at line %d", line))
		}
		arguments = append(arguments, identifier.(*ast.Identifier))
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirect(lexer.Colon) {
			return nil, errors.New(fmt.Sprintf("invalid lambda definition at line %d", parser.currentLine()))
		}
	}
	if !parser.matchDirect(lexer.Colon) {
		return nil, errors.New(fmt.Sprintf("invalid lambda definition at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	code, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := code.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("recevied a non expression in lambda body at line %d", line))
	}
	return &ast.LambdaExpression{
		Arguments: arguments,
		Code:      code.(ast.Expression),
	}, nil
}

func (parser *Parser) parseParentheses() (ast.Expression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if parser.matchDirect(lexer.CloseParentheses) {
		return nil, errors.New(fmt.Sprintf("syntax error: empty parentheses expression at line %d", parser.currentLine()))
	}
	line := parser.currentLine()
	firstExpression, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := firstExpression.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression as parentheses body at line %d", line))
	}
	if parser.matchDirect(lexer.CloseParentheses) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		return &ast.ParenthesesExpression{
			X: firstExpression.(ast.Expression),
		}, nil
	}
	if !parser.matchDirect(lexer.Comma) {
		return nil, errors.New(fmt.Sprintf("syntax error: empty parentheses expression at line %d", parser.currentLine()))
	}
	var values []ast.Expression
	values = append(values, firstExpression.(ast.Expression))
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var nextValue ast.Node
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		line = parser.currentLine()
		nextValue, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := nextValue.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as parentheses body at line %d", line))
		}
		values = append(values, nextValue.(ast.Expression))
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirect(lexer.CloseParentheses) {
			return nil, errors.New(fmt.Sprintf("syntax error: invalid tuple definition line %d", parser.currentLine()))
		}
	}
	if !parser.matchDirect(lexer.CloseParentheses) {
		return nil, errors.New(fmt.Sprintf("syntax error: tuple expression never closed%d", parser.currentLine()))
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
	var values []ast.Expression
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseSquareBracket) {
			break
		}
		line := parser.currentLine()
		value, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := value.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as array element at line %d", line))
		}
		values = append(values, value.(ast.Expression))
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !parser.matchDirect(lexer.CloseSquareBracket) {
			return nil, errors.New(fmt.Sprintf("invalid array definition at line %d", parser.currentLine()))
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
	var values []*ast.KeyValue
	var leftHandSide ast.Node
	var rightHandSide ast.Node
	var parsingError error
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseBrace) {
			break
		}
		line := parser.currentLine()
		leftHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := leftHandSide.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as a key in hash expression at line %d", line))
		}
		if !parser.matchDirect(lexer.Colon) {
			return nil, errors.New(fmt.Sprintf("syntax error: invalid hash definition at line %d", parser.currentLine()))
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		line = parser.currentLine()
		rightHandSide, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := rightHandSide.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as a value in hash expression at line %d", line))
		}
		values = append(values, &ast.KeyValue{
			Key:   leftHandSide.(ast.Expression),
			Value: rightHandSide.(ast.Expression),
		})
		if parser.matchDirect(lexer.Comma) {
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
	return &ast.HashExpression{
		Values: values,
	}, nil
}

func (parser *Parser) parseWhileLoop() (*ast.WhileLoopStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("invalid while loop declaration at line %d", parser.currentLine()))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid while loop declaration at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var whileChild ast.Node
	var body []ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("while statement never closed at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.WhileLoopStatement{
		Condition: condition.(ast.Expression),
		Body:      body,
	}, nil
}

func (parser *Parser) parseIfStatement() (*ast.IfStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("recevied a non expression condition for if statement at line %d", line))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid if statement declaration at line %d", parser.currentLine()))
	}
	// Parse If
	var body []ast.Node
	var bodyNode ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.Elif) ||
				parser.matchDirect(lexer.Else) ||
				parser.matchDirect(lexer.End) {
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
	// Parse Elifs
	var elifBlocks []*ast.ElifBlock
	if parser.matchDirect(lexer.Elif) {
	parsingElifLoop:
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.Else) ||
					parser.matchDirect(lexer.End) {
					break
				}
				continue
			}
			if !parser.matchDirect(lexer.Elif) {
				return nil, errors.New(fmt.Sprintf("invalid declaration of elif block at line %d", parser.currentLine()))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var elifCondition ast.Node
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := elifCondition.(ast.Expression); !ok {
				return nil, errors.New(fmt.Sprintf("invalid declaration of elif block at line %d", parser.currentLine()))
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid declaration of elif block at line %d", parser.currentLine()))
			}
			var elifBody []ast.Node
			var elifBodyNode ast.Node
			for ; !parser.complete; {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirect(lexer.Else) ||
						parser.matchDirect(lexer.End) {
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
			elifBlocks = append(elifBlocks, &ast.ElifBlock{
				Condition: elifCondition.(ast.Expression),
				Body:      elifBody,
			})
			if parser.matchDirect(lexer.Else) ||
				parser.matchDirect(lexer.End) {
				break parsingElifLoop
			}
		}
	}
	// Parse Default
	var elseBody []ast.Node
	if parser.matchDirect(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast.Node
		if !parser.matchDirect(lexer.NewLine) {
			return nil, errors.New(fmt.Sprintf("invalid definition of else block at line %d", parser.currentLine()))
		}
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("never closed if statement at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.IfStatement{
		Condition:  condition.(ast.Expression),
		Body:       body,
		ElifBlocks: elifBlocks,
		Else:       elseBody,
	}, nil
}

func (parser *Parser) parseUnlessStatement() (*ast.UnlessStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("recevied a non expression condition for if statement at line %d", line))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid if statement declaration at line %d", parser.currentLine()))
	}
	// Parse If
	var body []ast.Node
	var bodyNode ast.Node
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			if parser.matchDirect(lexer.Elif) ||
				parser.matchDirect(lexer.Else) ||
				parser.matchDirect(lexer.End) {
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
	// Parse Elifs
	var elifBlocks []*ast.ElifBlock
	if parser.matchDirect(lexer.Elif) {
	parsingElifLoop:
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.Else) ||
					parser.matchDirect(lexer.End) {
					break
				}
				continue
			}
			if !parser.matchDirect(lexer.Elif) {
				return nil, errors.New(fmt.Sprintf("invalid declaration of elif block at line %d", parser.currentLine()))
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var elifCondition ast.Node
			elifCondition, parsingError = parser.parseBinaryExpression(0)
			if parsingError != nil {
				return nil, parsingError
			}
			if _, ok := elifCondition.(ast.Expression); !ok {
				return nil, errors.New(fmt.Sprintf("invalid declaration of elif block at line %d", parser.currentLine()))
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid declaration of elif block at line %d", parser.currentLine()))
			}
			var elifBody []ast.Node
			var elifBodyNode ast.Node
			for ; !parser.complete; {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirect(lexer.Else) ||
						parser.matchDirect(lexer.End) {
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
			elifBlocks = append(elifBlocks, &ast.ElifBlock{
				Condition: elifCondition.(ast.Expression),
				Body:      elifBody,
			})
			if parser.matchDirect(lexer.Else) ||
				parser.matchDirect(lexer.End) {
				break parsingElifLoop
			}
		}
	}
	// Parse Default
	var elseBody []ast.Node
	if parser.matchDirect(lexer.Else) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		var elseBodyNode ast.Node
		if !parser.matchDirect(lexer.NewLine) {
			return nil, errors.New(fmt.Sprintf("invalid definition of else block at line %d", parser.currentLine()))
		}
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("never closed if statement at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.UnlessStatement{
		Condition:  condition.(ast.Expression),
		Body:       body,
		ElifBlocks: elifBlocks,
		Else:       elseBody,
	}, nil
}

/*
ToDo: Totally nasty, need to handle it like the implementation of the try statement
*/
func (parser *Parser) parseSwitchStatement() (*ast.SwitchStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	target, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := target.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression as switch target at line %d", line))
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid switch statement at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	// parse Cases
	var caseBlocks []*ast.CaseBlock
	if parser.matchDirect(lexer.Case) {
		for ; !parser.complete; {
			if parser.matchDirect(lexer.Default) ||
				parser.matchDirect(lexer.End) {
				break
			}
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			var cases []ast.Expression
			var caseTarget ast.Node
			for ; !parser.complete; {
				caseTarget, parsingError = parser.parseBinaryExpression(0)
				if parsingError != nil {
					return nil, parsingError
				}
				if _, ok := caseTarget.(ast.Expression); !ok {
					return nil, errors.New(fmt.Sprintf("recevied a non expression as case target at line %d", parser.currentLine()))
				}
				cases = append(cases, caseTarget.(ast.Expression))
				if parser.matchDirect(lexer.NewLine) {
					break
				} else if parser.matchDirect(lexer.Comma) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
				} else {
					return nil, errors.New(fmt.Sprintf("invalid case block definition at line %d", parser.currentLine()))
				}
			}
			if !parser.matchDirect(lexer.NewLine) {
				return nil, errors.New(fmt.Sprintf("invalid case block definition at line %d", parser.currentLine()))
			}
			// Case Body
			var caseBody []ast.Node
			var caseBodyNode ast.Node
			for ; !parser.complete; {
				if parser.matchKind(lexer.Separator) {
					tokenizingError = parser.next()
					if tokenizingError != nil {
						return nil, tokenizingError
					}
					if parser.matchDirect(lexer.Case) ||
						parser.matchDirect(lexer.Default) ||
						parser.matchDirect(lexer.End) {
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
			// Case block
			caseBlocks = append(caseBlocks, &ast.CaseBlock{
				Cases: cases,
				Body:  caseBody,
			})
		}
	}
	// Parse Default
	var defaultBody []ast.Node
	if parser.matchDirect(lexer.Default) {
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		if !parser.matchDirect(lexer.NewLine) {
			return nil, errors.New(fmt.Sprintf("invalid definition of default block at line %d", parser.currentLine()))
		}
		var defaultBodyNode ast.Node
		for ; !parser.complete; {
			if parser.matchKind(lexer.Separator) {
				tokenizingError = parser.next()
				if tokenizingError != nil {
					return nil, tokenizingError
				}
				if parser.matchDirect(lexer.End) {
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
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("Switch declaration never ended at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SwitchStatement{
		Target:     target.(ast.Expression),
		CaseBlocks: caseBlocks,
		Default:    defaultBody,
	}, nil
}

func (parser *Parser) parseStructStatement() (*ast.StructStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentLine()))
	}
	name := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var fields []*ast.Identifier
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		} else if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		} else if !parser.matchKind(lexer.IdentifierKind) {
			return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentLine()))
		}
		fields = append(fields, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("invalid struct definition at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.StructStatement{
		Name:   name,
		Fields: fields,
	}, nil
}

func (parser *Parser) parseDeferStatement() (*ast.DeferStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	methodInvocation, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	switch methodInvocation.(type) {
	case *ast.MethodInvocationExpression:
		return &ast.DeferStatement{
			X: methodInvocation.(*ast.MethodInvocationExpression),
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("no function call passed to go statement at line %d", parser.currentLine()))
	}
}

func (parser *Parser) parseGoStatement() (*ast.GoStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	methodInvocation, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	switch methodInvocation.(type) {
	case *ast.MethodInvocationExpression:
		return &ast.GoStatement{
			X: methodInvocation.(*ast.MethodInvocationExpression),
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("no function call passed to go statement at line %d", parser.currentLine()))
	}
}

func (parser *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var results []ast.Expression
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF) {
			break
		}
		line := parser.currentLine()
		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as return value at line %d", line))
		}
		results = append(results, result.(ast.Expression))
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF)) {
			return nil, errors.New(fmt.Sprintf("invalid return statement at line %d", parser.currentLine()))
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
	var results []ast.Expression
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF) {
			break
		}
		line := parser.currentLine()
		result, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := result.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as return value at line %d", line))
		}
		results = append(results, result.(ast.Expression))
		if parser.matchDirect(lexer.Comma) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
		} else if !(parser.matchKind(lexer.Separator) || parser.matchKind(lexer.EOF)) {
			return nil, errors.New(fmt.Sprintf("invalid return statement at line %d", parser.currentLine()))
		}
	}
	return &ast.YieldStatement{
		Results: results,
	}, nil
}

func (parser *Parser) parseSuperStatement() (*ast.SuperInvocationStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.OpenParentheses) {
		return nil, errors.New(fmt.Sprintf("invalid super call at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var arguments []ast.Expression
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		line := parser.currentLine()
		argument, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := argument.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression for super as argument at line %d", line))
		}
		arguments = append(arguments, argument.(ast.Expression))
		if !parser.matchDirect(lexer.Comma) {
			break
		}
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.SuperInvocationStatement{
		Arguments: arguments,
	}, nil
}

func (parser *Parser) parseRetryStatement() (*ast.RetryStatement, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.RetryStatement{}, nil
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

func (parser *Parser) parseEnumStatement() (*ast.EnumStatement, error) { // What about initializing it's identifiers with an specific value?
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchKind(lexer.IdentifierKind) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of enum statement at line %d", parser.currentLine()))
	}
	namespace := &ast.Identifier{
		Token: parser.currentToken,
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if !parser.matchDirect(lexer.NewLine) {
		return nil, errors.New(fmt.Sprintf("invalid declaration of enum statement at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var identifiers []*ast.Identifier
	for ; !parser.complete; {
		if parser.matchDirect(lexer.End) {
			break
		} else if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		} else if !parser.matchKind(lexer.IdentifierKind) {
			return nil, errors.New(fmt.Sprintf("invalid declaration of enum statement at line %d", parser.currentLine()))
		}
		identifiers = append(identifiers, &ast.Identifier{
			Token: parser.currentToken,
		})
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if !parser.matchDirect(lexer.End) {
		return nil, errors.New(fmt.Sprintf("enum never ended at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	return &ast.EnumStatement{
		Name:            namespace,
		EnumIdentifiers: identifiers,
	}, nil
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
			return parser.parseWhileLoop()
		case lexer.For:
			return parser.parseForLoop()
		case lexer.Until:
			return parser.parseUntilLoop()
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
		case lexer.Async:
			return parser.parseAsyncFunctionDefinitionStatement()
		case lexer.Struct:
			return parser.parseStructStatement()
		case lexer.Interface:
			return parser.parseInterfaceStatement()
		case lexer.Defer:
			return parser.parseDeferStatement()
		case lexer.Class:
			return parser.parseClassStatement()
		case lexer.Raise:
			return parser.parseRaiseStatement()
		case lexer.Try:
			return parser.parseTryStatement()
		case lexer.Go:
			return parser.parseGoStatement()
		case lexer.Return:
			return parser.parseReturnStatement()
		case lexer.Yield:
			return parser.parseYieldStatement()
		case lexer.Super:
			return parser.parseSuperStatement()
		case lexer.Retry:
			return parser.parseRetryStatement()
		case lexer.Break:
			return parser.parseBreakStatement()
		case lexer.Redo:
			return parser.parseRedoStatement()
		case lexer.Pass:
			return parser.parsePassStatement()
		case lexer.Enum:
			return parser.parseEnumStatement()
		case lexer.GoTo:
			return parser.parseGoToStatement()
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
	return nil, errors.New(fmt.Sprintf("unknown expression with token at line %d", parser.currentLine()))
}

func (parser *Parser) parseSelectorExpression(expression ast.Expression) (*ast.SelectorExpression, error) {
	selector := expression
	for ; !parser.complete; {
		if !parser.matchDirect(lexer.Dot) {
			break
		}
		tokenizingError := parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
		identifier := parser.currentToken
		if identifier.Kind != lexer.IdentifierKind {
			return nil, errors.New(fmt.Sprintf("invalid selector at token in line %d", identifier.Line))
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

func (parser *Parser) parseMethodInvocationExpression(expression ast.Expression) (*ast.MethodInvocationExpression, error) {
	var arguments []ast.Expression
	// The first token is open parentheses
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	for ; !parser.complete; {
		if parser.matchDirect(lexer.CloseParentheses) {
			break
		}
		line := parser.currentLine()
		argument, parsingError := parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := argument.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression as method call at line %d", line))
		}
		arguments = append(arguments, argument.(ast.Expression))
		if parser.matchDirect(lexer.Comma) {
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

func (parser *Parser) parseIndexExpression(expression ast.Expression) (*ast.IndexExpression, error) {
	tokenizationError := parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	var rightIndex ast.Node
	line := parser.currentLine()
	leftIndex, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := leftIndex.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression index at line %d", line))
	}
	if parser.matchDirect(lexer.Colon) {
		tokenizationError = parser.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		line = parser.currentLine()
		rightIndex, parsingError = parser.parseBinaryExpression(0)
		if parsingError != nil {
			return nil, parsingError
		}
		if _, ok := leftIndex.(ast.Expression); !ok {
			return nil, errors.New(fmt.Sprintf("received a non expression index at line %d", line))
		}
	}
	tokenizationError = parser.next()
	if tokenizationError != nil {
		return nil, tokenizationError
	}
	if rightIndex == nil {
		return &ast.IndexExpression{
			Source: expression,
			Index: [2]ast.Expression{
				leftIndex.(ast.Expression),
				nil,
			},
		}, nil
	}
	return &ast.IndexExpression{
		Source: expression,
		Index: [2]ast.Expression{
			leftIndex.(ast.Expression),
			rightIndex.(ast.Expression),
		},
	}, nil
}

func (parser *Parser) parseIfOneLiner(result ast.Expression) (*ast.IfOneLineExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression as one liner if expression condition at line %d", line))
	}
	if !parser.matchDirect(lexer.Else) {
		return &ast.IfOneLineExpression{
			Result:     result,
			Condition:  condition.(ast.Expression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var elseResult ast.Node
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := elseResult.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression as else result at line %d", line))
	}
	return &ast.IfOneLineExpression{
		Result:     result,
		Condition:  condition.(ast.Expression),
		ElseResult: elseResult.(ast.Expression),
	}, nil
}

func (parser *Parser) parseUnlessOneLiner(result ast.Expression) (*ast.UnlessOneLinerExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	condition, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression for unless condition at line %d", line))
	}
	if !parser.matchDirect(lexer.Else) {
		return &ast.UnlessOneLinerExpression{
			Result:     result,
			Condition:  condition.(ast.Expression),
			ElseResult: nil,
		}, nil
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var elseResult ast.Node
	line = parser.currentLine()
	elseResult, parsingError = parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := condition.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression output for else at line %d", line))
	}
	return &ast.UnlessOneLinerExpression{
		Result:     result,
		Condition:  condition.(ast.Expression),
		ElseResult: elseResult.(ast.Expression),
	}, nil
}

func (parser *Parser) parseGeneratorExpression(operation ast.Expression) (*ast.GeneratorExpression, error) {
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	var variables []*ast.Identifier
	numberOfVariables := 0
	for ; !parser.complete; {
		if parser.matchDirect(lexer.In) {
			break
		}
		variables = append(variables, &ast.Identifier{
			Token: parser.currentToken,
		})
		numberOfVariables++
		tokenizingError = parser.next()
		if tokenizingError != nil {
			return nil, tokenizingError
		}
	}
	if numberOfVariables == 0 {
		return nil, errors.New(fmt.Sprintf("syntax error: no receivers in generator defined at line %d", parser.currentLine()))
	}
	if !parser.matchDirect(lexer.In) {
		return nil, errors.New(fmt.Sprintf("syntax error: invalid generator syntax at line %d", parser.currentLine()))
	}
	tokenizingError = parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	source, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := source.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression as source in generator at line %d", line))
	}
	return &ast.GeneratorExpression{
		Operation: operation,
		Variables: variables,
		Source:    source.(ast.Expression),
	}, nil
}

func (parser *Parser) parseAssignmentStatement(leftHandSide ast.Expression) (*ast.AssignStatement, error) {
	assignmentToken := parser.currentToken
	tokenizingError := parser.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	line := parser.currentLine()
	rightHandSide, parsingError := parser.parseBinaryExpression(0)
	if parsingError != nil {
		return nil, parsingError
	}
	if _, ok := rightHandSide.(ast.Expression); !ok {
		return nil, errors.New(fmt.Sprintf("received a non expression as right hand side of assign statement at line %d", line))
	}
	return &ast.AssignStatement{
		LeftHandSide:   leftHandSide,
		AssignOperator: assignmentToken,
		RightHandSide:  rightHandSide.(ast.Expression),
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
			parsedNode, parsingError = parser.parseSelectorExpression(parsedNode.(ast.Expression))
		case lexer.OpenParentheses: // Is function Call
			parsedNode, parsingError = parser.parseMethodInvocationExpression(parsedNode.(ast.Expression))
		case lexer.OpenSquareBracket: // Is indexing
			parsedNode, parsingError = parser.parseIndexExpression(parsedNode.(ast.Expression))
		case lexer.For: // Generators
			parsedNode, parsingError = parser.parseGeneratorExpression(parsedNode.(ast.Expression))
		case lexer.If: // One line If
			parsedNode, parsingError = parser.parseIfOneLiner(parsedNode.(ast.Expression))
		case lexer.Unless: // One line Unless
			parsedNode, parsingError = parser.parseUnlessOneLiner(parsedNode.(ast.Expression))
		default:
			if parser.matchKind(lexer.Assignment) {
				parsedNode, parsingError = parser.parseAssignmentStatement(parsedNode.(ast.Expression))
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
	for ; !parser.complete; {
		if parser.matchKind(lexer.Separator) {
			tokenizingError = parser.next()
			if tokenizingError != nil {
				return nil, tokenizingError
			}
			continue
		}
		if parser.matchDirect(lexer.BEGIN) {
			beginStatement, parsingError = parser.parseBeginStatement()
			if result.Begin != nil {
				return nil, errors.New("multiple declarations of BEGIN statement at line")
			}
			result.Begin = beginStatement
		} else if parser.matchDirect(lexer.END) {
			endStatement, parsingError = parser.parseEndStatement()
			if result.End != nil {
				return nil, errors.New("multiple declarations of END statement at line")
			}
			result.End = endStatement
		} else {
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

func NewParser(lexer_ *lexer.Lexer) (*Parser, error) {
	parser := &Parser{
		lexer:        lexer_,
		complete:     false,
		currentToken: nil,
	}
	return parser, nil
}
