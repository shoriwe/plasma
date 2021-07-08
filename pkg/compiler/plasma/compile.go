package plasma

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler/ast"
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
	"github.com/shoriwe/gplasma/pkg/compiler/parser"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/tools"
	"github.com/shoriwe/gplasma/pkg/vm"
	"reflect"
	"strconv"
	"strings"
)

/*
	Compile to the Plasma stack VM
*/

type Options struct {
	Debug bool
}

type Compiler struct {
	parser  *parser.Parser
	options Options
}

func (c *Compiler) compileBegin(begin *ast.BeginStatement) ([]vm.Code, *errors.Error) {
	return c.compileBody(begin.Body)
}

func (c *Compiler) compileEnd(end *ast.EndStatement) ([]vm.Code, *errors.Error) {
	return c.compileBody(end.Body)
}

func (c *Compiler) compileLiteral(literal *ast.BasicLiteralExpression) ([]vm.Code, *errors.Error) {
	switch literal.DirectValue {
	case lexer.Integer:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		number, success := strconv.ParseInt(numberString, 10, 64)
		if success != nil {
			return nil, errors.New(literal.Token.Line, "Error parsing Integer", errors.GoRuntimeError)
		}
		return []vm.Code{vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number)}, nil
	case lexer.HexadecimalInteger:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 16, 64)
		if parsingError != nil {
			return nil, errors.New(literal.Token.Line, "Error parsing Hexadecimal Integer", errors.GoRuntimeError)
		}
		return []vm.Code{vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number)}, nil
	case lexer.OctalInteger:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 8, 64)
		if parsingError != nil {
			return nil, errors.New(literal.Token.Line, "Error parsing Octal Integer", errors.GoRuntimeError)
		}
		return []vm.Code{vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number)}, nil
	case lexer.BinaryInteger:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 2, 64)
		if parsingError != nil {
			return nil, errors.New(literal.Token.Line, "Error parsing Binary Integer", errors.GoRuntimeError)
		}
		return []vm.Code{vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number)}, nil
	case lexer.Float, lexer.ScientificFloat:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		number, parsingError := strconv.ParseFloat(numberString, 64)
		if parsingError != nil {
			return nil, errors.New(literal.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		return []vm.Code{vm.NewCode(vm.NewFloatOP, literal.Token.Line, number)}, nil
	case lexer.SingleQuoteString, lexer.DoubleQuoteString:
		return []vm.Code{vm.NewCode(
			vm.NewStringOP, literal.Token.Line,

			string(
				tools.ReplaceEscaped(
					[]rune(literal.Token.String[1:len(literal.Token.String)-1])),
			),
		),
		}, nil
	case lexer.ByteString:
		return []vm.Code{vm.NewCode(vm.NewBytesOP, literal.Token.Line,
			[]byte(
				string(
					tools.ReplaceEscaped(
						[]rune(literal.Token.String[2:len(literal.Token.String)-1]),
					),
				),
			),
		),
		}, nil
	case lexer.True:
		return []vm.Code{vm.NewCode(vm.NewTrueBoolOP, literal.Token.Line, nil)}, nil
	case lexer.False:
		return []vm.Code{vm.NewCode(vm.NewFalseBoolOP, literal.Token.Line, nil)}, nil
	case lexer.None:
		return []vm.Code{vm.NewCode(vm.GetNoneOP, literal.Token.Line, nil)}, nil
	}
	panic(errors.NewUnknownVMOperationError(literal.Token.DirectValue))
}

func (c *Compiler) compileTuple(tuple *ast.TupleExpression) ([]vm.Code, *errors.Error) {
	valuesLength := len(tuple.Values)
	var result []vm.Code
	for i := valuesLength - 1; i > -1; i-- {
		childExpression, valueCompilationError := c.compileExpression(true, tuple.Values[i])
		if valueCompilationError != nil {
			return nil, valueCompilationError
		}
		result = append(result, childExpression...)
	}
	return append(result, vm.NewCode(vm.NewTupleOP, errors.UnknownLine, len(tuple.Values))), nil
}

func (c *Compiler) compileArray(array *ast.ArrayExpression) ([]vm.Code, *errors.Error) {
	valuesLength := len(array.Values)
	var result []vm.Code
	for i := valuesLength - 1; i > -1; i-- {
		childExpression, valueCompilationError := c.compileExpression(true, array.Values[i])
		if valueCompilationError != nil {
			return nil, valueCompilationError
		}
		result = append(result, childExpression...)
	}
	return append(result, vm.NewCode(vm.NewArrayOP, errors.UnknownLine, len(array.Values))), nil
}

func (c *Compiler) compileHash(hash *ast.HashExpression) ([]vm.Code, *errors.Error) {
	valuesLength := len(hash.Values)
	var result []vm.Code
	for i := valuesLength - 1; i > -1; i-- {
		key, valueCompilationError := c.compileExpression(true, hash.Values[i].Value)
		if valueCompilationError != nil {
			return nil, valueCompilationError
		}
		result = append(result, key...)
		value, keyCompilationError := c.compileExpression(true, hash.Values[i].Key)
		if keyCompilationError != nil {
			return nil, keyCompilationError
		}
		result = append(result, value...)
	}
	return append(result, vm.NewCode(vm.NewHashOP, errors.UnknownLine, len(hash.Values))), nil
}

func (c *Compiler) compileUnaryExpression(unaryExpression *ast.UnaryExpression) ([]vm.Code, *errors.Error) {
	result, expressionCompileError := c.compileExpression(true, unaryExpression.X)
	if expressionCompileError != nil {
		return nil, expressionCompileError
	}
	switch unaryExpression.Operator.DirectValue {
	case lexer.NegateBits:
		result = append(result, vm.NewCode(vm.NegateBitsOP, unaryExpression.Operator.Line, nil))
	case lexer.Not, lexer.SignNot:
		result = append(result, vm.NewCode(vm.BoolNegateOP, unaryExpression.Operator.Line, nil))
	case lexer.Sub:
		result = append(result, vm.NewCode(vm.NegativeOP, unaryExpression.Operator.Line, nil))
	}
	return result, nil
}

