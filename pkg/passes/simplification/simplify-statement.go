package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifyStatement(stmt ast.Statement) ast2.Statement {
	switch s := stmt.(type) {
	case *ast.AssignStatement:
		return simplifyAssign(s)
	case *ast.DoWhileStatement:
		return simplifyDoWhile(s)
	case *ast.WhileLoopStatement:
		return simplifyWhile(s)
	case *ast.UntilLoopStatement:
		return simplifyUntil(s)
	case *ast.ForLoopStatement:
		return simplifyFor(s)
	case *ast.IfStatement:
		return simplifyIf(s)
	case *ast.UnlessStatement:
		return simplifyUnless(s)
	case *ast.SwitchStatement:
		return simplifySwitch(s)
	case *ast.ModuleStatement:
		return simplifyModule(s)
	case *ast.FunctionDefinitionStatement:
		return simplifyFunction(s)
	case *ast.GeneratorDefinitionStatement:
		return simplifyGeneratorDef(s)
	case *ast.InterfaceStatement:
		return simplifyInterface(s)
	case *ast.ClassStatement:
		return simplifyClass(s)
	case *ast.TryStatement:
		return simplifyTry(s)
	case *ast.ReturnStatement:
		return simplifyReturn(s)
	case *ast.YieldStatement:
		return simplifyYield(s)
	case *ast.ContinueStatement:
		return simplifyContinue(s)
	case *ast.BreakStatement:
		return simplifyBreak(s)
	case *ast.PassStatement:
		return simplifyPass(s)
	case *ast.RaiseStatement:
		return simplifyRaise(s)
	case *ast.BlockStatement:
		return simplifyBlock(s)
	case *ast.RequireStatement:
		return simplifyRequire(s)
	case *ast.DeleteStatement:
		return simplifyDelete(s)
	default:
		panic("unknown statement type")
	}
}
