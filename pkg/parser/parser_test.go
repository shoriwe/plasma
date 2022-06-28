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
	switch n := node.(type) {
	case *ast2.Program:
		result := ""
		if n.Begin != nil {
			result += walker(n.Begin)
		}
		if n.End != nil {
			if n.End != nil {
				result += "\n"
			}
			result += walker(n.End)
		}
		if n.Begin != nil || n.End != nil {
			result += "\n"
		}
		if len(n.Body) > 0 {
			for _, child := range n.Body {
				result += walker(child)
				result += "\n"
			}
		}
		return result
	case *ast2.BinaryExpression:
		return walker(n.LeftHandSide) +
			" " + n.Operator.String() +
			" " + walker(n.RightHandSide)
	case *ast2.BasicLiteralExpression:
		return n.Token.String()
	case *ast2.UnaryExpression:
		if n.Operator.DirectValue == lexer2.Not {
			return n.Operator.String() + " " + walker(n.X)
		}
		return n.Operator.String() + walker(n.X)
	case *ast2.SelectorExpression:
		return walker(n.X) + "." + n.Identifier.Token.String()
	case *ast2.Identifier:
		return n.Token.String()
	case *ast2.MethodInvocationExpression:
		result := walker(n.Function) + "("
		for index, child := range n.Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(child)
		}
		return result + ")"
	case *ast2.IndexExpression:
		result := walker(n.Source) + "["
		result += walker(n.Index)
		return result + "]"
	case *ast2.LambdaExpression:
		result := "lambda "
		for index, argument := range n.Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ": "
		return result + walker(n.Code)
	case *ast2.ParenthesesExpression:
		return "(" + walker(n.X) + ")"
	case *ast2.TupleExpression:
		result := "("
		for index, value := range n.Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + ")"
	case *ast2.ArrayExpression:
		result := "["
		for index, value := range n.Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + "]"
	case *ast2.HashExpression:
		result := "{"
		for index, value := range n.Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value.Key)
			result += ": " + walker(value.Value)
		}
		return result + "}"
	case *ast2.IfOneLinerExpression:
		result := walker(n.Result)
		result += " if " + walker(n.Condition)
		if n.ElseResult != nil {
			result += " else " + walker(n.ElseResult)
		}
		return result
	case *ast2.UnlessOneLinerExpression:
		result := walker(n.Result)
		result += " unless "
		result += walker(n.Condition)
		if n.ElseResult != nil {
			result += " else " + walker(n.ElseResult)
		}
		return result
	case *ast2.GeneratorExpression:
		result := walker(n.Operation)
		result += " for "
		for index, variable := range n.Receivers {
			if index != 0 {
				result += ", "
			}
			result += walker(variable)
		}
		result += " in "
		return "(" + result + walker(n.Source) + ")"
	case *ast2.AssignStatement:
		result := walker(n.LeftHandSide)
		result += " " + n.AssignOperator.String() + " "
		return result + walker(n.RightHandSide)
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
		for index, output := range n.Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast2.ReturnStatement:
		result := "return "
		for index, output := range n.Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast2.IfStatement:
		result := "if "
		result += walker(n.Condition)
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		for _, elifBlock := range n.ElifBlocks {
			result += "\nelif " + walker(elifBlock.Condition)
			for _, bodyNode := range elifBlock.Body {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if len(n.Else) > 0 {
			result += "\nelse"
			for _, elseNode := range n.Else {
				nodeString := walker(elseNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast2.UnlessStatement:
		result := "unless "
		result += walker(n.Condition)
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		for _, elifBlock := range n.ElifBlocks {
			result += "\nelif " + walker(elifBlock.Condition)
			for _, bodyNode := range elifBlock.Body {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if len(n.Else) > 0 {
			result += "\nelse"
			for _, elseNode := range n.Else {
				nodeString := walker(elseNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast2.SwitchStatement:
		result := "switch " + walker(n.Target)
		for _, caseBlock := range n.CaseBlocks {
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
		if n.Default != nil {
			result += "\ndefault"
			for _, elseChild := range n.Default {
				nodeString := walker(elseChild)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast2.WhileLoopStatement:
		result := "while " + walker(n.Condition)
		for _, child := range n.Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.UntilLoopStatement:
		result := "until " + walker(n.Condition)
		for _, child := range n.Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.ForLoopStatement:
		result := "for "
		for index, receiver := range n.Receivers {
			if index != 0 {
				result += ", "
			}
			result += walker(receiver)
		}
		result += " in " + walker(n.Source)
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.ModuleStatement:
		result := "module " + walker(n.Name)
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.ClassStatement:
		result := "class " + walker(n.Name)
		if n.Bases != nil {
			result += "("
			for index, base := range n.Bases {
				if index != 0 {
					result += ", "
				}
				result += walker(base)
			}
			result += ")"
		}
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.InterfaceStatement:
		result := "interface " + walker(n.Name)
		if n.Bases != nil {
			result += "("
			for index, base := range n.Bases {
				if index != 0 {
					result += ", "
				}
				result += walker(base)
			}
			result += ")"
		}
		for _, bodyNode := range n.MethodDefinitions {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.FunctionDefinitionStatement:
		result := "def " + walker(n.Name)
		result += "("
		for index, argument := range n.Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ")"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.GeneratorDefinitionStatement:
		result := "gen " + walker(n.Name)
		result += "("
		for index, argument := range n.Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ")"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.EndStatement:
		result := "END"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.BeginStatement:
		result := "BEGIN"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast2.RaiseStatement:
		return "raise " + walker(n.X)
	case *ast2.TryStatement:
		result := "try"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		for _, exceptBlock := range n.ExceptBlocks {
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
		if n.Else != nil {
			result += "\nelse"
			for _, bodyNode := range n.Else {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		if n.Finally != nil {
			result += "\nfinally"
			for _, bodyNode := range n.Finally {
				nodeString := walker(bodyNode)
				nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
				result += "\n\t" + nodeString
			}
		}
		return result + "\nend"
	case *ast2.DoWhileStatement:
		result := "do"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nwhile " + walker(n.Condition)
	case *ast2.SuperExpression:
		return "super " + walker(n.X)
	case *ast2.RequireStatement:
		return "require " + walker(n.X)
	case *ast2.DeleteStatement:
		return "delete " + walker(n.X)
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
	"(super DataType).Initialize",
	"require \"../my/script.pm\"",
	"delete a.b",
	"delete a[b]",
	"gen a()\n\tyield 1\nend",
	"1 implements Integer",
	"1 is Integer",
	"a_function(1)",
	"a_function(1, 2)",
	"a_function(1, 2, call((1, 2, 3, 4, 2 * (5 - 3))))",
	"if a > 2\n" +
		"\tcall()\n" +
		"elif a < 2\n" +
		"\tif a == 0\n" +
		"\t\tprint(\"\\\"a\\\" is zero\")\n" +
		"\telse\n" +
		"\t\tprint(\"\\\"a\\\" is non zero\")" +
		"\n\tend\n" +
		"end",
	"if a > 2\n" +
		"\tcall()\n" +
		"elif a < 2\n" +
		"\tprintln(1)\n" +
		"elif b == 2\n" +
		"\tprint(3)\n" +
		"else\n" +
		"\texit(1)\n" +
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
	"unless a > 2\n" +
		"\tcall()\n" +
		"elif a < 2\n" +
		"\tprintln(1)\n" +
		"elif b == 2\n" +
		"\tprint(3)\n" +
		"else\n" +
		"\texit(1)\n" +
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