func (c *Compiler) compileBinaryExpression(binaryExpression *ast.BinaryExpression) ([]vm.Code, *errors.Error) {
	var result []vm.Code
	// Compile first right hand side
	right, rightHandSideCompileError := c.compileExpression(true, binaryExpression.RightHandSide)
	if rightHandSideCompileError != nil {
		return nil, rightHandSideCompileError
	}
	result = append(result, right...)
	// Then left hand side
	left, leftHandSideCompileError := c.compileExpression(true, binaryExpression.LeftHandSide)
	if leftHandSideCompileError != nil {
		return nil, leftHandSideCompileError
	}
	result = append(result, left...)
	var operation uint8
	// Finally decide the instruction to use
	switch binaryExpression.Operator.DirectValue {
	case lexer.Add:
		operation = vm.AddOP
	case lexer.Sub:
		operation = vm.SubOP
	case lexer.Star:
		operation = vm.MulOP
	case lexer.Div:
		operation = vm.DivOP
	case lexer.FloorDiv:
		operation = vm.FloorDivOP
	case lexer.Modulus:
		operation = vm.ModOP
	case lexer.PowerOf:
		operation = vm.PowOP
	case lexer.BitwiseXor:
		operation = vm.BitXorOP
	case lexer.BitWiseAnd:
		operation = vm.BitAndOP
	case lexer.BitwiseOr:
		operation = vm.BitOrOP
	case lexer.BitwiseLeft:
		operation = vm.BitLeftOP
	case lexer.BitwiseRight:
		operation = vm.BitRightOP
	case lexer.And:
		operation = vm.AndOP
	case lexer.Or:
		operation = vm.OrOP
	case lexer.Xor:
		operation = vm.XorOP
	case lexer.Equals:
		operation = vm.EqualsOP
	case lexer.NotEqual:
		operation = vm.NotEqualsOP
	case lexer.GreaterThan:
		operation = vm.GreaterThanOP
	case lexer.LessThan:
		operation = vm.LessThanOP
	case lexer.GreaterOrEqualThan:
		operation = vm.GreaterThanOrEqualOP
	case lexer.LessOrEqualThan:
		operation = vm.LessThanOrEqualOP
	case lexer.In:
		operation = vm.ContainsOP
	default:
		panic(errors.NewUnknownVMOperationError(binaryExpression.Operator.DirectValue))
	}
	return append(result, vm.NewCode(operation, binaryExpression.Operator.Line, nil)), nil
}

func (c *Compiler) compileParenthesesExpression(parenthesesExpression *ast.ParenthesesExpression) ([]vm.Code, *errors.Error) {
	result, resultError := c.compileExpression(true, parenthesesExpression.X)
	if resultError != nil {
		return nil, resultError
	}
	result = append(result,
		vm.NewCode(vm.NewParenthesesOP, errors.UnknownLine, nil),
	)
	return result, nil
}

func (c *Compiler) compileIfOneLinerExpression(ifOneLineExpression *ast.IfOneLinerExpression) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, ifOneLineExpression.Condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}

	ifResult, resultCompilationError := c.compileExpression(true, ifOneLineExpression.Result)
	if resultCompilationError != nil {
		return nil, resultCompilationError
	}
	var elseResult []vm.Code
	if ifOneLineExpression.ElseResult != nil {
		var elseResultCompilationError *errors.Error
		elseResult, elseResultCompilationError = c.compileExpression(true, ifOneLineExpression.ElseResult)
		if elseResultCompilationError != nil {
			return nil, elseResultCompilationError
		}
	} else {
		elseResult = append(elseResult,
			vm.NewCode(vm.GetNoneOP, errors.UnknownLine, nil),
			vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
		)
	}
	elseResultLength := len(elseResult)
	result := condition
	result = append(result, vm.NewCode(vm.IfJumpOP, errors.UnknownLine, len(ifResult)+2))
	result = append(result, ifResult...)
	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, elseResultLength+1))
	result = append(result, elseResult...)
	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	return result, nil
}

func (c *Compiler) compileUnlessOneLinerExpression(ifOneLineExpression *ast.UnlessOneLinerExpression) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, ifOneLineExpression.Condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}

	ifResult, resultCompilationError := c.compileExpression(true, ifOneLineExpression.Result)
	if resultCompilationError != nil {
		return nil, resultCompilationError
	}

	var elseResult []vm.Code
	if ifOneLineExpression.ElseResult != nil {
		var elseResultCompilationError *errors.Error
		elseResult, elseResultCompilationError = c.compileExpression(true, ifOneLineExpression.ElseResult)
		if elseResultCompilationError != nil {
			return nil, elseResultCompilationError
		}
	} else {
		elseResult = append(elseResult,
			vm.NewCode(vm.GetNoneOP, errors.UnknownLine, nil),
			vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
		)
	}
	elseResultLength := len(elseResult)
	result := condition
	result = append(result, vm.NewCode(vm.UnlessJumpOP, errors.UnknownLine, len(ifResult)+1))
	result = append(result, ifResult...)
	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, elseResultLength+1))
	result = append(result, elseResult...)
	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))

	return result, nil
}

func (c *Compiler) compileIndexExpression(indexExpression *ast.IndexExpression) ([]vm.Code, *errors.Error) {
	source, sourceCompilationError := c.compileExpression(true, indexExpression.Source)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	index, indexCompilationError := c.compileExpression(true, indexExpression.Index)
	if indexCompilationError != nil {
		return nil, indexCompilationError
	}
	result := source
	result = append(result, index...)
	return append(result, vm.NewCode(vm.IndexOP, errors.UnknownLine, nil)), nil
}

func (c *Compiler) compileSelectorExpression(selectorExpression *ast.SelectorExpression) ([]vm.Code, *errors.Error) {
	source, sourceCompilationError := c.compileExpression(true, selectorExpression.X)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	return append(source, vm.NewCode(vm.SelectNameFromObjectOP, selectorExpression.Identifier.Token.Line, selectorExpression.Identifier.Token.String)), nil
}

func (c *Compiler) compileMethodInvocationExpression(methodInvocationExpression *ast.MethodInvocationExpression) ([]vm.Code, *errors.Error) {
	numberOfArguments := len(methodInvocationExpression.Arguments)
	var result []vm.Code
	for i := numberOfArguments - 1; i > -1; i-- {
		argument, argumentCompilationError := c.compileExpression(true, methodInvocationExpression.Arguments[i])
		if argumentCompilationError != nil {
			return nil, argumentCompilationError
		}
		result = append(result, argument...)
	}
	function, functionCompilationError := c.compileExpression(true, methodInvocationExpression.Function)
	if functionCompilationError != nil {
		return nil, functionCompilationError
	}
	result = append(result, function...)
	return append(result, vm.NewCode(vm.MethodInvocationOP, errors.UnknownLine, len(methodInvocationExpression.Arguments))), nil
}

func (c *Compiler) compileIdentifierExpression(identifier *ast.Identifier) ([]vm.Code, *errors.Error) {
	return []vm.Code{vm.NewCode(vm.GetIdentifierOP, identifier.Token.Line, identifier.Token.String)}, nil
}

