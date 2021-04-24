package parser

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
	"testing"
)

func walker(node ast.Node, deep int) string {
	switch node.(type) {
	case *ast.Program:
		result := ""
		for _, child := range node.(*ast.Program).Body {
			result += walker(child, deep+1)
			result += "\n"
		}
		return result
	case *ast.BinaryExpression:
		return walker(node.(*ast.BinaryExpression).LeftHandSide, deep+1) +
			" " + node.(*ast.BinaryExpression).Operator +
			" " + walker(node.(*ast.BinaryExpression).RightHandSide, deep+1)
	case *ast.BasicLiteralExpression:
		return node.(*ast.BasicLiteralExpression).String
	case *ast.UnaryExpression:
		return node.(*ast.UnaryExpression).Operator + walker(node.(*ast.UnaryExpression).X, deep+1)
	case *ast.SelectorExpression:
		return walker(node.(*ast.SelectorExpression).X, deep+1) + "." + node.(*ast.SelectorExpression).Identifier.String
	case *ast.Identifier:
		return node.(*ast.Identifier).String
	case *ast.MethodInvocationExpression:
		result := walker(node.(*ast.MethodInvocationExpression).Function, deep+1) + "("
		for index, child := range node.(*ast.MethodInvocationExpression).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(child, deep+1)
		}
		return result + ")"
	case *ast.IndexExpression:
		result := walker(node.(*ast.IndexExpression).Source, deep+1) + "["
		result += walker(node.(*ast.IndexExpression).Index[0], deep+1)
		if node.(*ast.IndexExpression).Index[1] != nil {
			result += ":" + walker(node.(*ast.IndexExpression).Index[1], deep+1)
		}
		return result + "]"
	case *ast.LambdaExpression:
		result := "lambda "
		for index, argument := range node.(*ast.LambdaExpression).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument, deep+1)
		}
		result += ": "
		return result + walker(node.(*ast.LambdaExpression).Code, deep+1)
	case *ast.ParenthesesExpression:
		return "(" + walker(node.(*ast.ParenthesesExpression).X, deep+1) + ")"
	case *ast.TupleExpression:
		result := "("
		for index, value := range node.(*ast.TupleExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value, deep+1)
		}
		return result + ")"
	case *ast.ArrayExpression:
		result := "["
		for index, value := range node.(*ast.ArrayExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value, deep+1)
		}
		return result + "]"
	case *ast.HashExpression:
		result := "{"
		for index, value := range node.(*ast.HashExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value.Key, deep+1)
			result += ":" + walker(value.Value, deep+1)
		}
		return result + "}"
	case *ast.OneLineIfExpression:
		result := walker(node.(*ast.OneLineIfExpression).Result, deep+1)
		result += " if " + walker(node.(*ast.OneLineIfExpression).Condition, deep+1)
		if node.(*ast.OneLineIfExpression).ElseResult != nil {
			result += " else " + walker(node.(*ast.OneLineIfExpression).ElseResult, deep+1)
		}
		return result
	case *ast.OneLineUnlessExpression:
		result := walker(node.(*ast.OneLineUnlessExpression).Result, deep+1)
		result += " unless "
		result += walker(node.(*ast.OneLineUnlessExpression).Condition, deep+1)
		if node.(*ast.OneLineUnlessExpression).ElseResult != nil {
			result += " else " + walker(node.(*ast.OneLineUnlessExpression).ElseResult, deep+1)
		}
		return result
	case *ast.GeneratorExpression:
		result := walker(node.(*ast.GeneratorExpression).Operation, deep+1)
		result += " for "
		for index, variable := range node.(*ast.GeneratorExpression).Variables {
			if index != 0 {
				result += ", "
			}
			result += walker(variable, deep+1)
		}
		result += " in "
		return result + walker(node.(*ast.GeneratorExpression).Source, deep+1)
	case *ast.AssignStatement:
		result := walker(node.(*ast.AssignStatement).LeftHandSide, deep+1)
		result += " " + node.(*ast.AssignStatement).AssignOperator + " "
		return result + walker(node.(*ast.AssignStatement).RightHandSide, deep+1)
	case *ast.RetryStatement:
		return "retry"
	case *ast.BreakStatement:
		return "break"
	case *ast.RedoStatement:
		return "redo"
	case *ast.PassStatement:
		return "pass"
	case *ast.YieldStatement:
		result := "yield "
		for index, output := range node.(*ast.YieldStatement).Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output, deep+1)
		}
		return result
	case *ast.ReturnStatement:
		result := "return "
		for index, output := range node.(*ast.ReturnStatement).Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output, deep+1)
		}
		return result
	case *ast.GoStatement:
		return "go " + walker(node.(*ast.GoStatement).X, deep+1)
	case *ast.SuperInvocationStatement:
		result := "super("
		for index, argument := range node.(*ast.SuperInvocationStatement).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument, deep+1)
		}
		return result + ")"
	}
	panic("unknown node type")
}

func walk(sample int, node ast.Node) {
	fmt.Printf("\nSample: %d\n", sample)
	fmt.Print(walker(node, 0))
}

func test(t *testing.T, samples []string) {
	for sampleIndex, sample := range samples {
		lex := lexer.NewLexer(sample)
		parser, parserCreationError := NewParser(lex)
		if parserCreationError != nil {
			t.Error(parserCreationError)
			return
		}
		program, parsingError := parser.Parse()
		if parsingError != nil {
			t.Error(parsingError)
			return
		}
		walk(sampleIndex+1, program)
	}
}

var basicSamples = []string{
	"1 + 2 * 3",
	"1 * 2 + 3",
	"1 >= 2 == 3 - 4 + 5 - 6 / 7 / 8 ** 9 - 10",
	"5--5",
	"hello.world.carro",
	"1.4.hello.world()",
	"hello(1)",
	"'Hello world'.index(str(12345)[010])",
	"'Hello world'.index(str(12345)[0:10])",
	"lambda x, y, z: print(x, y - z)",
	"lambda x: print((1+2) * 3)",
	"(1, 2, (3, (4, (((4+1))))))",
	"[1]",
	"{1: (1*2)/4, 2: 282}",
	"(1 if x < 4 else 2)",
	"True",
	"1 if True",
	"1 unless False",
	"(1 for 2 in (3, 4))",
	"\n\n\n\n\n\n\n1\n2\n3\n\n\n\n\n\n\n\n[4,\n\n\n5+\n6!=\n11]",
	"a = 234",
	"a[1] ~= 234",
	"2.a += [1]",
	"redo",
	"yield 1",
	"yield 1, 2 + 4, lambda x: x + 2, (1, 2 , 3, 4)",
	"return 1",
	"return 1, 2 + 4, lambda x: x + 2, (1, 2 , 3, 4)",
	"go super_duper()",
	"super(1)",
	"super(1, 2)",
	"super(1, 2, call((1, 2, 3, 4, 2 * (5 - 3))))",
}

func TestParseBasic(t *testing.T) {
	test(t, basicSamples)
}
