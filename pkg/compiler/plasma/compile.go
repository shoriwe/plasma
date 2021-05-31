package plasma

import (
	"github.com/shoriwe/gruby/pkg/cleanup"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
	"github.com/shoriwe/gruby/pkg/compiler/parser"
	"github.com/shoriwe/gruby/pkg/errors"
	"github.com/shoriwe/gruby/pkg/reader"
	"github.com/shoriwe/gruby/pkg/vm"
	"strconv"
	"strings"
)

/*
	Compile to the Plasma stack VM
*/

type Compiler struct {
	parser       *parser.Parser
	instructions []vm.Code
	index        int
	length       int
}

func (c *Compiler) pushInstruction(code vm.Code) {
	c.instructions = append(c.instructions, code)
	c.index++
	c.length++
}

func (c *Compiler) extendInstructions(code []vm.Code) {
	c.instructions = append(c.instructions, code...)
}

func (c *Compiler) compileBegin(begin *ast.BeginStatement) *errors.Error {
	if begin != nil {
		return c.compileBody(begin.Body)
	}
	return nil
}

func (c *Compiler) compileEnd(end *ast.EndStatement) *errors.Error {
	if end != nil {
		return c.compileBody(end.Body)
	}
	return nil
}

func (c *Compiler) compileLiteral(literal *ast.BasicLiteralExpression) *errors.Error {
	switch literal.DirectValue {
	case lexer.Integer:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		number, parsingError := strconv.ParseInt(numberString, 10, 64)
		if parsingError != nil {
			return errors.New(literal.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		c.pushInstruction(vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number))
	case lexer.HexadecimalInteger:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 16, 64)
		if parsingError != nil {
			return errors.New(literal.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		c.pushInstruction(vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number))
	case lexer.OctalInteger:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 8, 64)
		if parsingError != nil {
			return errors.New(literal.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		c.pushInstruction(vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number))
	case lexer.BinaryInteger:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 2, 64)
		if parsingError != nil {
			return errors.New(literal.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		c.pushInstruction(vm.NewCode(vm.NewIntegerOP, literal.Token.Line, number))
	case lexer.Float, lexer.ScientificFloat:
		numberString := literal.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		number, parsingError := cleanup.ParseFloat(numberString)
		if parsingError != nil {
			return errors.New(literal.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		c.pushInstruction(vm.NewCode(vm.NewFloatOP, literal.Token.Line, number))
	case lexer.SingleQuoteString, lexer.DoubleQuoteString:
		c.pushInstruction(
			vm.NewCode(
				vm.NewStringOP, literal.Token.Line,

				string(
					cleanup.ReplaceEscaped(
						[]rune(literal.Token.String[1:len(literal.Token.String)-1])),
				),
			),
		)
	case lexer.ByteString:
		c.pushInstruction(
			vm.NewCode(vm.NewBytesOP, literal.Token.Line,
				[]byte(
					string(
						cleanup.ReplaceEscaped(
							[]rune(literal.Token.String[2:len(literal.Token.String)-1]),
						),
					),
				),
			),
		)
	case lexer.True:
		c.pushInstruction(vm.NewCode(vm.NewTrueBoolOP, literal.Token.Line, nil))
	case lexer.False:
		c.pushInstruction(vm.NewCode(vm.NewFalseBoolOP, literal.Token.Line, nil))
	case lexer.None:
		c.pushInstruction(vm.NewCode(vm.GetNoneOP, literal.Token.Line, nil))
	}
	return nil
}

func (c *Compiler) compileTuple(tuple *ast.TupleExpression) *errors.Error {
	valuesLength := len(tuple.Values)
	for i := valuesLength - 1; i > -1; i-- {
		valueCompilationError := c.compileExpression(tuple.Values[i])
		if valueCompilationError != nil {
			return valueCompilationError
		}
	}
	c.instructions = append(c.instructions, vm.NewCode(vm.NewTupleOP, errors.UnknownLine, len(tuple.Values)))
	return nil
}

func (c *Compiler) compileArray(array *ast.ArrayExpression) *errors.Error {
	valuesLength := len(array.Values)
	for i := valuesLength - 1; i > -1; i-- {
		valueCompilationError := c.compileExpression(array.Values[i])
		if valueCompilationError != nil {
			return valueCompilationError
		}
	}
	c.instructions = append(c.instructions, vm.NewCode(vm.NewArrayOP, errors.UnknownLine, len(array.Values)))
	return nil
}

func (c *Compiler) compileHash(hash *ast.HashExpression) *errors.Error {
	valuesLength := len(hash.Values)
	for i := valuesLength - 1; i > -1; i-- {
		valueCompilationError := c.compileExpression(hash.Values[i].Value)
		if valueCompilationError != nil {
			return valueCompilationError
		}
		keyCompilationError := c.compileExpression(hash.Values[i].Key)
		if keyCompilationError != nil {
			return keyCompilationError
		}
	}
	c.instructions = append(c.instructions, vm.NewCode(vm.NewHashOP, errors.UnknownLine, len(hash.Values)))
	return nil
}

func (c *Compiler) compileUnaryExpression(unaryExpression *ast.UnaryExpression) *errors.Error {
	expressionCompileError := c.compileExpression(unaryExpression.X)
	if expressionCompileError != nil {
		return expressionCompileError
	}
	switch unaryExpression.Operator.DirectValue {
	case lexer.NegateBits:
		c.pushInstruction(vm.NewCode(vm.NegateBitsOP, unaryExpression.Operator.Line, nil))
	case lexer.Not, lexer.SignNot:
		c.pushInstruction(vm.NewCode(vm.BoolNegateOP, unaryExpression.Operator.Line, nil))
	case lexer.Sub:
		c.pushInstruction(vm.NewCode(vm.NegativeOP, unaryExpression.Operator.Line, nil))
	}
	return nil
}

func (c *Compiler) compileBinaryExpression(binaryExpression *ast.BinaryExpression) *errors.Error {
	// Compile first right hand side
	rightHandSideCompileError := c.compileExpression(binaryExpression.RightHandSide)
	if rightHandSideCompileError != nil {
		return rightHandSideCompileError
	}
	// Then left hand side
	leftHandSideCompileError := c.compileExpression(binaryExpression.LeftHandSide)
	if leftHandSideCompileError != nil {
		return leftHandSideCompileError
	}
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
	default:
		panic(errors.NewUnknownVMOperationError(operation))
	}
	c.pushInstruction(vm.NewCode(operation, binaryExpression.Operator.Line, nil))
	return nil
}

func (c *Compiler) compileParenthesesExpression(parenthesesExpression *ast.ParenthesesExpression) *errors.Error {
	return c.compileExpression(parenthesesExpression.X)
}

func (c *Compiler) compileIfOneLinerExpression(ifOneLineExpression *ast.IfOneLinerExpression) *errors.Error {
	return nil
}

func (c *Compiler) compileUnlessOneLinerExpression(ifOneLineExpression *ast.UnlessOneLinerExpression) *errors.Error {
	return nil
}

func (c *Compiler) compileIndexExpression(indexExpression *ast.IndexExpression) *errors.Error {
	sourceCompilationError := c.compileExpression(indexExpression.Source)
	if sourceCompilationError != nil {
		return sourceCompilationError
	}
	indexCompilationError := c.compileExpression(indexExpression.Index)
	if indexCompilationError != nil {
		return indexCompilationError
	}
	c.pushInstruction(vm.NewCode(vm.IndexOP, errors.UnknownLine, nil))
	return nil
}

func (c *Compiler) compileSelectorExpression(selectorExpression *ast.SelectorExpression) *errors.Error {
	sourceCompilationError := c.compileExpression(selectorExpression.X)
	if sourceCompilationError != nil {
		return sourceCompilationError
	}
	c.pushInstruction(vm.NewCode(vm.SelectNameFromObjectOP, selectorExpression.Identifier.Token.Line, selectorExpression.Identifier.Token.String))
	return nil
}

func (c *Compiler) compileMethodInvocationExpression(methodInvocationExpression *ast.MethodInvocationExpression) *errors.Error {
	numberOfArguments := len(methodInvocationExpression.Arguments)
	for i := numberOfArguments - 1; i > -1; i-- {
		argumentCompilationError := c.compileExpression(methodInvocationExpression.Arguments[i])
		if argumentCompilationError != nil {
			return argumentCompilationError
		}
	}
	functionCompilationError := c.compileExpression(methodInvocationExpression.Function)
	if functionCompilationError != nil {
		return functionCompilationError
	}
	c.pushInstruction(vm.NewCode(vm.MethodInvocationOP, errors.UnknownLine, len(methodInvocationExpression.Arguments)))
	return nil
}

func (c *Compiler) compileIdentifierExpression(identifier *ast.Identifier) *errors.Error {
	c.pushInstruction(vm.NewCode(vm.GetIdentifierOP, identifier.Token.Line, identifier.Token.String))
	return nil
}

func (c *Compiler) compileLambdaExpression(lambdaExpression *ast.LambdaExpression) *errors.Error {
	instructionsBackup := c.instructions
	c.instructions = nil
	lambdaCodeCompilationError := c.compileExpression(lambdaExpression.Code)
	if lambdaCodeCompilationError != nil {
		return lambdaCodeCompilationError
	}
	functionCode := c.instructions
	c.instructions = nil
	c.instructions = instructionsBackup
	c.pushInstruction(vm.NewCode(vm.NewFunctionOP, errors.UnknownLine, [2]int{len(functionCode) + 2, len(lambdaExpression.Arguments)}))
	var arguments []string
	for _, argument := range lambdaExpression.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	c.pushInstruction(vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	c.extendInstructions(functionCode)
	c.pushInstruction(vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1))
	return nil
}

func (c *Compiler) compileExpression(expression ast.Expression) *errors.Error {
	switch expression.(type) {
	case *ast.BasicLiteralExpression:
		return c.compileLiteral(expression.(*ast.BasicLiteralExpression))
	case *ast.TupleExpression:
		return c.compileTuple(expression.(*ast.TupleExpression))
	case *ast.ArrayExpression:
		return c.compileArray(expression.(*ast.ArrayExpression))
	case *ast.HashExpression:
		return c.compileHash(expression.(*ast.HashExpression))
	case *ast.UnaryExpression:
		return c.compileUnaryExpression(expression.(*ast.UnaryExpression))
	case *ast.BinaryExpression:
		return c.compileBinaryExpression(expression.(*ast.BinaryExpression))
	case *ast.ParenthesesExpression:
		return c.compileParenthesesExpression(expression.(*ast.ParenthesesExpression))
	case *ast.IfOneLinerExpression:
		return c.compileIfOneLinerExpression(expression.(*ast.IfOneLinerExpression))
	case *ast.UnlessOneLinerExpression:
		return c.compileUnlessOneLinerExpression(expression.(*ast.UnlessOneLinerExpression))
	case *ast.IndexExpression:
		return c.compileIndexExpression(expression.(*ast.IndexExpression))
	case *ast.SelectorExpression:
		return c.compileSelectorExpression(expression.(*ast.SelectorExpression))
	case *ast.MethodInvocationExpression:
		return c.compileMethodInvocationExpression(expression.(*ast.MethodInvocationExpression))
	case *ast.Identifier:
		return c.compileIdentifierExpression(expression.(*ast.Identifier))
	case *ast.LambdaExpression:
		return c.compileLambdaExpression(expression.(*ast.LambdaExpression))
	}
	return nil
}

// Statement compilation functions

func (c *Compiler) compileAssignStatementMiddleBinaryExpression(leftHandSide ast.Expression, assignOperator *lexer.Token) *errors.Error {
	leftHandSideCompilationError := c.compileExpression(leftHandSide)
	if leftHandSideCompilationError != nil {
		return leftHandSideCompilationError
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
	c.pushInstruction(vm.NewCode(operation, assignOperator.Line, nil))
	return nil
}

func (c *Compiler) compileIdentifierAssign(identifier *ast.Identifier) *errors.Error {
	c.pushInstruction(vm.NewCode(vm.AssignIdentifierOP, identifier.Token.Line, identifier.Token.String))
	return nil
}

func (c *Compiler) compileSelectorAssign(selectorExpression *ast.SelectorExpression) *errors.Error {
	sourceCompilationError := c.compileExpression(selectorExpression.X)
	if sourceCompilationError != nil {
		return sourceCompilationError
	}
	c.pushInstruction(vm.NewCode(vm.AssignSelectorOP, selectorExpression.Identifier.Token.Line, selectorExpression.Identifier.Token.String))
	return nil
}

func (c *Compiler) compileIndexAssign(indexExpression *ast.IndexExpression) *errors.Error {
	sourceCompilationError := c.compileExpression(indexExpression.Source)
	if sourceCompilationError != nil {
		return sourceCompilationError
	}
	indexCompilationError := c.compileExpression(indexExpression.Index)
	if indexCompilationError != nil {
		return indexCompilationError
	}
	c.pushInstruction(vm.NewCode(vm.AssignIndexOP, errors.UnknownLine, nil))
	return nil
}

func (c *Compiler) compileAssignStatement(assignStatement *ast.AssignStatement) *errors.Error {
	valueCompilationError := c.compileExpression(assignStatement.RightHandSide)
	if valueCompilationError != nil {
		return valueCompilationError
	}
	if assignStatement.AssignOperator.DirectValue != lexer.Assign {
		// Do something here to evaluate the operation
		middleOperationCompilationError := c.compileAssignStatementMiddleBinaryExpression(assignStatement.LeftHandSide, assignStatement.AssignOperator)
		if middleOperationCompilationError != nil {
			return middleOperationCompilationError
		}
	}
	switch assignStatement.LeftHandSide.(type) {
	case *ast.Identifier:
		return c.compileIdentifierAssign(assignStatement.LeftHandSide.(*ast.Identifier))
	case *ast.SelectorExpression:
		return c.compileSelectorAssign(assignStatement.LeftHandSide.(*ast.SelectorExpression))
	case *ast.IndexExpression:
		return c.compileIndexAssign(assignStatement.LeftHandSide.(*ast.IndexExpression))
	}
	// ToDo: Fix this return a better error
	return errors.NewUnknownVMOperationError(errors.UnknownLine)
}

func (c *Compiler) compileFunctionDefinition(functionDefinition *ast.FunctionDefinitionStatement) *errors.Error {
	instructionsBackup := c.instructions
	c.instructions = nil
	functionDefinitionBodyCompilationError := c.compileBody(functionDefinition.Body)
	if functionDefinitionBodyCompilationError != nil {
		return functionDefinitionBodyCompilationError
	}
	functionCode := c.instructions
	c.instructions = nil
	c.instructions = instructionsBackup
	c.pushInstruction(vm.NewCode(vm.NewFunctionOP, errors.UnknownLine, [2]int{len(functionCode) + 2, len(functionDefinition.Arguments)}))
	var arguments []string
	for _, argument := range functionDefinition.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	c.pushInstruction(vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	c.extendInstructions(functionCode)
	c.pushInstruction(vm.NewCode(vm.ReturnOP, errors.UnknownLine, 0))
	c.pushInstruction(vm.NewCode(vm.AssignIdentifierOP, functionDefinition.Name.Token.Line, functionDefinition.Name.Token.String))
	return nil
}

func (c *Compiler) compileReturnStatement(returnStatement *ast.ReturnStatement) *errors.Error {
	numberOfResults := len(returnStatement.Results)
	for i := numberOfResults - 1; i > -1; i-- {
		resultCompilationError := c.compileExpression(returnStatement.Results[i])
		if resultCompilationError != nil {
			return resultCompilationError
		}
	}
	c.pushInstruction(vm.NewCode(vm.ReturnOP, errors.UnknownLine, numberOfResults))
	return nil
}

func (c *Compiler) compileStatement(statement ast.Statement) *errors.Error {
	switch statement.(type) {
	case *ast.AssignStatement:
		return c.compileAssignStatement(statement.(*ast.AssignStatement))
	case *ast.FunctionDefinitionStatement:
		return c.compileFunctionDefinition(statement.(*ast.FunctionDefinitionStatement))
	case *ast.ReturnStatement:
		return c.compileReturnStatement(statement.(*ast.ReturnStatement))
	}
	return nil
}

func (c *Compiler) compile(node ast.Node) *errors.Error {
	switch node.(type) {
	case ast.Expression:
		return c.compileExpression(node.(ast.Expression))
	case ast.Statement:
		return c.compileStatement(node.(ast.Statement))
	}
	return nil
}

func (c *Compiler) compileBody(body []ast.Node) *errors.Error {
	for _, node := range body {
		compileError := c.compile(node)
		if compileError != nil {
			return compileError
		}
	}
	return nil
}

func (c *Compiler) CompileToArray() ([]vm.Code, *errors.Error) {
	codeAst, parsingError := c.parser.Parse()
	if parsingError != nil {
		return nil, parsingError
	}

	compileError := c.compileBegin(codeAst.Begin)
	if compileError != nil {
		return nil, compileError
	}

	compileError = c.compileBody(codeAst.Body)
	if compileError != nil {
		return nil, compileError
	}

	compileError = c.compileEnd(codeAst.End)
	if compileError != nil {
		return nil, compileError
	}
	return c.instructions, nil
}

func (c *Compiler) Compile() (*vm.Bytecode, *errors.Error) {
	_, compilationError := c.CompileToArray()
	if compilationError != nil {
		return nil, compilationError
	}
	return vm.NewBytecodeFromArray(c.instructions), nil
}

func NewCompiler(codeReader reader.Reader, ) *Compiler {
	return &Compiler{
		parser:       parser.NewParser(lexer.NewLexer(codeReader)),
		instructions: nil,
		index:        -1,
		length:       0,
	}
}
