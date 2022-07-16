package parser

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/lexer"
	reader2 "github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/test-samples/basic"
	"reflect"
	"strings"
	"testing"
)

func walker(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Program:
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
	case *ast.BinaryExpression:
		return walker(n.LeftHandSide) +
			" " + n.Operator.String() +
			" " + walker(n.RightHandSide)
	case *ast.BasicLiteralExpression:
		return n.Token.String()
	case *ast.UnaryExpression:
		if n.Operator.DirectValue == lexer.Not {
			return n.Operator.String() + " " + walker(n.X)
		}
		return n.Operator.String() + walker(n.X)
	case *ast.SelectorExpression:
		return walker(n.X) + "." + n.Identifier.Token.String()
	case *ast.Identifier:
		return n.Token.String()
	case *ast.MethodInvocationExpression:
		result := walker(n.Function) + "("
		for index, child := range n.Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(child)
		}
		return result + ")"
	case *ast.IndexExpression:
		result := walker(n.Source) + "["
		result += walker(n.Index)
		return result + "]"
	case *ast.LambdaExpression:
		result := "lambda "
		for index, argument := range n.Arguments {
			if index != 0 {
				result += ", "
			}
			result += walker(argument)
		}
		result += ": "
		return result + walker(n.Code)
	case *ast.ParenthesesExpression:
		return "(" + walker(n.X) + ")"
	case *ast.TupleExpression:
		result := "("
		for index, value := range n.Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + ")"
	case *ast.ArrayExpression:
		result := "["
		for index, value := range n.Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value)
		}
		return result + "]"
	case *ast.HashExpression:
		result := "{"
		for index, value := range n.Values {
			if index != 0 {
				result += ", "
			}
			result += walker(value.Key)
			result += ": " + walker(value.Value)
		}
		return result + "}"
	case *ast.IfOneLinerExpression:
		result := walker(n.Result)
		result += " if " + walker(n.Condition)
		if n.ElseResult != nil {
			result += " else " + walker(n.ElseResult)
		}
		return result
	case *ast.UnlessOneLinerExpression:
		result := walker(n.Result)
		result += " unless "
		result += walker(n.Condition)
		if n.ElseResult != nil {
			result += " else " + walker(n.ElseResult)
		}
		return result
	case *ast.GeneratorExpression:
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
	case *ast.AssignStatement:
		result := walker(n.LeftHandSide)
		result += " " + n.AssignOperator.String() + " "
		return result + walker(n.RightHandSide)
	case *ast.ContinueStatement:
		return "continue"
	case *ast.BreakStatement:
		return "break"
	case *ast.PassStatement:
		return "pass"
	case *ast.YieldStatement:
		result := "yield "
		for index, output := range n.Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast.ReturnStatement:
		result := "return "
		for index, output := range n.Results {
			if index != 0 {
				result += ", "
			}
			result += walker(output)
		}
		return result
	case *ast.IfStatement:
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
	case *ast.UnlessStatement:
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
	case *ast.SwitchStatement:
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
	case *ast.WhileLoopStatement:
		result := "while " + walker(n.Condition)
		for _, child := range n.Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast.UntilLoopStatement:
		result := "until " + walker(n.Condition)
		for _, child := range n.Body {
			nodeString := walker(child)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast.ForLoopStatement:
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
	case *ast.ModuleStatement:
		result := "module " + walker(n.Name)
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast.ClassStatement:
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
	case *ast.InterfaceStatement:
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
	case *ast.FunctionDefinitionStatement:
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
	case *ast.GeneratorDefinitionStatement:
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
	case *ast.EndStatement:
		result := "END"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast.BeginStatement:
		result := "BEGIN"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast.BlockStatement:
		result := "block"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nend"
	case *ast.DoWhileStatement:
		result := "do"
		for _, bodyNode := range n.Body {
			nodeString := walker(bodyNode)
			nodeString = strings.ReplaceAll(nodeString, "\n", "\n\t")
			result += "\n\t" + nodeString
		}
		return result + "\nwhile " + walker(n.Condition)
	case *ast.SuperExpression:
		return "super " + walker(n.X)
	case *ast.DeleteStatement:
		return "delete " + walker(n.X)
	case *ast.DeferStatement:
		return "defer " + walker(n.X)
	}
	panic("unknown node type: " + reflect.TypeOf(node).String())
}

func walk(node ast.Node) string {
	return walker(node)
}

func test(t *testing.T, samples map[string]string) {
	for sampleScript, sample := range samples {
		lex := lexer.NewLexer(reader2.NewStringReader(sample))
		parser := NewParser(lex)
		program, parsingError := parser.Parse()
		if parsingError != nil {
			t.Error(fmt.Sprintf("%s in sample %s", parsingError.Error(), sampleScript))
			return
		}
		result := walk(program)
		if len(result) != 0 {
			result = result[:len(result)-1]
		}
		if result == sample {
			t.Logf("\nSample: %s -> SUCCESS", sampleScript)
		} else {
			t.Errorf("\nSample: %s -> FAIL", sampleScript)
			fmt.Println(sample)
			fmt.Println(result)
		}
	}
}

func TestParseBasic(t *testing.T) {
	test(t, basic.Samples)
}
