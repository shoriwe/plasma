package parser

import (
	"fmt"
	ast2 "github.com/shoriwe/gplasma/pkg/ast"
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
	reader2 "github.com/shoriwe/gplasma/pkg/reader"
	"reflect"
	"strings"
	"testing"
)

func walker(node ast2.Node) string {
	switch node.(type) {
	case *ast2.Program:
		result := ""
		if node.(*ast2.Program).Begin != nil {
			result += walker(node.(*ast2.Program).Begin)
		}
		if node.(*ast2.Program).End != nil {
			if node.(*ast2.Program).End != nil {
				result += "\n"
			}
			result += walker(node.(*ast2.Program).End)
		}
		if node.(*ast2.Program).Begin != nil || node.(*ast2.Program).End != nil {
			result += "\n"
		}
		if len(node.(*ast2.Program).Body) > 0 {
			for _, child := range node.(*ast2.Program).Body {
				result += walker(child)
				result += "\n"
			}
		}
		return result
	case *ast2.BinaryExpression:
		return walker(node.(*ast2.BinaryExpression).LeftHandSide) +
			" " + node.(*ast2.BinaryExpression).Operator.String() +
			" " + walker(node.(*ast2.BinaryExpression).RightHandSide)
	case *ast2.BasicLiteralExpression:
		return node.(*ast2.BasicLiteralExpression).Token.String()
	case *ast2.UnaryExpression:
		if node.(*ast2.UnaryExpression).Operator.DirectValue == lexer2.Not {
			return node.(*ast2.UnaryExpression).Operator.String() + " " + walker(node.(*ast2.UnaryExpression).X)
		}
		return node.(*ast2.UnaryExpression).Operator.String() + walker(node.(*ast2.UnaryExpression).X)
	case *ast2.SelectorExpression:
		return walker(node.(*ast2.SelectorExpression).X) + "." + node.(*ast2.SelectorExpression).Identifier.Token.String()
	case *ast2.Identifier:
		return node.(*ast2.Identifier).Token.String()
	case *ast2.MethodInvocationExpression:
		result := walker(node.(*ast2.MethodInvocationExpression).Function) + "("
		for index, child := range node.(*ast2.MethodInvocationExpression).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(child)
		}
		return result + ")"
	case *ast2.IndexExpression:
		result := walker(node.(*ast2.IndexExpression).Source) + "["
		result += walker(node.(*ast2.IndexExpression).Index)
		return result + "]"
	case *ast2.LambdaExpression:
		result := "lambda "
		for index, argument := range node.(*ast2.LambdaExpression).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ": "
		return result + walker(node.(*ast2.LambdaExpression).Code)
	case *ast2.ParenthesesExpression:
		return "(" + walker(node.(*ast2.ParenthesesExpression).X) + ")"
	case *ast2.TupleExpression:
		result := "("
		for index, value := range node.(*ast2.TupleExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + ")"
	case *ast2.ArrayExpression:
		result := "["
		for index, value := range node.(*ast2.ArrayExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + "]"
	case *ast2.HashExpression:
		result := "{"
		for index, value := range node.(*ast2.HashExpression).Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value.Key)
			result += ": " + walker(value.Value)
		}
		return result + "}"
	case *ast2.IfOneLinerExpression:
		result := walker(node.(*ast2.IfOneLinerExpression).Result)
		result += " if " + walker(node.(*ast2.IfOneLinerExpression).Condition)
		if node.(*ast2.IfOneLinerExpression).ElseResult != nil {
			result += " else " + walker(node.(*ast2.IfOneLinerExpression).ElseResult)
		}
		return result
	case *ast2.UnlessOneLinerExpression:
		result := walker(node.(*ast2.UnlessOneLinerExpression).Result)
		result += " unless "
		result += walker(node.(*ast2.UnlessOneLinerExpression).Condition)
		if node.(*ast2.UnlessOneLinerExpression).ElseResult != nil {
			result += " else " + walker(node.(*ast2.UnlessOneLinerExpression).ElseResult)
		}
		return result
	case *ast2.GeneratorExpression:
		result := walker(node.(*ast2.GeneratorExpression).Operation)
		result += " for "
		for index, variable := range node.(*ast2.GeneratorExpression).Receivers {
			if index != 0 {
				result += ", "
			}
			result += walker(variable)
		}
		result += " in "
		return "(" + result + walker(node.(*ast2.GeneratorExpression).Source) + ")"
	case *ast2.AssignStatement:
		result := walker(node.(*ast2.AssignStatement).LeftHandSide)
		result += " " + node.(*ast2.AssignStatement).AssignOperator.String() + " "
		return result + walker(node.(*ast2.AssignStatement).RightHandSide)
	case *ast2.ContinueStatement:
		return "continue"
	case *ast2.BreakStatement:
		return "break"
	case *ast2.RedoStatement:
		return "redo"
	case *ast2.PassStatement:
		return "pass"
	case *ast2.YieldStatement:
		result := "yield "
		for index, output := range node.(*ast2.YieldStatement).Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast2.ReturnStatement:
		result := "return "
		for index, output := range node.(*ast2.ReturnStatement).Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast2.IfStatement:
		result := "if "
		result += walker(node.(*ast2.IfStatement).Condition)
		for _, bodyNode := range node.(*ast2.IfStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		if len(node.(*ast2.IfStatement).Else) > 0 {
			if _, isIf := node.(*ast2.IfStatement).Else[0].(*ast2.IfStatement); len(node.(*ast2.IfStatement).Else) == 1 && isIf {
				elifBlock := node.(*ast2.IfStatement).Else[0].(*ast2.IfStatement)
				result += "\nelif "
				result += walker(elifBlock.Condition)
				for _, bodyNode := range elifBlock.Body {
					nodeString := walker(bodyNode)
					nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
					result += "\n\t" + nodeString
				}
			} else {
				result += "\nelse"
				for _, elseNode := range node.(*ast2.IfStatement).Else {
					nodeString := walker(elseNode)
					nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
					result += "\n\t" + nodeString
				}
			}
		}
		return result + "\nend"
	case *ast2.UnlessStatement:
		result := "unless "
		result += walker(node.(*ast2.UnlessStatement).Condition)
		for _, bodyNode := range node.(*ast2.UnlessStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		if len(node.(*ast2.UnlessStatement).Else) > 0 {
			if _, isUnless := node.(*ast2.UnlessStatement).Else[0].(*ast2.UnlessStatement); len(node.(*ast2.UnlessStatement).Else) == 1 && isUnless {
				elifBlock := node.(*ast2.UnlessStatement).Else[0].(*ast2.UnlessStatement)
				result += "\nelif "
				result += walker(elifBlock.Condition)
				for _, bodyNode := range elifBlock.Body {
					nodeString := walker(bodyNode)
					nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
					result += "\n\t" + nodeString
				}
			} else {
				result += "\nelse"
				for _, elseNode := range node.(*ast2.UnlessStatement).Else {
					nodeString := walker(elseNode)
					nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
					result += "\n\t" + nodeString
				}
			}
		}
		return result + "\nend"
	case *ast2.SwitchStatement:
		result := "switch " + walker(node.(*ast2.SwitchStatement).Target)
		for _, caseBlock := range node.(*ast2.SwitchStatement).CaseBlocks {
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
		if node.(*ast2.SwitchStatement).Default != nil {
			result += "\ndefault"
			for _, elseChild := range node.(*ast2.SwitchStatement).Default {
				nodeString := walker(elseChild)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast2.WhileLoopStatement:
		result := "while " + walker(node.(*ast2.WhileLoopStatement).Condition)
		for _, child := range node.(*ast2.WhileLoopStatement).Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.UntilLoopStatement:
		result := "until " + walker(node.(*ast2.UntilLoopStatement).Condition)
		for _, child := range node.(*ast2.UntilLoopStatement).Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.ForLoopStatement:
		result := "for "
		for index, receiver := range node.(*ast2.ForLoopStatement).Receivers {
			if index != 0 {
				result += ", "
			}
			result += walker(receiver)
		}
		result += " in " + walker(node.(*ast2.ForLoopStatement).Source)
		for _, bodyNode := range node.(*ast2.ForLoopStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.ModuleStatement:
		result := "module " + walker(node.(*ast2.ModuleStatement).Name)
		for _, bodyNode := range node.(*ast2.ModuleStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.ClassStatement:
		result := "class " + walker(node.(*ast2.ClassStatement).Name)
		if node.(*ast2.ClassStatement).Bases != nil {
			result += "("
			for index, base := range node.(*ast2.ClassStatement).Bases {
				if index != 0 {
					result += ", "
				}
				result += walker(base)
			}
			result += ")"
		}
		for _, bodyNode := range node.(*ast2.ClassStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.InterfaceStatement:
		result := "interface " + walker(node.(*ast2.InterfaceStatement).Name)
		if node.(*ast2.InterfaceStatement).Bases != nil {
			result += "("
			for index, base := range node.(*ast2.InterfaceStatement).Bases {
				if index != 0 {
					result += ", "
				}
				result += walker(base)
			}
			result += ")"
		}
		for _, bodyNode := range node.(*ast2.InterfaceStatement).MethodDefinitions {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.FunctionDefinitionStatement:
		result := "def " + walker(node.(*ast2.FunctionDefinitionStatement).Name)
		result += "("
		for index, argument := range node.(*ast2.FunctionDefinitionStatement).Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ")"
		for _, bodyNode := range node.(*ast2.FunctionDefinitionStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.EndStatement:
		result := "END"
		for _, bodyNode := range node.(*ast2.EndStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.BeginStatement:
		result := "BEGIN"
		for _, bodyNode := range node.(*ast2.BeginStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.RaiseStatement:
		return "raise " + walker(node.(*ast2.RaiseStatement).X)
	case *ast2.TryStatement:
		result := "try"
		for _, bodyNode := range node.(*ast2.TryStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		for _, exceptBlock := range node.(*ast2.TryStatement).ExceptBlocks {
			result += "\nexcept "
			for index, target := range exceptBlock.Targets {
				if index != 0 {
					result += ", "
				}
				result += walker(target)
			}
			if exceptBlock.CaptureName != nil {
				result += " as " + walker(exceptBlock.CaptureName)
			}
			for _, bodyNode := range exceptBlock.Body {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if node.(*ast2.TryStatement).Else != nil {
			result += "\nelse"
			for _, bodyNode := range node.(*ast2.TryStatement).Else {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if node.(*ast2.TryStatement).Finally != nil {
			result += "\nfinally"
			for _, bodyNode := range node.(*ast2.TryStatement).Finally {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast2.DoWhileStatement:
		result := "do"
		for _, bodyNode := range node.(*ast2.DoWhileStatement).Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nwhile " + walker(node.(*ast2.DoWhileStatement).Condition)
	}
	panic("unknown node type: " + reflect.TypeOf(node).String())
}

func walk(node ast2.Node) string {
	return walker(node)
}

func test(t *testing.T, samples []string) {
	for sampleIndex, sample := range samples {
		lex := lexer2.NewLexer(reader2.NewStringReader(sample))
		parser := NewParser(lex)
		program, parsingError := parser.Parse()
		if parsingError != nil {
			t.Error(fmt.Sprintf("%s in sample %d", parsingError.Error(), sampleIndex))
			return
		}
		result := walk(program)
		if len(result) != 0 {
			result = result[:len(result)-1]
		}
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
	"'Hello world'.index(str(12345)[(0, 10)])",
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
	"(1 for a in (3, 4))",
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
	"switch Token.Kind\n" +
		"case Numeric, CommandOutput\n" +
		"\tbreak\n" +
		"case String\n" +
		"\tprint(\"I am a String\")\n" +
		"default\n" +
		"\tprint(\"errors\")\n" +
		"end",
	"while True\n" +
		"\tif a > b\n" +
		"\t\tbreak\n" +
		"\tend\n" +
		"\ta += 1\n" +
		"\tb -= 1\n" +
		"end",
	"[]",
	"for a, b, c in range(10)\n" +
		"\tprint(\"hello world!\")\n" +
		"end",
	"until (1 + 2 * 3 / a) > 5\n" +
		"\tif a > b\n" +
		"\t\tbreak\n" +
		"\tend\n" +
		"\ta += 1\n" +
		"\tb -= 1\n" +
		"end",
	"module something\n" +
		"\tclass Hello(a.b, a, c, Hello2)\n" +
		"\tend\n" +
		"\tclass Hello2(IHello)\n" +
		"\tend\n" +
		"\tinterface IHello\n" +
		"\t\tdef SayHello()\n" +
		"\t\t\tprint(\"Hello\")\n" +
		"\t\tend\n" +
		"\t\tdef SayHello()\n" +
		"\t\t\tprint(\"Hello\")\n" +
		"\t\tend\n" +
		"\tend\n" +
		"end",
	"BEGIN\n" +
		"\tNode = (1, 2)\n" +
		"end",
	"try\n" +
		"\tprint(variable)\n" +
		"except UndefinedIdentifier, AnyException as errors\n" +
		"\tprint(errors)\n" +
		"except NoToStringException as errors\n" +
		"\tprint(errors)\n" +
		"else\n" +
		"\tprint(\"Unknown *errors\")\n" +
		"\traise UnknownException()\n" +
		"finally\n" +
		"\tprint(\"Done\")\n" +
		"end",
	"do\n" +
		"\tprint(\"Hello\")\n" +
		"while a > b",
	"def fib(n)\n" +
		"\tif n == 0\n" +
		"\t\treturn 0\n" +
		"\tend\n" +
		"\tif n == 1\n" +
		"\t\treturn 1\n" +
		"\tend\n" +
		"\treturn fib(n - 1) + fib(n - 2)\n" +
		"end\n" +
		"println(fib(35))",
}

func TestParseBasic(t *testing.T) {
	test(t, basicSamples)
}
