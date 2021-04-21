package parser

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/ast"
	"github.com/shoriwe/gruby/pkg/compiler/lexer"
	"testing"
)

func walker(node ast.Node, deep int) {
	switch node.(type) {
	case *ast.Program:
		for _, child := range node.(*ast.Program).Body {
			walker(child, deep+1)
		}
	case *ast.BinaryExpression:
		fmt.Print(node.(*ast.BinaryExpression).Operator)
		walker(node.(*ast.BinaryExpression).LeftHandSide, deep+1)
		walker(node.(*ast.BinaryExpression).RightHandSide, deep+1)
	case *ast.BasicLiteralExpression:
		fmt.Print(node.(*ast.BasicLiteralExpression).String)
	case *ast.UnaryExpression:
		fmt.Print(node.(*ast.UnaryExpression).Operator)
		walker(node.(*ast.UnaryExpression).X, deep+1)
	case *ast.SelectorExpression:
		walker(node.(*ast.SelectorExpression).X, deep+1)
		fmt.Print("." + node.(*ast.SelectorExpression).Identifier.String)
	case *ast.Identifier:
		fmt.Print(node.(*ast.Identifier).String)
	case *ast.MethodInvocationExpression:
		walker(node.(*ast.MethodInvocationExpression).Function, deep+1)
		fmt.Print("(")
		for _, child := range node.(*ast.MethodInvocationExpression).Arguments {
			walker(child, deep+1)
		}
		fmt.Print(")")
	}
}

func walk(node ast.Node) {
	walker(node, 0)
	fmt.Println("")
}

func test(t *testing.T, samples []string) {
	for _, sample := range samples {
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
		walk(program)
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
}

func TestParseBasic(t *testing.T) {
	test(t, basicSamples)
}
