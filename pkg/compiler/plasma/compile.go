package plasma

import (
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
	parser      *parser.Parser
	programCode []vm.Code
}

func (c *Compiler) pushInstruction(code vm.Code) {
	c.programCode = append(c.programCode, code)
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
	case lexer.Float:
		break
	case lexer.ScientificFloat:
		break
	case lexer.SingleQuoteString, lexer.DoubleQuoteString:
		break
	case lexer.ByteString:
		break
	case lexer.True:
		break
	case lexer.False:
		break
	case lexer.None:
		break
	}
	return nil
}

func (c *Compiler) compile(node ast.Node) *errors.Error {
	switch node.(type) {
	case *ast.BasicLiteralExpression:
		return c.compileLiteral(node.(*ast.BasicLiteralExpression))
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

func (c *Compiler) Compile() (*vm.Bytecode, *errors.Error) {
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

	c.pushInstruction(vm.NewCode(vm.ReturnOP, errors.UnknownLine, nil))

	return vm.NewBytecodeFromArray(c.programCode), nil
}

func NewCompiler(codeReader reader.Reader, ) *Compiler {
	return &Compiler{
		parser:      parser.NewParser(lexer.NewLexer(codeReader)),
		programCode: nil,
	}
}