func (c *Compiler) compileLambdaExpression(lambdaExpression *ast.LambdaExpression) ([]vm.Code, *errors.Error) {
	var result []vm.Code
	functionCode, lambdaCodeCompilationError := c.compileExpression(true, lambdaExpression.Code)
	if lambdaCodeCompilationError != nil {
		return nil, lambdaCodeCompilationError
	}
	result = append(result, vm.NewCode(vm.NewLambdaFunctionOP, errors.UnknownLine, [2]int{len(functionCode) + 2, len(lambdaExpression.Arguments)}))
	var arguments []string
	for _, argument := range lambdaExpression.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	result = append(result, vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	result = append(result, functionCode...)
	result = append(result, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1))
	return result, nil
}

func (c *Compiler) compileGeneratorExpression(generatorExpression *ast.GeneratorExpression) ([]vm.Code, *errors.Error) {
	// Compile the HasNext function
	hasNextCode, hasNextCallCompilationError := c.compileMethodInvocationExpression(
		&ast.MethodInvocationExpression{
			Function: &ast.SelectorExpression{
				X: &ast.SelectorExpression{
					X: &ast.Identifier{
						Token: &lexer.Token{
							String: vm.Self,
						},
					},
					Identifier: &ast.Identifier{
						Token: &lexer.Token{
							String: vm.Source,
						},
					},
				},
				Identifier: &ast.Identifier{
					Token: &lexer.Token{
						String: vm.HasNext,
					},
				},
			},
			Arguments: []ast.Expression{},
		},
	)
	if hasNextCallCompilationError != nil {
		return nil, hasNextCallCompilationError
	}
	hasNextCode = append(hasNextCode, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	hasNextCode = append(hasNextCode, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1))

	// Compile the Next function
	/// First capture the next value
	nextCode, temporalVariable1AssignError := c.compileAssignStatement(
		&ast.AssignStatement{
			LeftHandSide: &ast.Identifier{
				Token: &lexer.Token{
					String: vm.TemporalVariable1,
				},
			},
			AssignOperator: &lexer.Token{
				String:      "=",
				DirectValue: lexer.Assign,
				Kind:        lexer.Assignment,
			},
			RightHandSide: &ast.MethodInvocationExpression{
				Function: &ast.SelectorExpression{
					X: &ast.SelectorExpression{
						X: &ast.Identifier{
							Token: &lexer.Token{
								String: vm.Self,
							},
						},
						Identifier: &ast.Identifier{
							Token: &lexer.Token{
								String: vm.Source,
							},
						},
					},
					Identifier: &ast.Identifier{
						Token: &lexer.Token{
							String: vm.Next,
						},
					},
				},
				Arguments: []ast.Expression{},
			},
		},
	)
	if temporalVariable1AssignError != nil {
		return nil, temporalVariable1AssignError
	}
	//// Then Unpack the received value
	if len(generatorExpression.Receivers) == 1 {
		unpacked, receiverAssignError := c.compileAssignStatement(
			&ast.AssignStatement{
				LeftHandSide: &ast.Identifier{
					Token: &lexer.Token{
						String: generatorExpression.Receivers[0].Token.String,
					},
				},
				AssignOperator: &lexer.Token{
					String:      "=",
					DirectValue: lexer.Assign,
					Kind:        lexer.Assignment,
				},
				RightHandSide: &ast.Identifier{
					Token: &lexer.Token{
						String: vm.TemporalVariable1,
					},
				},
			},
		)
		if receiverAssignError != nil {
			return nil, receiverAssignError
		}
		nextCode = append(nextCode, unpacked...)
	} else {
		unpacked, temporalVariableAssignError := c.compileAssignStatement(
			&ast.AssignStatement{
				LeftHandSide: &ast.Identifier{
					Token: &lexer.Token{
						String: vm.TemporalVariable2,
					},
				},
				AssignOperator: &lexer.Token{
					String:      "=",
					DirectValue: lexer.Assign,
					Kind:        lexer.Assignment,
				},
				RightHandSide: &ast.MethodInvocationExpression{
					Function: &ast.SelectorExpression{
						X: &ast.Identifier{
							Token: &lexer.Token{
								String: vm.TemporalVariable1,
							},
						},
						Identifier: &ast.Identifier{
							Token: &lexer.Token{
								String: vm.Iter,
							},
						},
					},
					Arguments: []ast.Expression{},
				},
			},
		)
		if temporalVariableAssignError != nil {
			return nil, temporalVariableAssignError
		}
		nextCode = append(nextCode, unpacked...)
		for _, receiver := range generatorExpression.Receivers {
			compiledReceiver, receiverAssignCompilationError := c.compileAssignStatement(
				&ast.AssignStatement{
					LeftHandSide: &ast.Identifier{
						Token: &lexer.Token{
							String: receiver.Token.String,
						},
					},
					AssignOperator: &lexer.Token{
						String:      "=",
						DirectValue: lexer.Assign,
						Kind:        lexer.Assignment,
					},
					RightHandSide: &ast.MethodInvocationExpression{
						Function: &ast.SelectorExpression{
							X: &ast.Identifier{
								Token: &lexer.Token{
									String: vm.TemporalVariable2,
								},
							},
							Identifier: &ast.Identifier{
								Token: &lexer.Token{
									String: vm.Next,
								},
							},
						},
						Arguments: []ast.Expression{},
					},
				},
			)
			if receiverAssignCompilationError != nil {
				return nil, receiverAssignCompilationError
			}
			nextCode = append(nextCode, compiledReceiver...)
		}
	}
	//// Then Evaluate the operation and return its result
	evaluateOperand, operationCompilationError := c.compileExpression(true, generatorExpression.Operation)
	if operationCompilationError != nil {
		return nil, operationCompilationError
	}
	nextCode = append(nextCode, evaluateOperand...)
	nextCode = append(nextCode, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1))

	// Finally set everything together
	source, sourceCompilationError := c.compileExpression(true, generatorExpression.Source)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	result := source
	result = append(result,
		vm.NewCode(vm.NewIteratorOP, errors.UnknownLine,
			[2]int{len(hasNextCode), len(nextCode)},
		),
	)
	result = append(result, hasNextCode...)
	result = append(result, nextCode...)
	return result, nil
}

