package parser

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
	"strings"
	"testing"
)

func walker(node ast.Node) string {
	switch node.(type) {
	case *ast.Program:
		result := ""
		for _, child := range node.(*ast.Program).Body {
			result += walker(child)
			result += "\n"
		}
		return result
	case *ast.BinaryExpression:
		return walker(node.(*ast.BinaryExpression).LeftHandSide) +
			" " + node.(*ast.BinaryExpression).Operator.String +
			" " + walker(node.(*ast.BinaryExpression).RightHandSide)
	case *ast.BasicLiteralExpression:
		return node.(*ast.BasicLiteralExpression).Token.String
	case *ast.UnaryExpression:
		if node.(*ast.UnaryExpression).Operator.DirectValue == lexer.Not {
			return node.(*ast.UnaryExpression).Operator.String + " " + walker(node.(*ast.UnaryExpression).X)
		}
		return node.(*ast.UnaryExpression).Operator.String + walker(node.(*ast.UnaryExpression).X)
	case *ast.SelectorExpression:
		return walker(node.(*ast.SelectorExpression).X) + "." + node.(*ast.SelectorExpression).Identifier.Token.String
	case *ast.Identifier:
		return node.(*ast.Identifier).Token.String
	case *ast.MethodInvocationExpression:
		result := walker(node.(*ast.MethodInvocationExpression).Function) + "("
		for index, child := range node.(*ast.MethodInvocationExpression).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(child)
		}
		return result + ")"
	case *ast.IndexExpression:
		result := walker(node.(*ast.IndexExpression).Source) + "["
		result += walker(node.(*ast.IndexExpression).Index[0])
		if node.(*ast.IndexExpression).Index[1] != nil {
			result += ":" + walker(node.(*ast.IndexExpression).Index[1])
		}
		return result + "]"
	case *ast.LambdaExpression:
		result := "lambda "
		for index, argument := range node.(*ast.LambdaExpression).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ": "
		return result + walker(node.(*ast.LambdaExpression).Code)
	case *ast.ParenthesesExpression:
		return "(" + walker(node.(*ast.ParenthesesExpression).X) + ")"
	case *ast.TupleExpression:
		result := "("
		for index, value := range node.(*ast.TupleExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + ")"
	case *ast.ArrayExpression:
		result := "["
		for index, value := range node.(*ast.ArrayExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + "]"
	case *ast.HashExpression:
		result := "{"
		for index, value := range node.(*ast.HashExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value.Key)
			result += ": " + walker(value.Value)
		}
		return result + "}"
	case *ast.OneLineIfExpression:
		result := walker(node.(*ast.OneLineIfExpression).Result)
		result += " if " + walker(node.(*ast.OneLineIfExpression).Condition)
		if node.(*ast.OneLineIfExpression).ElseResult != nil {
			result += " else " + walker(node.(*ast.OneLineIfExpression).ElseResult)
		}
		return result
	case *ast.OneLineUnlessExpression:
		result := walker(node.(*ast.OneLineUnlessExpression).Result)
		result += " unless "
		result += walker(node.(*ast.OneLineUnlessExpression).Condition)
		if node.(*ast.OneLineUnlessExpression).ElseResult != nil {
			result += " else " + walker(node.(*ast.OneLineUnlessExpression).ElseResult)
		}
		return result
	case *ast.GeneratorExpression:
		result := walker(node.(*ast.GeneratorExpression).Operation)
		result += " for "
		for index, variable := range node.(*ast.GeneratorExpression).Variables {
			if index != 0 {
				result += ", "
			}
			result += walker(variable)
		}
		result += " in "
		return result + walker(node.(*ast.GeneratorExpression).Source)
	case *ast.AssignStatement:
		result := walker(node.(*ast.AssignStatement).LeftHandSide)
		result += " " + node.(*ast.AssignStatement).AssignOperator.String + " "
		return result + walker(node.(*ast.AssignStatement).RightHandSide)
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
			result += walker(output)
		}
		return result
	case *ast.ReturnStatement:
		result := "return "
		for index, output := range node.(*ast.ReturnStatement).Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast.DeferStatement:
		return "defer " + walker(node.(*ast.DeferStatement).X)
	case *ast.GoStatement:
		return "go " + walker(node.(*ast.GoStatement).X)
	case *ast.SuperInvocationStatement:
		result := "super("
		for index, argument := range node.(*ast.SuperInvocationStatement).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		return result + ")"
	case *ast.IfStatement:
		result := "if "
		result += walker(node.(*ast.IfStatement).Condition)
		for _, bodyNode := range node.(*ast.IfStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		for _, elifBlock := range node.(*ast.IfStatement).ElifBlocks {
			result += "\nelif "
			result += walker(elifBlock.Condition)
			for _, bodyNode := range elifBlock.Body {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if node.(*ast.IfStatement).Else != nil {
			result += "\nelse"
			for _, elseNode := range node.(*ast.IfStatement).Else {
				nodeString := walker(elseNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast.UnlessStatement:
		result := "unless "
		result += walker(node.(*ast.UnlessStatement).Condition)
		for _, bodyNode := range node.(*ast.UnlessStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		for _, elifBlock := range node.(*ast.UnlessStatement).ElifBlocks {
			result += "\nelif "
			result += walker(elifBlock.Condition)
			for _, bodyNode := range elifBlock.Body {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if node.(*ast.UnlessStatement).Else != nil {
			result += "\nelse"
			for _, elseNode := range node.(*ast.UnlessStatement).Else {
				nodeString := walker(elseNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast.EnumStatement:
		result := "enum " + walker(node.(*ast.EnumStatement).Name)
		for _, identifier := range node.(*ast.EnumStatement).EnumIdentifiers {
			result += "\n\t" + walker(identifier)
		}
		return result + "\nend"
	case *ast.StructStatement:
		result := "struct " + walker(node.(*ast.StructStatement).Name)
		for _, identifier := range node.(*ast.StructStatement).Fields {
			result += "\n\t" + walker(identifier)
		}
		return result + "\nend"
	case *ast.SwitchStatement:
		result := "switch " + walker(node.(*ast.SwitchStatement).Target)
		for _, caseBlock := range node.(*ast.SwitchStatement).CaseBlocks {
			result += "\ncase "
			for index, caseTarget := range caseBlock.Cases {
				if index != 0 {
					result += ", "
				}
				result += walker(caseTarget)
			}
			for _, caseChild := range caseBlock.Body {
				nodeString := walker(caseChild)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if node.(*ast.SwitchStatement).Else != nil {
			result += "\nelse"
			for _, elseChild := range node.(*ast.SwitchStatement).Else {
				nodeString := walker(elseChild)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast.WhileLoopStatement:
		result := "while " + walker(node.(*ast.WhileLoopStatement).Condition)
		for _, child := range node.(*ast.WhileLoopStatement).Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	}
	panic("unknown node type")
}

func walk(node ast.Node) string {
	return walker(node)
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
		result := walk(program)
		result = result[:len(result)-1]
		if result == sample {
			t.Logf("\nSample: %d -> SUCCESS", sampleIndex+1)
		} else {
			t.Errorf("\nSample: %d -> FAIL", sampleIndex+1)
			fmt.Println(sample)
			fmt.Println(result)
		}
	}
}

var basicSamples = []string{
	"1 + 2 * 3",
	"1 * 2 + 3",
	"1 >= 2 == 3 - 4 + 5 - 6 / 7 / 8 ** 9 - 10",
	"5 - -5",
	"hello.world.carro",
	"1.4.hello.world()",
	"hello(1)",
	"'Hello world'.index(str(12345)[010])",
	"'Hello world'.index(str(12345)[0:10])",
	"lambda x, y, z: print(x, y - z)",
	"lambda x: print((1 + 2) * 3)",
	"(1, 2, (3, (4, (((4 + 1))))))",
	"[1]",
	"{1: (1 * 2) / 4, 2: 282}",
	"(1 if x < 4 else 2)",
	"True",
	"not True",
	"1 if True",
	"!True",
	"1 unless False",
	"1 in (1, 2, 3)",
	"(1 for 2 in (3, 4))",
	"1\n2\n3\n[4, 5 + 6 != 11]",
	"a = 234",
	"a[1] ~= 234",
	"2.a += [1]",
	"a and b",
	"a xor b",
	"a or not b",
	"redo",
	"yield 1",
	"yield 1, 2 + 4, lambda x: x + 2, (1, 2, 3, 4)",
	"return 1",
	"return 1, 2 + 4, lambda x: x + 2, (1, 2, 3, 4)",
	"go super_duper()",
	"defer a()",
	"super(1)",
	"super(1, 2)",
	"super(1, 2, call((1, 2, 3, 4, 2 * (5 - 3))))",
	"if a > 2\n" +
		"\tcall()\n" +
		"elif a < 2\n" +
		"\tif a == 0\n" +
		"\t\tprint(\"\\\"a\\\" is zero\")\n" +
		"\telse\n" +
		"\t\tprint(\"\\\"a\\\" is non zero\")" +
		"\n\tend\n" +
		"end",
	"unless a > 2\n" +
		"\tcall()\n" +
		"elif a < 2\n" +
		"\tif 1 if a < 2 else None\n" +
		"\t\tprint(\"\\\"a\\\" is zero\")\n" +
		"\telse\n" +
		"\t\tprint(\"\\\"a\\\" is non zero\")\n" +
		"\tend\n" +
		"\tif 1 == 2\n" +
		"\t\tprint(2)\n" +
		"\tend\n" +
		"end",
	"enum Tokens\n" +
		"\tString\n" +
		"\tFloat\n" +
		"\tInteger\n" +
		"end",
	"struct ListNode\n" +
		"\tValue\n" +
		"\tLeft\n" +
		"\tRight\n" +
		"end",
	"switch Token.Kind\n" +
		"case Numeric\n" +
		"\tbreak\n" +
		"case String\n" +
		"\tprint(\"I am a String\")\n" +
		"end",
	"while True\n" +
		"\tif a > b\n" +
		"\t\tbreak\n" +
		"\tend\n" +
		"\ta += 1\n" +
		"\tb -= 1\n" +
		"end",
}

func TestParseBasic(t *testing.T) {
	test(t, basicSamples)
}