func (c *Compiler) compileExpression(pushExpression bool, expression ast.Expression) ([]vm.Code, *errors.Error) {
	var result []vm.Code
	var resultError *errors.Error
	switch expression.(type) {
	case *ast.BasicLiteralExpression:
		result, resultError = c.compileLiteral(expression.(*ast.BasicLiteralExpression))
	case *ast.TupleExpression:
		result, resultError = c.compileTuple(expression.(*ast.TupleExpression))
	case *ast.ArrayExpression:
		result, resultError = c.compileArray(expression.(*ast.ArrayExpression))
	case *ast.HashExpression:
		result, resultError = c.compileHash(expression.(*ast.HashExpression))
	case *ast.UnaryExpression:
		result, resultError = c.compileUnaryExpression(expression.(*ast.UnaryExpression))
	case *ast.BinaryExpression:
		result, resultError = c.compileBinaryExpression(expression.(*ast.BinaryExpression))
	case *ast.ParenthesesExpression:
		result, resultError = c.compileParenthesesExpression(expression.(*ast.ParenthesesExpression))
	case *ast.IfOneLinerExpression:
		result, resultError = c.compileIfOneLinerExpression(expression.(*ast.IfOneLinerExpression))
	case *ast.UnlessOneLinerExpression:
		result, resultError = c.compileUnlessOneLinerExpression(expression.(*ast.UnlessOneLinerExpression))
	case *ast.IndexExpression:
		result, resultError = c.compileIndexExpression(expression.(*ast.IndexExpression))
	case *ast.SelectorExpression:
		result, resultError = c.compileSelectorExpression(expression.(*ast.SelectorExpression))
	case *ast.MethodInvocationExpression:
		result, resultError = c.compileMethodInvocationExpression(expression.(*ast.MethodInvocationExpression))
	case *ast.Identifier:
		result, resultError = c.compileIdentifierExpression(expression.(*ast.Identifier))
	case *ast.LambdaExpression:
		result, resultError = c.compileLambdaExpression(expression.(*ast.LambdaExpression))
	case *ast.GeneratorExpression:
		result, resultError = c.compileGeneratorExpression(expression.(*ast.GeneratorExpression))
	default:
		panic(reflect.TypeOf(expression))
	}
	if pushExpression {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, resultError
}

// Statement compilation functions

func (c *Compiler) compileAssignStatementMiddleBinaryExpression(leftHandSide ast.Expression, assignOperator *lexer.Token) ([]vm.Code, *errors.Error) {
	result, leftHandSideCompilationError := c.compileExpression(true, leftHandSide)
	if leftHandSideCompilationError != nil {
		return nil, leftHandSideCompilationError
	}
	// Finally decide the instruction to use
	var operation uint8
	switch assignOperator.DirectValue {
	case lexer.AddAssign:
		operation = vm.AddOP
	case lexer.SubAssign:
		operation = vm.SubOP
	case lexer.StarAssign:
		operation = vm.MulOP
	case lexer.DivAssign:
		operation = vm.DivOP
	case lexer.ModulusAssign:
		operation = vm.ModOP
	case lexer.PowerOfAssign:
		operation = vm.PowOP
	case lexer.BitwiseXorAssign:
		operation = vm.BitXorOP
	case lexer.BitWiseAndAssign:
		operation = vm.BitAndOP
	case lexer.BitwiseOrAssign:
		operation = vm.BitOrOP
	case lexer.BitwiseLeftAssign:
		operation = vm.BitLeftOP
	case lexer.BitwiseRightAssign:
		operation = vm.BitRightOP
	default:
		panic(errors.NewUnknownVMOperationError(operation))
	}
	return append(result, vm.NewCode(operation, assignOperator.Line, nil)), nil
}

func (c *Compiler) compileIdentifierAssign(identifier *ast.Identifier) ([]vm.Code, *errors.Error) {
	return []vm.Code{vm.NewCode(vm.AssignIdentifierOP, identifier.Token.Line, identifier.Token.String)}, nil
}

func (c *Compiler) compileSelectorAssign(selectorExpression *ast.SelectorExpression) ([]vm.Code, *errors.Error) {
	result, sourceCompilationError := c.compileExpression(true, selectorExpression.X)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	return append(result, vm.NewCode(vm.AssignSelectorOP, selectorExpression.Identifier.Token.Line, selectorExpression.Identifier.Token.String)), nil
}

func (c *Compiler) compileIndexAssign(indexExpression *ast.IndexExpression) ([]vm.Code, *errors.Error) {
	result, sourceCompilationError := c.compileExpression(true, indexExpression.Source)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	index, indexCompilationError := c.compileExpression(true, indexExpression.Index)
	if indexCompilationError != nil {
		return nil, indexCompilationError
	}
	result = append(result, index...)
	return append(result, vm.NewCode(vm.AssignIndexOP, errors.UnknownLine, nil)), nil
}

func (c *Compiler) compileAssignStatement(assignStatement *ast.AssignStatement) ([]vm.Code, *errors.Error) {
	result, valueCompilationError := c.compileExpression(true, assignStatement.RightHandSide)
	if valueCompilationError != nil {
		return nil, valueCompilationError
	}
	if assignStatement.AssignOperator.DirectValue != lexer.Assign {
		// Do something here to evaluate the operation
		assignOperation, middleOperationCompilationError := c.compileAssignStatementMiddleBinaryExpression(assignStatement.LeftHandSide, assignStatement.AssignOperator)
		if middleOperationCompilationError != nil {
			return nil, middleOperationCompilationError
		}
		result = append(result, assignOperation...)
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	var leftHandSide []vm.Code
	var leftHandSideCompilationError *errors.Error
	switch assignStatement.LeftHandSide.(type) {
	case *ast.Identifier:
		leftHandSide, leftHandSideCompilationError = c.compileIdentifierAssign(assignStatement.LeftHandSide.(*ast.Identifier))
	case *ast.SelectorExpression:
		leftHandSide, leftHandSideCompilationError = c.compileSelectorAssign(assignStatement.LeftHandSide.(*ast.SelectorExpression))
	case *ast.IndexExpression:
		leftHandSide, leftHandSideCompilationError = c.compileIndexAssign(assignStatement.LeftHandSide.(*ast.IndexExpression))
	default:
		panic(reflect.TypeOf(assignStatement.LeftHandSide))
	}
	if leftHandSideCompilationError != nil {
		return nil, leftHandSideCompilationError
	}
	return append(result, leftHandSide...), nil
}

func (c *Compiler) compileFunctionDefinition(functionDefinition *ast.FunctionDefinitionStatement) ([]vm.Code, *errors.Error) {
	functionCode, functionDefinitionBodyCompilationError := c.compileBody(functionDefinition.Body)
	if functionDefinitionBodyCompilationError != nil {
		return nil, functionDefinitionBodyCompilationError
	}
	var result []vm.Code
	result = append(result, vm.NewCode(vm.NewFunctionOP, errors.UnknownLine, [2]int{len(functionCode) + 2, len(functionDefinition.Arguments)}))
	var arguments []string
	for _, argument := range functionDefinition.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	result = append(result, vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	result = append(result, functionCode...)
	result = append(result, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 0))
	return append(result, vm.NewCode(vm.AssignIdentifierOP, functionDefinition.Name.Token.Line, functionDefinition.Name.Token.String)), nil
}

func (c *Compiler) compileReturnStatement(returnStatement *ast.ReturnStatement) ([]vm.Code, *errors.Error) {
	numberOfResults := len(returnStatement.Results)
	var result []vm.Code
	for i := numberOfResults - 1; i > -1; i-- {
		returnResult, resultCompilationError := c.compileExpression(true, returnStatement.Results[i])
		if resultCompilationError != nil {
			return nil, resultCompilationError
		}
		result = append(result, returnResult...)
	}
	return append(result, vm.NewCode(vm.ReturnOP, errors.UnknownLine, numberOfResults)), nil
}

type ElifInformation struct {
	Condition       []vm.Code
	ConditionLength int
	Body            []vm.Code
	BodyLength      int
}

func (c *Compiler) compileIfStatement(ifStatement *ast.IfStatement) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, ifStatement.Condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	totalLength := 0
	body, bodyCompilationError := c.compileBody(ifStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	bodyLength := len(body)
	totalLength += bodyLength
	numberOfElifBlocks := 0
	var elifBlocks []ElifInformation
	for _, elifBlock := range ifStatement.ElifBlocks {
		elifCondition, elifConditionCompilationError := c.compileExpression(true, elifBlock.Condition)
		if elifConditionCompilationError != nil {
			return nil, elifConditionCompilationError
		}
		elifBody, elifBodyCompilationError := c.compileBody(elifBlock.Body)
		if elifBodyCompilationError != nil {
			return nil, elifBodyCompilationError
		}
		elifConditionLength := len(elifCondition)
		elifLength := len(elifBody)
		totalLength += elifConditionLength + 1 + elifLength + 1
		elifBlocks = append(elifBlocks, ElifInformation{
			Condition:       elifCondition,
			ConditionLength: elifConditionLength,
			Body:            elifBody,
			BodyLength:      elifLength,
		})
		numberOfElifBlocks++
	}
	elseBodyLength := 0
	var elseBody []vm.Code
	var elseCompilationError *errors.Error
	if ifStatement.Else != nil {
		elseBody, elseCompilationError = c.compileBody(ifStatement.Else)
		if elseCompilationError != nil {
			return nil, elseCompilationError
		}
		elseBodyLength = len(elseBody)
	}
	totalLength += elseBodyLength
	result := condition
	jump := 0
	if numberOfElifBlocks > 0 || elseBodyLength > 0 {
		jump = 1
	}
	result = append(result, vm.NewCode(vm.IfJumpOP, errors.UnknownLine, bodyLength+jump))
	result = append(result, body...)
	if numberOfElifBlocks > 0 {
		totalLength -= bodyLength
		result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, totalLength))
		for _, elifBlock := range elifBlocks {
			result = append(result, elifBlock.Condition...)
			result = append(result, vm.NewCode(vm.IfJumpOP, errors.UnknownLine, elifBlock.BodyLength+1))
			result = append(result, elifBlock.Body...)
			totalLength -= elifBlock.ConditionLength - 1 - elifBlock.BodyLength - 1
			result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, totalLength))
		}
		if elseBodyLength > 0 {
			result = append(result, elseBody...)
		}
	} else if elseBodyLength > 0 {
		result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, elseBodyLength))
		result = append(result, elseBody...)
	}
	return result, nil
}

func (c *Compiler) compileUnlessStatement(unlessStatement *ast.UnlessStatement) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, unlessStatement.Condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	totalLength := 0
	body, bodyCompilationError := c.compileBody(unlessStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	bodyLength := len(body)
	totalLength += bodyLength
	numberOfElifBlocks := 0
	var elifBlocks []ElifInformation
	for _, elifBlock := range unlessStatement.ElifBlocks {
		elifCondition, elifConditionCompilationError := c.compileExpression(true, elifBlock.Condition)
		if elifConditionCompilationError != nil {
			return nil, elifConditionCompilationError
		}
		elifBody, elifBodyCompilationError := c.compileBody(elifBlock.Body)
		if elifBodyCompilationError != nil {
			return nil, elifBodyCompilationError
		}
		elifConditionLength := len(elifCondition)
		elifLength := len(elifBody)
		totalLength += elifConditionLength + 1 + elifLength + 1
		elifBlocks = append(elifBlocks, ElifInformation{
			Condition:       elifCondition,
			ConditionLength: elifConditionLength,
			Body:            elifBody,
			BodyLength:      elifLength,
		})
		numberOfElifBlocks++
	}
	elseBodyLength := 0
	var elseBody []vm.Code
	var elseCompilationError *errors.Error
	if unlessStatement.Else != nil {
		elseBody, elseCompilationError = c.compileBody(unlessStatement.Else)
		if elseCompilationError != nil {
			return nil, elseCompilationError
		}
		elseBodyLength = len(elseBody)
	}
	totalLength += elseBodyLength
	result := condition
	jump := 0
	if numberOfElifBlocks > 0 || elseBodyLength > 0 {
		jump = 1
	}
	result = append(result, vm.NewCode(vm.UnlessJumpOP, errors.UnknownLine, bodyLength+jump))
	result = append(result, body...)
	if numberOfElifBlocks > 0 {
		totalLength -= bodyLength
		result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, totalLength))
		for _, elifBlock := range elifBlocks {
			result = append(result, elifBlock.Condition...)
			result = append(result, vm.NewCode(vm.UnlessJumpOP, errors.UnknownLine, elifBlock.BodyLength+1))
			result = append(result, elifBlock.Body...)
			totalLength -= elifBlock.ConditionLength - 1 - elifBlock.BodyLength - 1
			result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, totalLength))
		}
		if elseBodyLength > 0 {
			result = append(result, elseBody...)
		}
	} else if elseBodyLength > 0 {
		result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, elseBodyLength))
		result = append(result, elseBody...)
	}
	return result, nil
}

func (c *Compiler) compileRedoStatement() ([]vm.Code, *errors.Error) {
	return []vm.Code{vm.NewCode(vm.RedoOP, errors.UnknownLine, nil)}, nil
}

func (c *Compiler) compileBreakStatement() ([]vm.Code, *errors.Error) {
	return []vm.Code{vm.NewCode(vm.BreakOP, errors.UnknownLine, nil)}, nil
}

func (c *Compiler) compileContinueStatement() ([]vm.Code, *errors.Error) {
	return []vm.Code{vm.NewCode(vm.ContinueOP, errors.UnknownLine, nil)}, nil
}

func (c *Compiler) compilePassStatement() ([]vm.Code, *errors.Error) {
	return []vm.Code{vm.NewCode(vm.NOP, errors.UnknownLine, nil)}, nil
}

func (c *Compiler) compileDoWhileStatement(doWhileStatement *ast.DoWhileStatement) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, doWhileStatement.Condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	conditionLength := len(condition)
	body, bodyCompilationError := c.compileBody(doWhileStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	bodyLength := len(body)
	for index, instruction := range body {
		if instruction.Instruction.OpCode == vm.BreakOP && instruction.Value == nil {
			body[index].Value = (bodyLength - index) + conditionLength
		} else if instruction.Instruction.OpCode == vm.ContinueOP && instruction.Value == nil {
			body[index].Value = (bodyLength - index) - 1
		} else if instruction.Instruction.OpCode == vm.RedoOP && instruction.Value == nil {
			body[index].Value = -(index + 1)
		}
	}
	result := body
	result = append(result, condition...)
	result = append(result,
		vm.NewCode(vm.UnlessJumpOP, errors.UnknownLine, -(bodyLength+conditionLength+1)),
	)
	return result, nil
}

func (c *Compiler) compileWhileLoopStatement(whileStatement *ast.WhileLoopStatement) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, whileStatement.Condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	conditionLength := len(condition)
	body, bodyCompilationError := c.compileBody(whileStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	bodyLength := len(body)
	for index, instruction := range body {
		if instruction.Instruction.OpCode == vm.BreakOP && instruction.Value == nil {
			body[index].Value = bodyLength - index
		} else if instruction.Instruction.OpCode == vm.ContinueOP && instruction.Value == nil {
			body[index].Value = -(conditionLength + index + 2)
		} else if instruction.Instruction.OpCode == vm.RedoOP && instruction.Value == nil {
			body[index].Value = -(index + 1)
		}
	}
	result := condition
	result = append(result, vm.NewCode(vm.IfJumpOP, errors.UnknownLine, bodyLength+1))
	result = append(result, body...)
	result = append(result,
		vm.NewCode(vm.ContinueOP, errors.UnknownLine,
			-(conditionLength+1+bodyLength+1),
		),
	)
	return result, nil
}

func (c *Compiler) compileUntilLoopStatement(untilLoop *ast.UntilLoopStatement) ([]vm.Code, *errors.Error) {
	condition, conditionCompilationError := c.compileExpression(true, untilLoop.Condition)
	conditionLength := len(condition)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	body, bodyCompilationError := c.compileBody(untilLoop.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	bodyLength := len(body)
	for index, instruction := range body {
		if instruction.Instruction.OpCode == vm.BreakOP && instruction.Value == nil {
			body[index].Value = bodyLength - index
		} else if instruction.Instruction.OpCode == vm.ContinueOP && instruction.Value == nil {
			body[index].Value = -(conditionLength + index + 2)
		} else if instruction.Instruction.OpCode == vm.RedoOP && instruction.Value == nil {
			body[index].Value = -(index + 1)
		}
	}
	result := condition
	result = append(result, vm.NewCode(vm.UnlessJumpOP, errors.UnknownLine, bodyLength+1))
	result = append(result, body...)
	result = append(result,
		vm.NewCode(vm.ContinueOP, errors.UnknownLine,
			-(conditionLength+1+bodyLength+1),
		),
	)
	return result, nil
}

func (c *Compiler) compileForLoopStatement(forStatement *ast.ForLoopStatement) ([]vm.Code, *errors.Error) {
	source, sourceCompilationError := c.compileExpression(true, forStatement.Source)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	var receivers []string
	for _, receiver := range forStatement.Receivers {
		receivers = append(receivers, receiver.Token.String)
	}
	body, compilationError := c.compileBody(forStatement.Body)
	if compilationError != nil {
		return nil, compilationError
	}
	bodyLength := len(body)
	for index, instruction := range body {
		if instruction.Instruction.OpCode == vm.BreakOP && instruction.Value == nil {
			body[index].Value = bodyLength - index
		} else if instruction.Instruction.OpCode == vm.ContinueOP && instruction.Value == nil {
			body[index].Value = -(index + 3)
		} else if instruction.Instruction.OpCode == vm.RedoOP && instruction.Value == nil {
			body[index].Value = -(index + 2)
		}
	}
	result := source
	result = append(result,
		vm.NewCode(vm.SetupLoopOP, errors.UnknownLine,
			[2]interface{}{
				receivers, bodyLength + 2,
			},
		),
	)
	result = append(result, vm.NewCode(vm.UnpackForLoopOP, errors.UnknownLine, nil))
	result = append(result, vm.NewCode(vm.LoadForReloadOP, errors.UnknownLine, nil))
	result = append(result, body...)
	result = append(result, vm.NewCode(vm.ContinueOP, errors.UnknownLine, -(bodyLength+3)))
	return result, nil
}

type exceptBlock struct {
	Targets       []vm.Code
	TargetsLength int
	Receiver      string
	Body          []vm.Code
	BodyLength    int
}

func (c *Compiler) compileTryStatement(tryStatement *ast.TryStatement) ([]vm.Code, *errors.Error) {
	body, bodyCompilationError := c.compileBody(tryStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	bodyLength := len(body)
	totalLength := bodyLength
	var exceptBlocks []*exceptBlock
	var numberOfExceptBlocks int
	for _, except := range tryStatement.ExceptBlocks {
		targets, targetCompilationError := c.compileExpression(true,
			&ast.TupleExpression{
				Values: except.Targets,
			},
		)
		if targetCompilationError != nil {
			return nil, targetCompilationError
		}
		targetsLength := len(targets)

		totalLength += targetsLength + 1

		exceptBlockBody, exceptBlockBodyCompilationError := c.compileBody(except.Body)
		if exceptBlockBodyCompilationError != nil {
			return nil, exceptBlockBodyCompilationError
		}
		exceptBlockBodyLength := len(exceptBlockBody)
		totalLength += exceptBlockBodyLength + 1
		receiver := vm.JunkVariable
		if except.CaptureName != nil {
			receiver = except.CaptureName.Token.String
		}
		exceptBlocks = append(exceptBlocks,
			&exceptBlock{
				Targets:       targets,
				TargetsLength: targetsLength,
				Receiver:      receiver,
				Body:          exceptBlockBody,
				BodyLength:    exceptBlockBodyLength,
			},
		)
		numberOfExceptBlocks++
	}
	elseBody, elseCompilationError := c.compileBody(tryStatement.Else)
	if elseCompilationError != nil {
		return nil, elseCompilationError
	}
	elseLength := len(elseBody)
	totalLength += elseLength

	finallyBody, finallyCompilationError := c.compileBody(tryStatement.Finally)
	if finallyCompilationError != nil {
		return nil, finallyCompilationError
	}
	finallyLength := len(finallyBody)

	var result []vm.Code
	result = append(result, vm.NewCode(vm.SetupTryOP, errors.UnknownLine, bodyLength))
	result = append(result, body...)
	if numberOfExceptBlocks > 0 {
		totalLength -= bodyLength
		result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, totalLength))
		for _, except := range exceptBlocks {
			result = append(result, except.Targets...)
			result = append(result,
				vm.NewCode(vm.ExceptOP, errors.UnknownLine,
					[2]interface{}{except.Receiver, except.BodyLength + 1},
				),
			)
			result = append(result, except.Body...)
			totalLength -= except.TargetsLength + 1 + except.BodyLength + 1
			result = append(result,
				vm.NewCode(vm.JumpOP, errors.UnknownLine, totalLength),
			)
		}
		if elseLength > 0 {
			// result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, elseLength))
			result = append(result, elseBody...)
			if finallyLength > 0 {
				result = append(result, finallyBody...)
			}
		} else if finallyLength > 0 {
			result = append(result, finallyBody...)
		}
	} else if elseLength > 0 {
		result = append(result, vm.NewCode(vm.JumpOP, errors.UnknownLine, elseLength))
		result = append(result, elseBody...)
		if finallyLength > 0 {
			result = append(result, finallyBody...)
		}
	} else if finallyLength > 0 {
		result = append(result, finallyBody...)
	}
	result = append(result, vm.NewCode(vm.PopTryOP, errors.UnknownLine, nil))
	return result, nil
}

func (c *Compiler) compileModuleStatement(moduleStatement *ast.ModuleStatement) ([]vm.Code, *errors.Error) {
	moduleBody, moduleBodyCompilationError := c.compileBody(moduleStatement.Body)
	if moduleBodyCompilationError != nil {
		return nil, moduleBodyCompilationError
	}
	var result []vm.Code
	result = append(result,
		vm.NewCode(vm.NewModuleOP, moduleStatement.Name.Token.Line,
			vm.ModuleInformation{
				Name:       moduleStatement.Name.Token.String,
				CodeLength: len(moduleBody),
			},
		),
	)
	result = append(result, moduleBody...)
	return result, nil
}

func (c *Compiler) compileRaiseStatement(raise *ast.RaiseStatement) ([]vm.Code, *errors.Error) {
	result, expressionCompilationError := c.compileExpression(true, raise.X)
	if expressionCompilationError != nil {
		return nil, expressionCompilationError
	}
	result = append(result, vm.NewCode(vm.RaiseOP, errors.UnknownLine, nil))
	return result, nil
}

func (c *Compiler) compileClassFunctionDefinition(functionDefinition *ast.FunctionDefinitionStatement) ([]vm.Code, *errors.Error, bool) {
	functionCode, functionDefinitionBodyCompilationError := c.compileBody(functionDefinition.Body)
	if functionDefinitionBodyCompilationError != nil {
		return nil, functionDefinitionBodyCompilationError, false
	}
	var result []vm.Code
	result = append(result, vm.NewCode(vm.NewClassFunctionOP, errors.UnknownLine, [2]int{len(functionCode) + 2, len(functionDefinition.Arguments)}))
	var arguments []string
	for _, argument := range functionDefinition.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	result = append(result, vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	result = append(result, functionCode...)
	result = append(result, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 0))
	result = append(result, vm.NewCode(vm.AssignIdentifierOP, functionDefinition.Name.Token.Line, functionDefinition.Name.Token.String))
	return result, nil, functionDefinition.Name.Token.String == vm.Initialize
}

func (c *Compiler) compileClassBody(body []ast.Node) ([]vm.Code, *errors.Error) {
	foundInitialize := false
	var isInitialize bool
	var nodeCode []vm.Code
	var compilationError *errors.Error
	var result []vm.Code
	for _, node := range body {
		switch node.(type) {
		case ast.Expression:
			nodeCode, compilationError = c.compileExpression(true, node.(ast.Expression))
		case ast.Statement:
			if _, ok := node.(*ast.FunctionDefinitionStatement); ok {
				nodeCode, compilationError, isInitialize = c.compileClassFunctionDefinition(node.(*ast.FunctionDefinitionStatement))
				if isInitialize && !foundInitialize {
					foundInitialize = true
				}
			} else {
				nodeCode, compilationError = c.compileStatement(node.(ast.Statement))
			}
		}
		if compilationError != nil {
			return nil, compilationError
		}
		result = append(result, nodeCode...)
	}
	if !foundInitialize {
		nodeCode, _, _ = c.compileClassFunctionDefinition(
			&ast.FunctionDefinitionStatement{
				Name: &ast.Identifier{
					Token: &lexer.Token{
						String: vm.Initialize,
					},
				},
				Arguments: nil,
				Body:      nil,
			},
		)
		result = append(result, nodeCode...)
	}
	return result, nil
}

func (c *Compiler) compileInterfaceStatement(interfaceStatement *ast.InterfaceStatement) ([]vm.Code, *errors.Error) {
	bases := make([]ast.Expression, len(interfaceStatement.Bases))
	copy(bases, interfaceStatement.Bases)
	for i, j := 0, len(bases)-1; i < j; i, j = i+1, j-1 {
		bases[i], bases[j] = bases[j], bases[i]
	}
	result, basesCompilationError := c.compileExpression(true,
		&ast.TupleExpression{
			Values: bases,
		},
	)
	if basesCompilationError != nil {
		return nil, basesCompilationError
	}
	var interfaceMethods []ast.Node
	for _, functionDefinition := range interfaceStatement.MethodDefinitions {
		interfaceMethods = append(interfaceMethods, functionDefinition)
	}
	body, bodyCompilationError := c.compileClassBody(
		interfaceMethods,
	)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	result = append(result,
		vm.NewCode(vm.NewClassOP, interfaceStatement.Name.Token.Line,
			vm.ClassInformation{
				Name:       interfaceStatement.Name.Token.String,
				BodyLength: len(body),
			},
		),
	)
	result = append(result, body...)
	return result, nil
}

func (c *Compiler) compileClassStatement(classStatement *ast.ClassStatement) ([]vm.Code, *errors.Error) {
	bases := make([]ast.Expression, len(classStatement.Bases))
	copy(bases, classStatement.Bases)
	for i, j := 0, len(bases)-1; i < j; i, j = i+1, j-1 {
		bases[i], bases[j] = bases[j], bases[i]
	}
	result, basesCompilationError := c.compileExpression(true,
		&ast.TupleExpression{
			Values: bases,
		},
	)
	if basesCompilationError != nil {
		return nil, basesCompilationError
	}
	body, bodyCompilationError := c.compileClassBody(classStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	result = append(result,
		vm.NewCode(vm.NewClassOP, classStatement.Name.Token.Line,
			vm.ClassInformation{
				Name:       classStatement.Name.Token.String,
				BodyLength: len(body),
			},
		),
	)
	result = append(result, body...)
	return result, nil
}

func (c *Compiler) compileSwitchStatement(switchStatement *ast.SwitchStatement) ([]vm.Code, *errors.Error) {
	target, targetCompilationError := c.compileExpression(true, switchStatement.Target)
	if targetCompilationError != nil {
		return nil, targetCompilationError
	}

	var cases []struct {
		references []vm.Code
		body       []vm.Code
	}
	for _, case_ := range switchStatement.CaseBlocks {
		references, referencesCompilationError := c.compileExpression(true,
			&ast.TupleExpression{
				Values: case_.Cases,
			},
		)
		if referencesCompilationError != nil {
			return nil, referencesCompilationError
		}
		body, bodyCompilationError := c.compileBody(case_.Body)
		if bodyCompilationError != nil {
			return nil, bodyCompilationError
		}
		cases = append(cases,
			struct {
				references []vm.Code
				body       []vm.Code
			}{
				references: references,
				body:       body,
			},
		)
	}
	defaultBody, defaultCompilationError := c.compileBody(switchStatement.Default)
	if defaultCompilationError != nil {
		return nil, defaultCompilationError
	}
	var switchBody []vm.Code
	// Construct the switch
	totalLength := len(defaultBody) + 1
	for _, caseBlock := range cases {
		bodyLength := len(caseBlock.body)
		referencesLength := len(caseBlock.references)
		totalLength += referencesLength + 1 + bodyLength + 1
		switchBody = append(switchBody, caseBlock.references...)
		switchBody = append(switchBody, vm.NewCode(vm.CaseOP, errors.UnknownLine, bodyLength+1))
		switchBody = append(switchBody, caseBlock.body...)
		switchBody = append(switchBody, vm.NewCode(vm.JumpOP, errors.UnknownLine, nil))
	}
	switchBody = append(switchBody, vm.NewCode(vm.PopOP, errors.UnknownLine, nil))
	switchBody = append(switchBody, defaultBody...)

	for index, code := range switchBody {
		if code.Instruction.OpCode == vm.JumpOP && code.Value == nil {
			switchBody[index].Value = totalLength - index - 1
		}
	}
	result := target
	result = append(result, switchBody...)
	return result, nil
}

func (c *Compiler) compileStatement(statement ast.Statement) ([]vm.Code, *errors.Error) {
	switch statement.(type) {
	case *ast.AssignStatement:
		return c.compileAssignStatement(statement.(*ast.AssignStatement))
	case *ast.FunctionDefinitionStatement:
		return c.compileFunctionDefinition(statement.(*ast.FunctionDefinitionStatement))
	case *ast.ReturnStatement:
		return c.compileReturnStatement(statement.(*ast.ReturnStatement))
	case *ast.IfStatement:
		return c.compileIfStatement(statement.(*ast.IfStatement))
	case *ast.UnlessStatement:
		return c.compileUnlessStatement(statement.(*ast.UnlessStatement))
	case *ast.DoWhileStatement:
		return c.compileDoWhileStatement(statement.(*ast.DoWhileStatement))
	case *ast.RedoStatement:
		return c.compileRedoStatement()
	case *ast.BreakStatement:
		return c.compileBreakStatement()
	case *ast.ContinueStatement:
		return c.compileContinueStatement()
	case *ast.PassStatement:
		return c.compilePassStatement()
	case *ast.WhileLoopStatement:
		return c.compileWhileLoopStatement(statement.(*ast.WhileLoopStatement))
	case *ast.UntilLoopStatement:
		return c.compileUntilLoopStatement(statement.(*ast.UntilLoopStatement))
	case *ast.ForLoopStatement:
		return c.compileForLoopStatement(statement.(*ast.ForLoopStatement))
	case *ast.TryStatement:
		return c.compileTryStatement(statement.(*ast.TryStatement))
	case *ast.ModuleStatement:
		return c.compileModuleStatement(statement.(*ast.ModuleStatement))
	case *ast.RaiseStatement:
		return c.compileRaiseStatement(statement.(*ast.RaiseStatement))
	case *ast.ClassStatement:
		return c.compileClassStatement(statement.(*ast.ClassStatement))
	case *ast.InterfaceStatement:
		return c.compileInterfaceStatement(statement.(*ast.InterfaceStatement))
	case *ast.SwitchStatement:
		return c.compileSwitchStatement(statement.(*ast.SwitchStatement))
	}
	panic(reflect.TypeOf(statement))
}

func (c *Compiler) compile(node ast.Node) ([]vm.Code, *errors.Error) {
	switch node.(type) {
	case ast.Expression:
		return c.compileExpression(false, node.(ast.Expression))
	case ast.Statement:
		return c.compileStatement(node.(ast.Statement))
	}
	panic(reflect.TypeOf(node))
}

func (c *Compiler) compileBody(body []ast.Node) ([]vm.Code, *errors.Error) {
	var result []vm.Code
	for _, node := range body {
		nodeCode, compileError := c.compile(node)
		if compileError != nil {
			return nil, compileError
		}
		result = append(result, nodeCode...)
	}
	return result, nil
}

func (c *Compiler) CompileToArray() ([]vm.Code, *errors.Error) {
	codeAst, parsingError := c.parser.Parse()
	if parsingError != nil {
		return nil, parsingError
	}
	var result []vm.Code
	if codeAst.Begin != nil {
		begin, compileError := c.compileBegin(codeAst.Begin)
		if compileError != nil {
			return nil, compileError
		}
		result = append(result, begin...)
	}
	body, compileError := c.compileBody(codeAst.Body)
	if compileError != nil {
		return nil, compileError
	}
	result = append(result, body...)

	if codeAst.End != nil {
		var end []vm.Code
		end, compileError = c.compileEnd(codeAst.End)
		if compileError != nil {
			return nil, compileError
		}
		result = append(result, end...)
	}
	if c.options.Debug {
		fmt.Println("---- Compiled Code ----")
		for i, ins := range result {
			fmt.Println(i, ins.Instruction, ins.Value)
		}
		fmt.Println("---- Compiled Code End ----")
	}
	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	result = append(result, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1))
	return result, nil
}

func (c *Compiler) Compile() (*vm.Bytecode, *errors.Error) {
	result, compilationError := c.CompileToArray()
	if compilationError != nil {
		return nil, compilationError
	}
	return vm.NewBytecodeFromArray(result), nil
}

func NewCompiler(
	codeReader reader.Reader,
	options Options,
) *Compiler {
	return &Compiler{
		parser:  parser.NewParser(lexer.NewLexer(codeReader)),
		options: options,
	}
}
